package loan

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"

	"gitlab.com/2024/Juni/amartha-billing-srv2/infrastructure/repository"
)

var (
	errorValidation    = errors.New("validation request")
	errorFromDatabase  = errors.New("from database")
	errorDataNotExists = errors.New("data is not exists")
)

func (l *loanService) FetchOutstanding(
	ctx context.Context,
	uid string) (*FetchOutstandingResponse, error) {
	if uid == "" {
		return nil, errorValidation
	}

	loans, err := l.loanRepository.FindLoans(
		ctx, &repository.LoanEntity{
			Statuses: []string{"PENDING", "CLOSED"},
			UserID:   uid,
			DueDate:  l.generate.Time(),
		},
	)

	if errors.Is(err, repository.ErrorNoRows) {
		return nil, errorDataNotExists
	}

	if err != nil {
		return nil, errorFromDatabase
	}

	totalClosed := 0
	totalPending := 0
	pendingAmountOutstanding := decimal.NewFromFloat(float64(0))

	for _, val := range loans {
		if val.Status == "PENDING" {
			pendingAmountOutstanding = pendingAmountOutstanding.Add(val.Amount)
			totalPending += 1
		}

		if val.Status == "CLOSED" {
			totalClosed += 1
		}
	}

	//meaning : the customer already paid all the outstanding
	if totalClosed == len(loans) {
		return &FetchOutstandingResponse{
			RemainingOutstanding: decimal.NewFromFloat(float64(0)),
			IsDelinquent:         false,
		}, nil
	}

	if totalPending > 2 {
		return &FetchOutstandingResponse{
			RemainingOutstanding: pendingAmountOutstanding,
			IsDelinquent:         true,
		}, nil
	}

	return &FetchOutstandingResponse{
		RemainingOutstanding: pendingAmountOutstanding,
		IsDelinquent:         false,
	}, nil
}

func (l *loanService) Payment(
	ctx context.Context,
	paymentRequest PaymentRequest) error {
	//TODO implement me
	panic("implement me")
}
