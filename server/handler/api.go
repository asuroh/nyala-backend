package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"kriyapeople/pkg/jwe"
	"kriyapeople/pkg/jwt"
	"kriyapeople/pkg/str"
	"kriyapeople/usecase"

	"database/sql"
	"github.com/go-playground/universal-translator"
	validator "gopkg.in/go-playground/validator.v9"
)

// Handler ...
type Handler struct {
	ContractUC *usecase.ContractUC
	DB         *sql.DB
	EnvConfig  map[string]string
	Validate   *validator.Validate
	Translator ut.Translator
	Jwt        jwt.Credential
	Jwe        jwe.Credential
}

// Bind bind the API request payload (body) into request struct.
func (h Handler) Bind(r *http.Request, input interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&input)

	return err
}

// emptyJSONArr ...
func emptyJSONArr() []map[string]interface{} {
	return []map[string]interface{}{}
}

// SendSuccess send success into response with 200 http code.
func SendSuccess(w http.ResponseWriter, payload interface{}, meta interface{}) {
	RespondWithJSON(w, 200, 200, "Success", payload, meta)
}

// SendBadRequest send bad request into response with 400 http code.
func SendBadRequest(w http.ResponseWriter, message string) {
	RespondWithJSON(w, 400, 400, message, emptyJSONArr(), emptyJSONArr())
}

// SendRequestValidationError Send validation error response to consumers.
func (h Handler) SendRequestValidationError(w http.ResponseWriter, validationErrors validator.ValidationErrors) {
	errorResponse := map[string][]string{}
	errorTranslation := validationErrors.Translate(h.Translator)
	for _, err := range validationErrors {
		errKey := str.Underscore(err.StructField())
		errorResponse[errKey] = append(
			errorResponse[errKey],
			strings.Replace(errorTranslation[err.Namespace()], err.StructField(), "[]", -1),
		)
	}

	RespondWithJSON(w, 400, 405, "validation error", errorResponse, emptyJSONArr())
}

// RespondWithJSON write json response format
func RespondWithJSON(w http.ResponseWriter, httpCode int, statCode int, message string, payload interface{}, meta interface{}) {
	respPayload := map[string]interface{}{
		"stat_code": statCode,
		"stat_msg":  message,
		"meta":      meta,
		"data":      payload,
	}

	response, _ := json.Marshal(respPayload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	w.Write(response)
}

// requestIDFromContextInterface ...
func requestIDFromContextInterface(ctx context.Context, key string) (res map[string]interface{}) {
	if ctx.Value(key) != nil {
		res = ctx.Value(key).(map[string]interface{})
	}

	return res
}

// requestKeyFromContextInterface ...
func requestKeyFromContextInterface(ctx context.Context, key, objectKey string) (res string) {
	user := requestIDFromContextInterface(ctx, key)
	if user[objectKey] != nil {
		res = user[objectKey].(string)
	}

	return res
}
