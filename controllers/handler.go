package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/eneipereira/go-order-service/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AppHandler func(w http.ResponseWriter, r *http.Request) error

func ErrorMiddleware(h AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err == nil {
			return
		}

		log.Printf("API Error: %v", err)

		switch {

		case errors.Is(err, model.ErrCustomerNotFound),
			errors.Is(err, model.ErrProductNotFound),
			errors.Is(err, model.ErrOrderNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)

		case errors.Is(err, model.ErrEmailAlreadyExists),
			errors.Is(err, model.ErrInsufficientStock),
			errors.Is(err, model.ErrInvalidStatusChange):
			http.Error(w, err.Error(), http.StatusConflict)

		case isValidationError(err):
			http.Error(w, err.Error(), http.StatusBadRequest)

		default:
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func getIDFromRequest(r *http.Request) (string, error) {
	idStr := chi.URLParam(r, "id")
	if _, err := uuid.Parse(idStr); err != nil {
		return "", model.ErrInvalidUUID
	}
	return idStr, nil
}

func getQueryParamAsInt(r *http.Request, paramName string, defaultValue int) (int, error) {
	paramStr := r.URL.Query().Get(paramName)
	if paramStr == "" {
		return defaultValue, nil
	}
	paramInt, err := strconv.Atoi(paramStr)
	if err != nil {
		return 0, model.ErrInvalidQueryParam
	}
	return paramInt, nil
}

func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func isValidationError(err error) bool {
	return isValidationErrorCust(err) || isValidationErrorProd(err) || errors.Is(err, model.ErrInvalidQueryParam) || errors.Is(err, model.ErrInvalidUUID)
}

func isValidationErrorCust(err error) bool {
	return errors.Is(err, model.ErrCustomerNameRequired) ||
		errors.Is(err, model.ErrCustomerNameTooShort) ||
		errors.Is(err, model.ErrCustomerNameTooLong) ||
		errors.Is(err, model.ErrCustomerEmailRequired) ||
		errors.Is(err, model.ErrCustomerEmailInvalid) ||
		errors.Is(err, model.ErrCustomerEmailTooLong) ||
		errors.Is(err, model.ErrCustomerPhoneRequired) ||
		errors.Is(err, model.ErrCustomerPhoneInvalid) ||
		errors.Is(err, model.ErrCustomerPhoneTooLong) ||
		errors.Is(err, model.ErrCustomerPasswordRequired) ||
		errors.Is(err, model.ErrCustomerPasswordTooShort) ||
		errors.Is(err, model.ErrCustomerPasswordTooLong) ||
		errors.Is(err, model.ErrInvalidQuantity) ||
		errors.Is(err, model.ErrInvalidUUID)
}

func isValidationErrorProd(err error) bool {
	return errors.Is(err, model.ErrProductNameRequired) ||
		errors.Is(err, model.ErrProductNameTooShort) ||
		errors.Is(err, model.ErrProductNameTooLong) ||
		errors.Is(err, model.ErrProductPriceRequired) ||
		errors.Is(err, model.ErrProductPriceTooLow) ||
		errors.Is(err, model.ErrProductStockRequired) ||
		errors.Is(err, model.ErrProductStockTooLow)
}
