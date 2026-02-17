package repository

import (
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

type RoleRepositoryTestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	repo RoleRepository
}

func (s *RoleRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	s.Require().NoError(err)
	s.mock = mock

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	s.Require().NoError(err)

	s.repo = NewRoleRepository(gormDB)
}

func (s *RoleRepositoryTestSuite) TestCreate_Success() {
	role := &entity.Role{
		Name: "admin",
	}
	permissionIDs := []uint{1}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `roles`")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Flexible expectations for associations
	s.mock.ExpectQuery("SELECT .* FROM `permissions` .*").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "create-user"))

	s.mock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectCommit()

	err := s.repo.Create(role, permissionIDs)
	s.Require().NoError(err)
}

func (s *RoleRepositoryTestSuite) TestFindByID_Success() {
	id := uint(1)
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at"}).
		AddRow(id, "admin", "", now, now)

	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `roles` WHERE `roles`.`id` = ? ORDER BY `roles`.`id` LIMIT ?")).
		WithArgs(id, 1).
		WillReturnRows(rows)

	// Preload Permissions (Many-to-Many)
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `role_has_permissions` WHERE `role_has_permissions`.`role_id` = ?")).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"role_id", "permission_id"}).AddRow(id, 1))

	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `permissions` WHERE `permissions`.`id` = ?")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "create-user"))

	role, err := s.repo.FindByID(id)
	s.Require().NoError(err)
	s.Require().NotNil(role)
	assert.Equal(s.T(), "admin", role.Name)
}

func TestRoleRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RoleRepositoryTestSuite))
}
