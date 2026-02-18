package handlers

import (
	"fmt"
	"net/http"
	"wallet-service/pkg/helpers"
)

func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := helpers.Envelope{"error": message}

	err := helpers.WriteJSON(w, status, env, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {

	message := "the server encountered a problem and could not process your request"

	ErrorResponse(w, r, http.StatusInternalServerError, message)
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	ErrorResponse(w, r, http.StatusNotFound, message)
}

func MethodNotAllowResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this response", r.Method)
	ErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	ErrorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
