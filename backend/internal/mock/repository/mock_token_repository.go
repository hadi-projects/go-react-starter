package mock_repository

import (
	context "context"
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
func (m *MockTokenRepository) Create(ctx context.Context, token *entity.PasswordResetToken) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockTokenRepositoryMockRecorder) Create(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTokenRepository)(nil).Create), ctx, token)
}

// Delete mocks base method.
func (m *MockTokenRepository) Delete(ctx context.Context, token *entity.PasswordResetToken) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockTokenRepositoryMockRecorder) Delete(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTokenRepository)(nil).Delete), ctx, token)
}

// DeleteByUserID mocks base method.
func (m *MockTokenRepository) DeleteByUserID(ctx context.Context, userID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByUserID", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByUserID indicates an expected call of DeleteByUserID.
func (mr *MockTokenRepositoryMockRecorder) DeleteByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByUserID", reflect.TypeOf((*MockTokenRepository)(nil).DeleteByUserID), ctx, userID)
}

// FindByToken mocks base method.
func (m *MockTokenRepository) FindByToken(ctx context.Context, token string) (*entity.PasswordResetToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByToken", ctx, token)
	ret0, _ := ret[0].(*entity.PasswordResetToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByToken indicates an expected call of FindByToken.
func (mr *MockTokenRepositoryMockRecorder) FindByToken(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByToken", reflect.TypeOf((*MockTokenRepository)(nil).FindByToken), ctx, token)
}
