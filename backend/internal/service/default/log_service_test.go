package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hadi-projects/go-react-starter/config"
	"github.com/hadi-projects/go-react-starter/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LogServiceTestSuite struct {
	suite.Suite
	service LogService
	tempDir string
}

func (s *LogServiceTestSuite) SetupSuite() {
	dir, err := os.MkdirTemp("", "logtest")
	s.Require().NoError(err)
	s.tempDir = dir
}

func (s *LogServiceTestSuite) TearDownSuite() {
	os.RemoveAll(s.tempDir)
}

func (s *LogServiceTestSuite) SetupTest() {
	cfg := &config.Config{
		Log: config.LogConfig{
			Dir: s.tempDir,
		},
	}
	s.service = NewLogService(cfg)
}

func (s *LogServiceTestSuite) TestGetLogs_Success() {
	// Create a dummy audit log file
	logLine := `{"level":"info","action":"user_login","message":"login successful","time":"2026-02-17T07:30:00Z","user_id":1}`
	err := os.WriteFile(filepath.Join(s.tempDir, "audit.log"), []byte(logLine+"\n"), 0644)
	s.Require().NoError(err)

	query := dto.LogQuery{Type: "audit"}
	logs, err := s.service.GetLogs(query)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), logs, 1)
	assert.Equal(s.T(), "user_login", logs[0].Action)
	assert.Equal(s.T(), uint(1), *logs[0].UserID)
}

func TestLogServiceTestSuite(t *testing.T) {
	suite.Run(t, new(LogServiceTestSuite))
}
