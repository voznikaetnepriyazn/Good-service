package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"time"

	"github.com/valyala/fasthttp"
)

func RequestID(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		id := fmt.Sprintf("%d", time.Now().UnixNano())

		ctx.SetUserValue("requestID", id)

		ctx.Response.Header.Set("X-Request-ID", id)

		next(ctx)
	}
}

func Recoverer(log *slog.Logger) func(fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			defer func() {
				if err := recover(); err != nil {
					stackBuf := make([]byte, 4096)
					stackSize := runtime.Stack(stackBuf, false)
					stackTrace := string(stackBuf[:stackSize])

					log.Error("panic recovered",
						slog.Any("error", err),
						slog.String("stack_trace", stackTrace),
						slog.Any("method", ctx.Method()),
						slog.Any("path", ctx.Path()),
					)

					ctx.SetStatusCode(fasthttp.StatusInternalServerError)
					ctx.SetContentType("application/json")
					ctx.SetBodyString(`"error":"internal server error"}`)
				}
			}()

			next(ctx)
		}
	}
}

func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if reqID, ok := ctx.Value("requestID").(string); ok {
		return reqID
	}

	return ""
}
