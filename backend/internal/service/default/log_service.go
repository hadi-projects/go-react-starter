package service

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hadi-projects/go-react-starter/config"
	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	"github.com/xuri/excelize/v2"
)

type LogService interface {
	GetLogs(query dto.LogQuery) ([]dto.LogResponse, error)
	Export(query dto.LogQuery, format string) ([]byte, string, error)
}

type logService struct {
	config *config.Config
}

func NewLogService(config *config.Config) LogService {
	return &logService{config: config}
}

func (s *logService) GetLogs(query dto.LogQuery) ([]dto.LogResponse, error) {
	var allLogs []dto.LogResponse

	filesToRead := []string{}
	if query.Type == "" {
		query.Type = "all"
	}
	if query.Type == "all" || query.Type == "audit" {
		filesToRead = append(filesToRead, "audit.log")
	}
	if query.Type == "all" || query.Type == "system" {
		filesToRead = append(filesToRead, "system.log", "db.log", "redis.log", "rate_limit.log")
	}

	for _, fileName := range filesToRead {
		filePath := filepath.Join(s.config.Log.Dir, fileName)
		fmt.Printf("DEBUG: Reading log file: %s\n", filePath)
		logs, err := s.readLogFile(filePath, strings.TrimSuffix(fileName, ".log"))
		if err != nil {
			fmt.Printf("DEBUG: Error reading %s: %v\n", filePath, err)
			// Skip if file doesn't exist yet
			continue
		}
		fmt.Printf("DEBUG: Successfully read %d logs from %s\n", len(logs), filePath)
		allLogs = append(allLogs, logs...)
	}

	// Filter by UserID if provided
	if query.UserID != 0 {
		var filteredLogs []dto.LogResponse
		for _, log := range allLogs {
			if log.UserID != nil && *log.UserID == query.UserID {
				filteredLogs = append(filteredLogs, log)
			}
		}
		allLogs = filteredLogs
	}

	// Sort logs by time descending
	sort.Slice(allLogs, func(i, j int) bool {
		return allLogs[i].Time.After(allLogs[j].Time)
	})

	return allLogs, nil
}

func (s *logService) readLogFile(filePath string, logType string) ([]dto.LogResponse, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var logs []dto.LogResponse
	scanner := bufio.NewScanner(file)
	// Increase buffer to 10MB to handle large log lines (system.log can have very long lines)
	const maxScanTokenSize = 10 * 1024 * 1024 // 10MB
	buf := make([]byte, maxScanTokenSize)
	scanner.Buffer(buf, maxScanTokenSize)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var raw map[string]interface{}
		if err := json.Unmarshal([]byte(line), &raw); err != nil {
			continue
		}

		log := dto.LogResponse{
			Type: logType,
		}

		if val, ok := raw["level"].(string); ok {
			log.Level = val
		}
		if val, ok := raw["action"].(string); ok {
			log.Action = val
		}
		if val, ok := raw["message"].(string); ok {
			log.Message = val
		}
		if val, ok := raw["email"].(string); ok {
			log.Email = val
		}
		if val, ok := raw["request_id"].(string); ok {
			log.RequestID = val
		}

		// Populate Source
		if val, ok := raw["method"].(string); ok && val != "" {
			log.Source = strings.ToLower(val)
		} else {
			log.Source = logType
		}

		// Handle user_id (it could be uint or float64 from json)
		if val, ok := raw["user_id"]; ok {
			switch v := val.(type) {
			case float64:
				u := uint(v)
				log.UserID = &u
			case uint:
				log.UserID = &v
			}
		}

		if val, ok := raw["time"].(string); ok {
			if t, err := json.Marshal(val); err == nil {
				if err := json.Unmarshal(t, &log.Time); err != nil {
					fmt.Printf("DEBUG: Failed to unmarshal time (time): %v, val: %s\n", err, val)
				}
			}
		} else if val, ok := raw["timestamp"].(string); ok {
			// System logs use 'timestamp' instead of 'time'
			if t, err := json.Marshal(val); err == nil {
				if err := json.Unmarshal(t, &log.Time); err != nil {
					fmt.Printf("DEBUG: Failed to unmarshal time (timestamp): %v, val: %s\n", err, val)
				}
			}
		} else {
			fmt.Printf("DEBUG: No time or timestamp found in log: %v\n", raw)
		}

		// Collect other fields into Details
		delete(raw, "level")
		delete(raw, "action")
		delete(raw, "message")
		delete(raw, "user_id")
		delete(raw, "email")
		delete(raw, "time")
		delete(raw, "timestamp")
		delete(raw, "request_id")
		log.Details = raw

		logs = append(logs, log)
	}

	return logs, nil
}

func (s *logService) Export(query dto.LogQuery, format string) ([]byte, string, error) {
	logs, err := s.GetLogs(query)
	if err != nil {
		return nil, "", err
	}

	if format == "csv" {
		return s.generateCSV(logs)
	}
	return s.generateExcel(logs)
}

func (s *logService) generateCSV(logs []dto.LogResponse) ([]byte, string, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	header := []string{"Time", "Type", "Source", "Action", "Message", "Email", "Request ID"}
	if err := writer.Write(header); err != nil {
		return nil, "", err
	}

	for _, l := range logs {
		row := []string{
			l.Time.Format("2006-01-02 15:04:05"),
			l.Type,
			l.Source,
			l.Action,
			l.Message,
			l.Email,
			l.RequestID,
		}
		if err := writer.Write(row); err != nil {
			return nil, "", err
		}
	}

	writer.Flush()
	return buf.Bytes(), "system_logs.csv", nil
}

func (s *logService) generateExcel(logs []dto.LogResponse) ([]byte, string, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Logs"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	headers := []string{"Time", "Type", "Source", "Action", "Message", "Email", "Request ID"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, l := range logs {
		row := []interface{}{
			l.Time.Format("2006-01-02 15:04:05"),
			l.Type,
			l.Source,
			l.Action,
			l.Message,
			l.Email,
			l.RequestID,
		}
		for j, val := range row {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+2)
			f.SetCellValue(sheet, cell, val)
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, "", err
	}

	return buf.Bytes(), "system_logs.xlsx", nil
}
