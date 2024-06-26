package loan

import (
	"context"

	"github.com/shopspring/decimal"

	"gitlab.com/2024/Juni/amartha-billing-srv2/common"
	"gitlab.com/2024/Juni/amartha-billing-srv2/configuration"
	"gitlab.com/2024/Juni/amartha-billing-srv2/infrastructure/repository"
)

type (
	loanService struct {
		cfg            configuration.Configuration
		loanRepository repository.LoanRepository
		generate       common.Generate
	}

	FetchOutstandingResponse struct {
		RemainingOutstanding decimal.Decimal `json:"remaining_outstanding,omitempty"`
		IsDelinquent         bool            `json:"is_delinquent"`
	}

	PaymentRequest struct {
		UserID string  `json:"user_id,omitempty"`
		Amount float64 `json:"amount,omitempty"`
	}

	Service interface {
		FetchOutstanding(ctx context.Context, uid string) (*FetchOutstandingResponse, error)

		Payment(ctx context.Context, paymentRequest *PaymentRequest) error
	}
)

func NewLoanService(
	loanRepository repository.LoanRepository) Service {
	return &loanService{
		loanRepository: loanRepository,
		generate:       common.NewGenerate(),
	}
}
