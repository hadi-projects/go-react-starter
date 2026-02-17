package repository

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	repo UserRepository
}

func (s *UserRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	s.Require().NoError(err)
	s.mock = mock

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	s.Require().NoError(err)

	s.repo = NewUserRepository(gormDB)
}

func (s *UserRepositoryTestSuite) TestCreate_Success() {
	user := &entity.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password",
		RoleID:   1,
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err := s.repo.Create(user)
	s.Require().NoError(err)
}

func (s *UserRepositoryTestSuite) TestFindByEmail_Success() {
	email := "test@example.com"
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role_id", "created_at", "updated_at"}).
		AddRow(1, "Test User", email, "hashed", 1, now, now)

	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE email = ? ORDER BY `users`.`id` LIMIT ?")).
		WithArgs(email, 1).
		WillReturnRows(rows)

	// Preload Role
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `roles` WHERE `roles`.`id` = ?")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "admin"))

	// Preload Permissions for Role (Many-to-Many)
	// 1. Fetch from join table
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `role_has_permissions` WHERE `role_has_permissions`.`role_id` = ?")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"role_id", "permission_id"}).AddRow(1, 1))

	// 2. Fetch permissions by ID
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `permissions` WHERE `permissions`.`id` = ?")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "create-user"))

	user, err := s.repo.FindByEmail(email)
	s.Require().NoError(err)
	s.Require().NotNil(user)
	assert.Equal(s.T(), email, user.Email)
}

func (s *UserRepositoryTestSuite) TestUpdate_Success() {
	user := &entity.User{
		ID:    1,
		Name:  "Updated User",
		Email: "updated@example.com",
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err := s.repo.Update(user)
	s.Require().NoError(err)
}

func (s *UserRepositoryTestSuite) TestDelete_Success() {
	userID := uint(1)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `users`")).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err := s.repo.Delete(userID)
	s.Require().NoError(err)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
