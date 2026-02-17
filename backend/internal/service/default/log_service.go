package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hadi-projects/go-react-starter/config"
	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
)

type LogService interface {
	GetLogs(query dto.LogQuery) ([]dto.LogResponse, error)
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
	if query.Type == "all" || query.Type == "auth" {
		filesToRead = append(filesToRead, "auth.log")
	}
	if query.Type == "all" || query.Type == "audit" {
		filesToRead = append(filesToRead, "audit.log")
	}
	if query.Type == "all" || query.Type == "system" {
		filesToRead = append(filesToRead, "system.log")
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
