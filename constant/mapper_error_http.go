package constant

import (
	"net/http"
)

type BillingSrvHttpError int

const (
	Success BillingSrvHttpError = iota
	Validation
	DataNotFound
	GeneralError
	PaymentAmountShouldBeEquals
)

var HttpRc = map[BillingSrvHttpError]string{
	Success:                     "0000",
	Validation:                  "0001",
	DataNotFound:                "0002",
	PaymentAmountShouldBeEquals: "0003",
	GeneralError:                "9999",
}

var HttpRcDescription = map[BillingSrvHttpError]string{
	Success:                     "Successful",
	Validation:                  "one or more field should not be empty",
	DataNotFound:                "data is not exist",
	PaymentAmountShouldBeEquals: "amount of payment should be exact",
	GeneralError:                "General error",
}

var BillingCodeToHttpCode = map[string]int{
	"0000": http.StatusOK,
	"0001": http.StatusBadRequest,
	"0002": http.StatusNotFound,
	"0003": http.StatusBadRequest,
	"9999": http.StatusInternalServerError,
}
