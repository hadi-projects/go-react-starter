package mock_repository

import (
	reflect "reflect"

	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
	gomock "go.uber.org/mock/gomock"
)

// MockTokenRepository is a mock of TokenRepository interface.
type MockTokenRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTokenRepositoryMockRecorder
}

// MockTokenRepositoryMockRecorder is the mock recorder for MockTokenRepository.
type MockTokenRepositoryMockRecorder struct {
	mock *MockTokenRepository
}

// NewMockTokenRepository creates a new mock instance.
func NewMockTokenRepository(ctrl *gomock.Controller) *MockTokenRepository {
	mock := &MockTokenRepository{ctrl: ctrl}
	mock.recorder = &MockTokenRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenRepository) EXPECT() *MockTokenRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockTokenRepository) Create(token *entity.PasswordResetToken) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", token)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockTokenRepositoryMockRecorder) Create(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTokenRepository)(nil).Create), token)
}

// Delete mocks base method.
func (m *MockTokenRepository) Delete(token *entity.PasswordResetToken) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", token)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockTokenRepositoryMockRecorder) Delete(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTokenRepository)(nil).Delete), token)
}

// DeleteByUserID mocks base method.
func (m *MockTokenRepository) DeleteByUserID(userID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByUserID", userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByUserID indicates an expected call of DeleteByUserID.
func (mr *MockTokenRepositoryMockRecorder) DeleteByUserID(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByUserID", reflect.TypeOf((*MockTokenRepository)(nil).DeleteByUserID), userID)
}

// FindByToken mocks base method.
func (m *MockTokenRepository) FindByToken(token string) (*entity.PasswordResetToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByToken", token)
	ret0, _ := ret[0].(*entity.PasswordResetToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByToken indicates an expected call of FindByToken.
func (mr *MockTokenRepositoryMockRecorder) FindByToken(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByToken", reflect.TypeOf((*MockTokenRepository)(nil).FindByToken), token)
}
