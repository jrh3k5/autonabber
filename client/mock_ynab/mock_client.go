// Code generated by MockGen. DO NOT EDIT.
// Source: client/ynab/client.go

// Package mock_ynab is a generated GoMock package.
package mock_ynab

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/jrh3k5/autonabber/client/ynab/model"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// GetBudgets mocks base method.
func (m *MockClient) GetBudgets() ([]*model.Budget, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBudgets")
	ret0, _ := ret[0].([]*model.Budget)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBudgets indicates an expected call of GetBudgets.
func (mr *MockClientMockRecorder) GetBudgets() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBudgets", reflect.TypeOf((*MockClient)(nil).GetBudgets))
}

// GetCategories mocks base method.
func (m *MockClient) GetCategories(budget *model.Budget) ([]*model.BudgetCategoryGroup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCategories", budget)
	ret0, _ := ret[0].([]*model.BudgetCategoryGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCategories indicates an expected call of GetCategories.
func (mr *MockClientMockRecorder) GetCategories(budget interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCategories", reflect.TypeOf((*MockClient)(nil).GetCategories), budget)
}

// GetMonthlyAverageSpent mocks base method.
func (m *MockClient) GetMonthlyAverageSpent(budget *model.Budget, category *model.BudgetCategory, monthLookback int) (int64, int16, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMonthlyAverageSpent", budget, category, monthLookback)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(int16)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetMonthlyAverageSpent indicates an expected call of GetMonthlyAverageSpent.
func (mr *MockClientMockRecorder) GetMonthlyAverageSpent(budget, category, monthLookback interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMonthlyAverageSpent", reflect.TypeOf((*MockClient)(nil).GetMonthlyAverageSpent), budget, category, monthLookback)
}

// GetReadyToAssign mocks base method.
func (m *MockClient) GetReadyToAssign(budget *model.Budget) (int64, int16, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReadyToAssign", budget)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(int16)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetReadyToAssign indicates an expected call of GetReadyToAssign.
func (mr *MockClientMockRecorder) GetReadyToAssign(budget interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReadyToAssign", reflect.TypeOf((*MockClient)(nil).GetReadyToAssign), budget)
}

// SetBudget mocks base method.
func (m *MockClient) SetBudget(budget *model.Budget, category *model.BudgetCategory, newDollars int64, newCents int16) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetBudget", budget, category, newDollars, newCents)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetBudget indicates an expected call of SetBudget.
func (mr *MockClientMockRecorder) SetBudget(budget, category, newDollars, newCents interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBudget", reflect.TypeOf((*MockClient)(nil).SetBudget), budget, category, newDollars, newCents)
}
