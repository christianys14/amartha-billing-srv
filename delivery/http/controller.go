package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (b *billingHandler) routeBilling(r *mux.Router) {
	r.HandleFunc("/v1/customer/outstanding/{userID}", b.loanSrv.FindOutstanding).
		Methods(http.MethodGet)

	r.HandleFunc("/v1/customer/payment", b.loanSrv.Payment).
		Methods(http.MethodPost)
}
