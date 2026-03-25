package logger

import (
	"log/slog"
	"time"

	"github.com/voznikaetnepriyazn/Good-service/internal/http-server/middleware"

	"github.com/valyala/fasthttp"
)

func New(log *slog.Logger) func(fasthttp.RequestHandler) fasthttp.RequestHandler {

	log = log.With(
		slog.String("component", "middleware/logger"),
	)

	log.Info("logger middleware enabled")

	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			entry := log.With(
				slog.Any("method", ctx.Method()),
				slog.Any("path", ctx.Path()),
				slog.String("remote_addr", ctx.RemoteAddr().String()),
				slog.Any("user_agent", ctx.Request.Header.Peek("User-Agent")),
			)

			ctx.SetUserValue("logger", entry)

			start := time.Now()

			entry.Info("request completed",
				slog.Int("status", ctx.Response.StatusCode()),
				slog.Int("bytes", len(ctx.Response.Body())),
				slog.Duration("duration", time.Since(start)),
			)

			next(ctx)
		}
	}
}

func LogQuery(log *slog.Logger, op string, next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		ctx.SetUserValue("logger", log)
		next(ctx)
	}
}

func FromCtx(ctx *fasthttp.RequestCtx) *slog.Logger {
	if v := ctx.UserValue("logger"); v != nil {
		if log, ok := v.(*slog.Logger); ok {
			return log
		}
	}
	return slog.Default()
}
