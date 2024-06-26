package loan

import (
	"context"
	"errors"
	"log"
	"runtime/debug"

	"github.com/shopspring/decimal"

	"gitlab.com/2024/Juni/amartha-billing-srv2/infrastructure/repository"
)

var (
	errorValidation         = errors.New("validation request")
	errorFromDatabase       = errors.New("from database")
	errorDataNotExists      = errors.New("data is not exists")
	errorAmountShouldBeSame = errors.New("amount should be equals")
)

func (l *loanService) FetchOutstanding(
	ctx context.Context,
	uid string) (rsp *FetchOutstandingResponse, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Println("unidentified error (yet)", string(debug.Stack()))
			err = errorFromDatabase
			return
		}
	}()

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

	return l.identifyOutstanding(loans)
}

func (l *loanService) Payment(
	ctx context.Context,
	paymentRequest *PaymentRequest) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Println("unidentified error (yet)", string(debug.Stack()))
			err = errorFromDatabase
			return
		}
	}()

	if paymentRequest.UserID == "" || paymentRequest.Amount == 0 {
		return errorValidation
	}

	loans, errFindLoan := l.loanRepository.FindLoans(
		ctx, &repository.LoanEntity{
			Statuses: []string{"PENDING"},
			UserID:   paymentRequest.UserID,
			DueDate:  l.generate.Time(),
		},
	)

	if errors.Is(errFindLoan, repository.ErrorNoRows) {
		return errorDataNotExists
	}

	if errFindLoan != nil {
		return errorFromDatabase
	}

	return l.makePayment(ctx, paymentRequest, loans)
}

func (l *loanService) identifyOutstanding(
	loans []*repository.LoanEntity) (*FetchOutstandingResponse, error) {
	totalClosed, totalPending := 0, 0
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

func (l *loanService) makePayment(
	ctx context.Context,
	paymentRequest *PaymentRequest,
	loans []*repository.LoanEntity) error {
	amount := decimal.NewFromFloat(paymentRequest.Amount)
	totalAmount := decimal.NewFromFloat(float64(0))

	var loanIDs []uint64
	for _, loan := range loans {
		totalAmount = totalAmount.Add(loan.Amount)
		loanIDs = append(loanIDs, loan.ID)
	}

	if amount.LessThan(totalAmount) || amount.GreaterThan(totalAmount) {
		return errorAmountShouldBeSame
	}

	//integrate with 3rd party for debit the money customer
	//after the result of debit is success, then update the loan.
	//push notif (if any)
	//sent related marketing purposed, or any other activities.

	errUpdate := l.loanRepository.UpdateLoan(
		ctx, &repository.LoanEntityUpdate{
			IDs:    loanIDs,
			Status: "PAID",
		},
	)

	if errUpdate != nil {
		log.Println("failed update loan -> ", errUpdate)
		return errorFromDatabase
	}

	return nil
}
