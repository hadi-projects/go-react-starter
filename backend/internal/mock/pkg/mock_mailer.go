package mock_pkg

import (
	"context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockMailer is a mock of Mailer interface.
type MockMailer struct {
	ctrl     *gomock.Controller
	recorder *MockMailerMockRecorder
}

// MockMailerMockRecorder is the mock recorder for MockMailer.
type MockMailerMockRecorder struct {
	mock *MockMailer
}

// NewMockMailer creates a new mock instance.
func NewMockMailer(ctrl *gomock.Controller) *MockMailer {
	mock := &MockMailer{ctrl: ctrl}
	mock.recorder = &MockMailerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMailer) EXPECT() *MockMailerMockRecorder {
	return m.recorder
}

// SendEmail mocks base method.
func (m *MockMailer) SendEmail(ctx context.Context, to, subject, body string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendEmail", ctx, to, subject, body)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendEmail indicates an expected call of SendEmail.
func (mr *MockMailerMockRecorder) SendEmail(ctx, to, subject, body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendEmail", reflect.TypeOf((*MockMailer)(nil).SendEmail), ctx, to, subject, body)
}
