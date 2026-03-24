package decodejson

import (
	"encoding/json"
	"log/slog"

	"github.com/voznikaetnepriyazn/Good-service/internal/lib/logger/sl"

	"github.com/valyala/fasthttp"
)

func DecodeJSON(ctx *fasthttp.RequestCtx, req interface{}, log *slog.Logger) bool {
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		slog.Error("failed to decode request body", sl.Err(err))

		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetContentType("application/json")
		ctx.SetBodyString(`error":"failed to decode request"}`)

		slog.Error("failed to send error response", sl.Err(err))

		return false
	}

	slog.Info("request body decoded", slog.Any("request", req))
	return true
}
