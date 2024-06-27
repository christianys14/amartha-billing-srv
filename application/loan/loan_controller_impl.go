package loan

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"gitlab.com/2024/Juni/amartha-billing-srv2/common"
	"gitlab.com/2024/Juni/amartha-billing-srv2/constant"
)

func (l *loanController) FindOutstanding(
	writer http.ResponseWriter,
	req *http.Request) {
	query := mux.Vars(req)
	userID, errEscape := escapeSpecialCharacter(query["userID"])

	if errEscape != nil {
		common.ToErrorResponse(
			writer,
			constant.HttpRc[constant.Validation],
			constant.HttpRcDescription[constant.Validation],
		)
		return
	}

	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, time.Minute)
	defer cancelFunc()

	result, err := l.srv.FetchOutstanding(ctx, userID)
	if err != nil {
		if errors.Is(err, errorValidation) {
			common.ToErrorResponse(
				writer,
				constant.HttpRc[constant.Validation],
				constant.HttpRcDescription[constant.Validation],
			)
			return
		}

		if errors.Is(err, errorDataNotExists) {
			common.ToErrorResponse(
				writer,
				constant.HttpRc[constant.DataNotFound],
				constant.HttpRcDescription[constant.DataNotFound],
			)
			return
		}

		common.ToErrorResponse(
			writer,
			constant.HttpRc[constant.GeneralError],
			constant.HttpRcDescription[constant.GeneralError],
		)
		return
	}

	common.ToSuccessResponse(writer, nil, result)
}

func (l *loanController) Payment(
	writer http.ResponseWriter,
	req *http.Request) {
	var paymentRequest PaymentRequest
	err := decodeJSONBody(writer, req, &paymentRequest)

	if err != nil {
		log.Println("validation decode json body -> ", err)

		common.ToErrorResponse(
			writer,
			constant.HttpRc[constant.Validation],
			constant.HttpRcDescription[constant.Validation],
		)
		return
	}

	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, time.Minute)
	defer cancelFunc()

	errPayment := l.srv.Payment(ctx, &paymentRequest)

	if errPayment != nil {
		if errors.Is(errPayment, errorValidation) {
			common.ToErrorResponse(
				writer,
				constant.HttpRc[constant.Validation],
				constant.HttpRcDescription[constant.Validation],
			)
			return
		}

		if errors.Is(errPayment, errorAmountShouldBeSame) {
			common.ToErrorResponse(
				writer,
				constant.HttpRc[constant.PaymentAmountShouldBeEquals],
				constant.HttpRcDescription[constant.PaymentAmountShouldBeEquals],
			)
			return
		}

		if errors.Is(errPayment, errorNoPendingOutstanding) {
			common.ToErrorResponse(
				writer,
				constant.HttpRc[constant.ZeroOutstanding],
				constant.HttpRcDescription[constant.ZeroOutstanding],
			)
			return
		}

		common.ToErrorResponse(
			writer,
			constant.HttpRc[constant.GeneralError],
			constant.HttpRcDescription[constant.GeneralError],
		)
		return
	}

	common.ToSuccessResponse(writer, nil, nil)
}

type malformedRequestError struct {
	message string
}

func (m malformedRequestError) Error() string {
	return m.message
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	contentType := r.Header.Get("Content-Type")

	if contentType != "application/json" {
		message := "Content-Type header is not application/json"
		return &malformedRequestError{message: message}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			message := fmt.Sprintf(
				"Request body contains badly-formed JSON (at position %d)",
				syntaxError.Offset)
			return &malformedRequestError{message: message}

		case errors.Is(err, io.ErrUnexpectedEOF):
			message := "Request body contains badly-formed JSON"
			return &malformedRequestError{message: message}

		case errors.As(err, &unmarshalTypeError):
			message := fmt.Sprintf(
				"Request body contains an invalid value for the %q field (at position %d)",
				unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequestError{message: message}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			message := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequestError{message: message}

		case errors.Is(err, io.EOF):
			message := "Request body must not be empty"
			return &malformedRequestError{message: message}

		case err.Error() == "http: request body too large":
			message := "Request body must not be larger than 1MB"
			return &malformedRequestError{message: message}

		default:
			return &malformedRequestError{message: err.Error()}
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		message := "Request body must only contain a single JSON object"
		return &malformedRequestError{message: message}
	}

	return nil
}

func escapeSpecialCharacter(string string) (string, error) {
	reg, err := regexp.Compile(`[!?;{}|<>%'=]`)
	if err != nil {
		return string, err
	}

	newString := reg.ReplaceAllString(string, "")
	return newString, nil
}
