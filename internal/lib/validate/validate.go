package valid

import (
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/valyala/fasthttp"
)

var validate = validator.New()

func Validate(ctx *fasthttp.RequestCtx, req interface{}, log *slog.Logger) bool {
	if err := validate.Struct(req); err != nil {
		var validateErr validator.ValidationErrors
		if !errors.As(err, &validateErr) {
			log.Error("unknown validation error", slog.Any("error", err))

			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.SetContentType("application/json")
			ctx.SetBodyString(`error":"failed to validate request"}`)

			return false
		}

		log.Warn("validation failed", slog.Any("errors", FormatValidationError(validateErr)))

		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(map[string]any{
			"error":   "validation failed",
			"details": FormatValidationError(validateErr),
		})
		ctx.SetBody(body)

		return false
	}
	return true
}

func FormatValidationError(err validator.ValidationErrors) map[string]string {
	errors := make(map[string]string)
	for _, e := range err {
		errors[e.Field()] = e.Error()
	}
	return errors
}
