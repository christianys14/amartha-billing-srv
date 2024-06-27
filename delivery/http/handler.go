package http

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"gitlab.com/2024/Juni/amartha-billing-srv2/application/loan"
	"gitlab.com/2024/Juni/amartha-billing-srv2/configuration"
)

type billingHandler struct {
	configuration configuration.Configuration
	loanSrv       loan.Controller
}

func NewBillingHandler(
	configuration configuration.Configuration,
	loanSrv loan.Controller) *billingHandler {
	return &billingHandler{
		configuration: configuration,
		loanSrv:       loanSrv,
	}
}

func (b *billingHandler) showVersion() {
	version := b.configuration.GetString("app.billing.version")
	log.Println("show-billing-version -> ", version)
}

func (b *billingHandler) BuildHttp(router *mux.Router) http.Handler {
	b.showVersion()

	b.routeBilling(router)

	return router
}
