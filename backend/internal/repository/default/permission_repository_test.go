package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type PermissionRepositoryTestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	repo PermissionRepository
}

func (s *PermissionRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	s.Require().NoError(err)
	s.mock = mock

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	s.Require().NoError(err)

	s.repo = NewPermissionRepository(gormDB)
}

func (s *PermissionRepositoryTestSuite) TestCreate_Success() {
	permission := &entity.Permission{
		Name: "test-permission",
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `permissions`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err := s.repo.Create(context.Background(), permission)
	s.Require().NoError(err)
}

func (s *PermissionRepositoryTestSuite) TestFindByID_Success() {
	id := uint(1)
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at"}).
		AddRow(id, "test-permission", "desc", now, now)

	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `permissions` WHERE `permissions`.`id` = ?")).
		WithArgs(id, 1).
		WillReturnRows(rows)

	permission, err := s.repo.FindByID(context.Background(), id)
	s.Require().NoError(err)
	s.Require().NotNil(permission)
	assert.Equal(s.T(), "test-permission", permission.Name)
}

func TestPermissionRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PermissionRepositoryTestSuite))
}
