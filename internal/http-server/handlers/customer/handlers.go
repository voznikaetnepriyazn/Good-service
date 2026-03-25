package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/voznikaetnepriyazn/Good-service/internal/http-server/middleware/logger"
	uuidparam "github.com/voznikaetnepriyazn/Good-service/internal/http-server/middleware/uuid"
	response "github.com/voznikaetnepriyazn/Good-service/internal/lib/api/response"
	decodejson "github.com/voznikaetnepriyazn/Good-service/internal/lib/decode"
	"github.com/voznikaetnepriyazn/Good-service/internal/lib/logger/sl"
	valid "github.com/voznikaetnepriyazn/Good-service/internal/lib/validate"
	"github.com/voznikaetnepriyazn/Good-service/internal/models/good"
	"github.com/voznikaetnepriyazn/Good-service/internal/storage"

	"github.com/valyala/fasthttp"
)

type Response struct {
	URL string `json:"url" validate:"required, url"`
}

type RequestFullStruct struct {
	Good good.Good
}

type Request struct {
	response.Response
	URL string `json:"url" validate:"required, url"`
}

type Crud interface {
	NewAdd(log *slog.Logger, adder storage.GoodService) fasthttp.RequestHandler
	NewDelete(log *slog.Logger, deleter storage.GoodService) fasthttp.RequestHandler
	NewGetAll(log *slog.Logger, get storage.GoodService) fasthttp.RequestHandler
	NewGetById(log *slog.Logger, get storage.GoodService) fasthttp.RequestHandler
	NewUpdate(log *slog.Logger, update storage.GoodService) fasthttp.RequestHandler
	NewIsOrderCreated(log *slog.Logger, ord storage.GoodService) fasthttp.RequestHandler
}

func NewAdd(log *slog.Logger, adder storage.GoodService) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		const op = "handlers.url.add.New"

		log := logger.FromCtx(ctx)

		log.Info("handling request")

		var req RequestFullStruct

		if !decodejson.DecodeJSON(ctx, &req, log) {
			log.Error("can not decode json")
		}

		if !valid.Validate(ctx, &req, log) {
			log.Error("can not validate")
		}

		//проверка на уже существующее значение
		id, err := adder.AddURL(req.Good)
		if errors.Is(err, storage.ErrUrlExist) {
			log.Warn("url already exists", slog.Any("url", req.Good))

			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.SetContentType("application/json")
			ctx.SetBodyString(`"url already exist"`)
		}

		//прочие ошибки
		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.SetContentType("application/json")
			ctx.SetBodyString(`"failed to add url"`)
		}

		log.Info("url added", slog.Any("id", id))

		responseOK(ctx)
	}
}

func responseOK(ctx *fasthttp.RequestCtx) error {
	ctx.SetStatusCode(fasthttp.StatusAccepted)
	ctx.SetContentType("application/json")
	ctx.SetBodyString(`"status: ok"`)

	return nil
}

func NewDelete(log *slog.Logger, deleter storage.GoodService) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		const op = "handlers.url.delete.New"

		log := logger.FromCtx(ctx)

		log.Info("handling request")

		uuidParam, ok := uuidparam.UUIDFromCtx(ctx, "id")
		if ok {
			log.Info("took id from context")
		}

		err := deleter.DeleteURL(uuidParam)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("url not found", "id", uuidParam)

			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.SetContentType("application/json")
			ctx.SetBodyString(`"not found"`)
		}

		if err != nil {
			log.Error("failed to delete url", sl.Err(err))

			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.SetContentType("application/json")
			ctx.SetBodyString(`"internal server error"`)
		}

		log.Info("deleted url", slog.Any("deleted", uuidParam))

		responseOK(ctx)
	}
}

func NewGetAll(log *slog.Logger, get storage.GoodService) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		const op = "handlers.url.getById.New"

		log := logger.FromCtx(ctx)

		log.Info("handling request")

		resURL, err := get.GetAllURL()
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("urls not found")

			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.SetContentType("application/json")
			ctx.SetBodyString(`"not found"`)
		}

		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.SetContentType("application/json")
		}

		log.Info("got urls", slog.Any("urls", resURL))

		ctx.SetStatusCode(fasthttp.StatusAccepted)
		ctx.SetContentType("application/json")

		body, err := json.Marshal(map[string]interface{}{
			"urls": resURL,
		})
		if err != nil {
			ctx.Error("Marshal error", fasthttp.StatusInternalServerError)
			return
		}
		ctx.SetBody(body)
	}
}

func NewGetById(log *slog.Logger, get storage.GoodService) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		const op = "handlers.url.getById.New"

		log := logger.FromCtx(ctx)

		log.Info("handling request")

		uuidParam, ok := uuidparam.UUIDFromCtx(ctx, "id")
		if ok {
			log.Info("took id from context")
		}
		resURL, err := get.GetByIdURL(uuidParam)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("url not found", slog.String("id", uuidParam.String()))

			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.SetContentType("application/json")
			ctx.SetBodyString(`"not found"`)
		}

		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.SetContentType("application/json")
			ctx.SetBodyString(`"internal server errror"`)
		}

		log.Info("got url", slog.Any("url", resURL))

		ctx.SetStatusCode(fasthttp.StatusAccepted)
		ctx.SetContentType("application/json")

		body, err := json.Marshal(map[string]any{
			"urls": resURL,
		})
		if err != nil {
			ctx.Error("Marshal error", fasthttp.StatusInternalServerError)
			return
		}
		ctx.SetBody(body)
	}
}

func NewUpdate(log *slog.Logger, update storage.GoodService) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		const op = "handlers.url.update.New"

		log := logger.FromCtx(ctx)

		log.Info("handling request")

		var req RequestFullStruct

		if !decodejson.DecodeJSON(ctx, &req, log) {
			log.Error("can not decode json")
		}

		if !valid.Validate(ctx, &req, log) {
			log.Error("can not validate")
		}

		err := update.UpdateURL(req.Good)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("url not found", "id", req)

			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.SetContentType("application/json")
			ctx.SetBodyString(`"not found"`)
		}

		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.SetContentType("application/json")
			ctx.SetBodyString(`"internal server error"`)
		}

		log.Info("updated url", slog.Any("url", req))

		responseOK(ctx)
	}
}

func NewGetListOfGoodsByType(log *slog.Logger, ord storage.GoodService) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		const op = "handlers.url.getListOfGoodsByType.New"

		log := logger.FromCtx(ctx)

		log.Info("handling request")

		uuidParam, ok := uuidparam.UUIDFromCtx(ctx, "id")
		if ok {
			log.Info("took id from context")
		}

		resURL, err := ord.GetListOfGoodsByType(uuidParam)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("url not found", slog.String("id", uuidParam.String()))

			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.SetContentType("application/json")
			ctx.SetBodyString(`"not found"`)
		}

		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.SetContentType("application/json")
			ctx.SetBodyString(`"internal server errror"`)
		}

		log.Info("got urls", slog.Any("url", resURL))

		responseOK(ctx)
	}
}

func NewIsAvaliableForOrder(log *slog.Logger, ord storage.GoodService) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		const op = "handlers.url.isAvaliableForOrder.New"

		log := logger.FromCtx(ctx)

		log.Info("handling request")

		uuidParam, ok := uuidparam.UUIDFromCtx(ctx, "id")
		if ok {
			log.Info("took id from context")
		}

		resURL, err := ord.IsAvaliableForOrder(uuidParam)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("url not found", slog.String("id", uuidParam.String()))

			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.SetContentType("application/json")
			ctx.SetBodyString(`"not found id"`)
		}

		if err != nil {
			log.Error("failed to found url", sl.Err(err))

			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.SetContentType("application/json")
			ctx.SetBodyString(`"internal server errror"`)
		}

		log.Info("avaliable for order", slog.Any("url", resURL))

		responseOK(ctx)
	}
}

func NewRestOfGood(log *slog.Logger, ord storage.GoodService) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		const op = "handlers.url.restOfGood.New"

		log := logger.FromCtx(ctx)

		log.Info("handling request")

		uuidParam, ok := uuidparam.UUIDFromCtx(ctx, "id")
		if ok {
			log.Info("took id from context")
		}

		resURL, err := ord.RestOfGood(uuidParam)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("url not found", slog.String("id", uuidParam.String()))

			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.SetContentType("application/json")
			ctx.SetBodyString(`"not found"`)
		}

		if err != nil {
			log.Error("failed to found url", sl.Err(err))

			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.SetContentType("application/json")
			ctx.SetBodyString(`"internal server errror"`)
		}

		log.Info("rest", slog.Any("url", resURL))

		responseOK(ctx)
	}
}
