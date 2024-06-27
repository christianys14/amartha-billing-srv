package loan

import (
	"net/http"
)

type (
	loanController struct {
		srv Service
	}

	Controller interface {
		FindOutstanding(writer http.ResponseWriter, req *http.Request)

		Payment(writer http.ResponseWriter, req *http.Request)
	}
)

func NewLoanController(srv Service) Controller {
	return &loanController{
		srv: srv,
	}
}
