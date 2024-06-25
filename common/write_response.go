package common

import (
	"encoding/json"
	"log"
	"net/http"

	"gitlab.com/2024/Juni/amartha-billing-srv2/constant"
)

const (
	contentType = "Content-type"
	application = "application/json"
)

type BillingResponse struct {
	Rc         string      `json:"rc,omitempty"`
	Message    string      `json:"message,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

func NewBillingResponse(
	rc string,
	message string,
	pagination interface{},
	data interface{}) *BillingResponse {
	return &BillingResponse{
		Rc:         rc,
		Message:    message,
		Pagination: pagination,
		Data:       data,
	}
}

func responseWrite(rw http.ResponseWriter, data interface{}, statusCode int) {
	responseByte, err := json.Marshal(data)
	if err != nil {
		log.Println("error during encode responseWrite", err)
	}

	rw.WriteHeader(statusCode)
	_, err = rw.Write(responseByte)
}

func ToSuccessResponse(writer http.ResponseWriter, pagination interface{}, data interface{}) {
	rc := constant.HttpRc[constant.Success]
	rcDesc := constant.HttpRcDescription[constant.Success]
	httpRes := constant.BillingCodeToHttpCode[rc]

	responseWrite(
		writer,
		NewBillingResponse(rc, rcDesc, pagination, data),
		httpRes,
	)
}

func ToErrorResponse(writer http.ResponseWriter, rc, rcDesc string) {
	httpRes := constant.BillingCodeToHttpCode[rc]

	responseWrite(
		writer,
		NewBillingResponse(rc, rcDesc, nil, nil),
		httpRes,
	)
}
