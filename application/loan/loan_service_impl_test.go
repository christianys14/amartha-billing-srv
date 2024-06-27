package loan

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gitlab.com/2024/Juni/amartha-billing-srv2/infrastructure/repository"
	mocks2 "gitlab.com/2024/Juni/amartha-billing-srv2/mocks/infrastructure/repository"
)

func Test_loanService_FetchOutstanding(t *testing.T) {
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

	tests := []struct {
		name     string
		args     args
		want     *FetchOutstandingResponse
		wantErr  error
		mockFunc func()
	}{
		{
			name: "given intentionally panic," +
				"when fetchOutstanding," +
				"then return error",
			args: args{
				uid: "asd",
			},
			want:    nil,
			wantErr: errorFromDatabase,
			mockFunc: func() {
				mockLoanRepo.
					On(
						"FindLoans", mock.Anything, repository.LoanEntity{
							Statuses: []string{"PENDING", "CLOSED"},
							UserID:   "asd",
							DueDate:  time.Now(),
						}).
					Return(nil, nil).
					Once()
			},
		},
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
					Return(
						[]*repository.LoanEntity{
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
						}, nil).
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
					Return(
						[]*repository.LoanEntity{
							{
								Status: "PAID",
								Amount: decimal.NewFromFloat(float64(10)),
							},
							{
								Status: "PAID",
								Amount: decimal.NewFromFloat(float64(13)),
							},
						}, nil).
					Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				l := NewLoanService(mockLoanRepo)
				tt.mockFunc()

				got, err := l.FetchOutstanding(context.Background(), tt.args.uid)

				if got != nil {
					assert.Equal(t, tt.want.RemainingOutstanding, got.RemainingOutstanding)
					assert.Equal(t, tt.want.IsDelinquent, got.IsDelinquent)
				}

				assert.Equal(t, tt.wantErr, err)
			})
	}
}

func Test_loanService_Payment(t *testing.T) {
	mockLoanRepo := &mocks2.LoanRepository{}

	payReq := &PaymentRequest{
		UserID: "abc",
		Amount: float64(25),
	}

	type args struct {
		paymentRequest *PaymentRequest
	}
	tests := []struct {
		name     string
		args     args
		wantErr  error
		mockFunc func()
	}{
		{
			name: "given intentionally panic," +
				"when payment," +
				"then return error",
			args: args{
				paymentRequest: payReq,
			},
			wantErr: errorFromDatabase,
			mockFunc: func() {
				mockLoanRepo.
					On(
						"FindLoans", mock.Anything, repository.LoanEntity{
							Statuses: []string{"PENDING", "CLOSED"},
							UserID:   "asd",
							DueDate:  time.Now(),
						}).
					Return(nil, nil).
					Once()
			},
		},
		{
			name: "given not passed the validation," +
				"when payment," +
				"then return error",
			args: args{
				paymentRequest: &PaymentRequest{},
			},
			wantErr: errorValidation,
			mockFunc: func() {
			},
		},
		{
			name: "given no rows after looking for from db," +
				"when payment," +
				"then return error",
			args: args{
				paymentRequest: payReq,
			},
			wantErr: errorNoPendingOutstanding,
			mockFunc: func() {
				mockLoanRepo.
					On("FindLoans", mock.Anything, mock.Anything).
					Return(nil, nil).
					Once()
			},
		},
		{
			name: "given unknown error after looking for from db," +
				"when payment," +
				"then return error",
			args: args{
				paymentRequest: payReq,
			},
			wantErr: errorFromDatabase,
			mockFunc: func() {
				mockLoanRepo.
					On("FindLoans", mock.Anything, mock.Anything).
					Return(nil, errors.New("new error")).
					Once()
			},
		},
		{
			name: "given the validation because total amount > pending amount outstanding," +
				"when payment," +
				"then return error",
			args: args{
				paymentRequest: payReq,
			},
			wantErr: errorAmountShouldBeSame,
			mockFunc: func() {
				mockLoanRepo.
					On("FindLoans", mock.Anything, mock.Anything).
					Return(
						[]*repository.LoanEntity{
							{
								Status: "PENDING",
								Amount: decimal.NewFromFloat(float64(10)),
							},
							{
								Status: "PENDING",
								Amount: decimal.NewFromFloat(float64(13)),
							},
						}, nil).
					Once()
			},
		},
		{
			name: "given update loan is failed unknown error from database," +
				"when payment," +
				"then return error",
			args: args{
				paymentRequest: payReq,
			},
			wantErr: errorFromDatabase,
			mockFunc: func() {
				mockLoanRepo.
					On("FindLoans", mock.Anything, mock.Anything).
					Return(
						[]*repository.LoanEntity{
							{
								Status: "PENDING",
								Amount: decimal.NewFromFloat(float64(20)),
							},
							{
								Status: "PENDING",
								Amount: decimal.NewFromFloat(float64(5)),
							},
						}, nil).
					Once()

				mockLoanRepo.
					On("UpdateLoan", mock.Anything, mock.Anything).
					Return(errors.New("mock error")).
					Once()
			},
		},
		{
			name: "given update loan is success," +
				"when payment," +
				"then return nil",
			args: args{
				paymentRequest: payReq,
			},
			wantErr: nil,
			mockFunc: func() {
				mockLoanRepo.
					On("FindLoans", mock.Anything, mock.Anything).
					Return(
						[]*repository.LoanEntity{
							{
								Status: "PENDING",
								Amount: decimal.NewFromFloat(float64(20)),
							},
							{
								Status: "PENDING",
								Amount: decimal.NewFromFloat(float64(5)),
							},
						}, nil).
					Once()

				mockLoanRepo.
					On("UpdateLoan", mock.Anything, mock.Anything).
					Return(nil).
					Once()
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockFunc()
				l := NewLoanService(mockLoanRepo)

				err := l.Payment(context.Background(), tt.args.paymentRequest)
				assert.Equal(t, tt.wantErr, err)
			})
	}
}
