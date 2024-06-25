package loan

import (
	"context"
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gitlab.com/2024/Juni/amartha-billing-srv2/infrastructure/repository"
	mocks "gitlab.com/2024/Juni/amartha-billing-srv2/mocks/configuration"
	mocks2 "gitlab.com/2024/Juni/amartha-billing-srv2/mocks/infrastructure/repository"
)

func Test_loanService_FetchOutstanding(t *testing.T) {
	mockConfig := &mocks.Configuration{}
	mockLoanRepo := &mocks2.LoanRepository{}

	type args struct {
		uid string
	}

	lePendingCriteria := []*repository.LoanEntity{
		{
			Status: "PENDING",
			Amount: decimal.NewFromFloat(float64(10)),
		},
		{
			Status: "PENDING",
			Amount: decimal.NewFromFloat(float64(13)),
		},
		{
			Status: "PENDING",
			Amount: decimal.NewFromFloat(float64(5)),
		},
		{
			Status: "PENDING",
			Amount: decimal.NewFromFloat(float64(2)),
		},
	}

	pendingAmountOutstanding := decimal.NewFromFloat(float64(0))

	for _, val := range lePendingCriteria {
		if val.Status == "PENDING" {
			pendingAmountOutstanding = pendingAmountOutstanding.Add(val.Amount)
		}
	}

	leClosedCriteria := []*repository.LoanEntity{
		{
			Status: "CLOSED",
			Amount: decimal.NewFromFloat(float64(10)),
		},
		{
			Status: "CLOSED",
			Amount: decimal.NewFromFloat(float64(13)),
		},
		{
			Status: "CLOSED",
			Amount: decimal.NewFromFloat(float64(5)),
		},
	}

	lePaidCriteria := []*repository.LoanEntity{
		{
			Status: "PAID",
			Amount: decimal.NewFromFloat(float64(10)),
		},
		{
			Status: "PAID",
			Amount: decimal.NewFromFloat(float64(13)),
		},
	}

	tests := []struct {
		name     string
		args     args
		want     *FetchOutstandingResponse
		wantErr  error
		mockFunc func()
	}{
		{
			name: "given uid is empty string," +
				"when fetchOutstanding," +
				"then return error",
			args: args{
				uid: "",
			},
			want:    nil,
			wantErr: errorValidation,
			mockFunc: func() {

			},
		},
		{
			name: "given has error: sql no rows," +
				"when fetchOutstanding:findLoans," +
				"then return error",
			args: args{
				uid: "abc",
			},
			want:    nil,
			wantErr: errorDataNotExists,
			mockFunc: func() {
				mockLoanRepo.
					On("FindLoans", mock.Anything, mock.Anything).
					Return(nil, repository.ErrorNoRows).
					Once()
			},
		},
		{
			name: "given has error: unknown from database," +
				"when fetchOutstanding:findLoans," +
				"then return error",
			args: args{
				uid: "abc",
			},
			want:    nil,
			wantErr: errorFromDatabase,
			mockFunc: func() {
				mockLoanRepo.
					On("FindLoans", mock.Anything, mock.Anything).
					Return(nil, errors.New("new error")).
					Once()
			},
		},
		{
			name: "given valid request and works as expected," +
				"when fetchOutstanding:findLoans," +
				"then return outstanding pending",
			args: args{
				uid: "abc",
			},
			want: &FetchOutstandingResponse{
				RemainingOutstanding: pendingAmountOutstanding,
				IsDelinquent:         true,
			},
			wantErr: nil,
			mockFunc: func() {
				mockLoanRepo.
					On("FindLoans", mock.Anything, mock.Anything).
					Return(lePendingCriteria, nil).
					Once()
			},
		},
		{
			name: "given valid request and works as expected," +
				"when fetchOutstanding:findLoans," +
				"then return outstanding closed",
			args: args{
				uid: "abc",
			},
			want: &FetchOutstandingResponse{
				RemainingOutstanding: decimal.NewFromFloat(float64(0)),
				IsDelinquent:         false,
			},
			wantErr: nil,
			mockFunc: func() {
				mockLoanRepo.
					On("FindLoans", mock.Anything, mock.Anything).
					Return(leClosedCriteria, nil).
					Once()
			},
		},
		{
			name: "given valid request and works as expected," +
				"when fetchOutstanding:findLoans," +
				"then return outstanding paid",
			args: args{
				uid: "abc",
			},
			want: &FetchOutstandingResponse{
				RemainingOutstanding: decimal.NewFromFloat(float64(0)),
				IsDelinquent:         false,
			},
			wantErr: nil,
			mockFunc: func() {
				mockLoanRepo.
					On("FindLoans", mock.Anything, mock.Anything).
					Return(lePaidCriteria, nil).
					Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				l := NewLoanService(mockConfig, mockLoanRepo)
				tt.mockFunc()

				got, err := l.FetchOutstanding(context.Background(), tt.args.uid)

				if got != nil {
					assert.Equal(t, got.RemainingOutstanding, tt.want.RemainingOutstanding)
					assert.Equal(t, got.IsDelinquent, tt.want.IsDelinquent)
				}

				assert.Equal(t, err, tt.wantErr)
			})
	}
}
