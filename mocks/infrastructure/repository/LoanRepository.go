// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	repository "gitlab.com/2024/Juni/amartha-billing-srv2/infrastructure/repository"

	sql "database/sql"
)

// LoanRepository is an autogenerated mock type for the LoanRepository type
type LoanRepository struct {
	mock.Mock
}

// FindLoans provides a mock function with given fields: ctx, loanEntity
func (_m *LoanRepository) FindLoans(ctx context.Context, loanEntity *repository.LoanEntity) ([]*repository.LoanEntity, error) {
	ret := _m.Called(ctx, loanEntity)

	if len(ret) == 0 {
		panic("no return value specified for FindLoans")
	}

	var r0 []*repository.LoanEntity
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *repository.LoanEntity) ([]*repository.LoanEntity, error)); ok {
		return rf(ctx, loanEntity)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *repository.LoanEntity) []*repository.LoanEntity); ok {
		r0 = rf(ctx, loanEntity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*repository.LoanEntity)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *repository.LoanEntity) error); ok {
		r1 = rf(ctx, loanEntity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveLoans provides a mock function with given fields: ctx, tx, loanEntity
func (_m *LoanRepository) SaveLoans(ctx context.Context, tx *sql.Tx, loanEntity ...*repository.LoanEntity) error {
	_va := make([]interface{}, len(loanEntity))
	for _i := range loanEntity {
		_va[_i] = loanEntity[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, tx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for SaveLoans")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *sql.Tx, ...*repository.LoanEntity) error); ok {
		r0 = rf(ctx, tx, loanEntity...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateLoan provides a mock function with given fields: ctx, loanEntity
func (_m *LoanRepository) UpdateLoan(ctx context.Context, loanEntity *repository.LoanEntityUpdate) error {
	ret := _m.Called(ctx, loanEntity)

	if len(ret) == 0 {
		panic("no return value specified for UpdateLoan")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *repository.LoanEntityUpdate) error); ok {
		r0 = rf(ctx, loanEntity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewLoanRepository creates a new instance of LoanRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLoanRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *LoanRepository {
	mock := &LoanRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}