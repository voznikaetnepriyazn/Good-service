package uuidparam

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

func UUIDParam(paramName string, log *slog.Logger, next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		alias := string(ctx.QueryArgs().Peek("paramName"))
		if alias == "" {
			slog.Warn("id is empty")

			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.SetContentType("application/json")
			ctx.SetBodyString(`"id is empty"`)
		}

		id, err := uuid.Parse(alias)
		if err != nil {
			slog.Warn("invalid uuid format", slog.String("id", alias), slog.Any("error", err))

			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.SetContentType("application/json")
			ctx.SetBodyString(`"invalid id format"`)
		}

		ctx.SetUserValue("uuid_"+paramName, id)
		next(ctx)
	}
}

func UUIDFromCtx(ctx *fasthttp.RequestCtx, paramName string) (uuid.UUID, bool) {
	if v := ctx.UserValue("uuid_" + paramName); v != nil {
		if id, ok := v.(uuid.UUID); ok {
			return id, true
		}
	}
	return uuid.Nil, false
}
