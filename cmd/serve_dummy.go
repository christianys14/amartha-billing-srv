package cmd

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"

	"gitlab.com/2024/Juni/amartha-billing-srv2/common"
	"gitlab.com/2024/Juni/amartha-billing-srv2/configuration"
	"gitlab.com/2024/Juni/amartha-billing-srv2/infrastructure/repository"
)

var serveDummy = &cobra.Command{
	Use:   "serveDummy",
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

		numberOfCustomers := int(cfg.GetInt("custom.dummy.customers"))
		numberOfWeeks := int(cfg.GetInt("custom.weeks"))

		var loans []*repository.LoanEntity
		for i := 0; i < numberOfCustomers; i++ {
			userID := common.NewGenerate().Uuid()
			currentTime := time.Date(2023, 8, 23, 18, 58, 0, 0, time.UTC)

			amount := float64(5000000) / float64(numberOfWeeks)
			amountWithFee := amount + (amount * 0.1)

			le := repository.LoanEntity{
				UserID:    userID,
				Amount:    decimal.NewFromFloat(amountWithFee),
				CreatedAt: currentTime,
				Version:   0,
				UpdatedAt: currentTime,
			}

			for j := 0; j < numberOfWeeks; j++ {
				d := currentTime.AddDate(0, 0, j*7)
				status := func() string {
					//no outstanding
					if i == 1 {
						return "CLOSED"
					}

					//partially paid
					if j < 5 || i == 2 {
						return "PAID"
					}

					return "PENDING"
				}()

				loans = append(
					loans, &repository.LoanEntity{
						Status:    status,
						UserID:    userID,
						DueDate:   d,
						Amount:    le.Amount,
						CreatedAt: currentTime,
						Version:   0,
						UpdatedAt: currentTime,
					},
				)
			}
		}

		ctx := context.Background()
		ctx, cancelFunc := context.WithTimeout(ctx, time.Minute)
		defer cancelFunc()

		tx, errTx := masterDB.Begin()
		defer commitOrRollback(tx, &errTx)

		errInsert := loanRepository.SaveLoans(ctx, tx, loans...)
		if errInsert != nil {
			log.Println("failed during saveLoans -> ", errInsert)
			return
		}
	},
}

func commitOrRollback(tx *sql.Tx, sqlErr *error) {
	var err error
	if *sqlErr != nil {
		err = tx.Rollback()
		if err != nil {
			log.Println("failed when rollback -> ", err)
		}
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println("failed when commit -> ", err)
	}
}
