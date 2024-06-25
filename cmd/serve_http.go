package cmd

import (
	"context"
	"errors"
	"log"
	http2 "net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"gitlab.com/2024/Juni/amartha-billing-srv2/application/loan"
	"gitlab.com/2024/Juni/amartha-billing-srv2/configuration"
	"gitlab.com/2024/Juni/amartha-billing-srv2/delivery/http"
	"gitlab.com/2024/Juni/amartha-billing-srv2/infrastructure/repository"
)

var serveHttp = &cobra.Command{
	Use:   "serveHttp",
	Short: "Turn on amartha billing service HTTP Rest API",
	Long:  "Cobra CLI : turn on Billing service HTTP Rest API",
	Run: func(cmd *cobra.Command, args []string) {
		//init configuration and credential
		cfg, cre := fetchConfiguration()

		//init database master
		initDB := configuration.NewStoreImpl(cre)
		masterDB, err := initDB.InitDBMaster()

		if err != nil {
			panic(err)
		}

		loanRepository := repository.NewLoanRepository(masterDB)
		loanService := loan.NewLoanService(cfg, loanRepository)
		loanController := loan.NewLoanController(loanService)

		billingHttpServerAddress := cfg.GetString("server.address.http")
		router := mux.NewRouter()

		billingHandler := http.NewBillingHandler(cfg, loanController).BuildHttp(router)
		billingHttpServer := http2.Server{
			Addr:    billingHttpServerAddress,
			Handler: billingHandler,
		}

		go func() {
			log.Println(
				"[Billing Controller HTTP] server started. Listening on port",
				billingHttpServerAddress)

			if err := billingHttpServer.ListenAndServe(); err != nil &&
				!errors.Is(err, http2.ErrServerClosed) {
				log.Println("error on close http : " + err.Error())
			}
		}()

		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-done
		if err := billingHttpServer.Shutdown(context.Background()); err != nil {
			log.Println("[Billing Controller HTTP], shutdown has error", err)
		} else {
			log.Println("[Billing Controller HTTP] server stopped.")
		}
	},
}
