package handler

import (
	"github.com/kataras/iris/v12/context"
	"github.com/stellarisJAY/goim/internal/api/middleware"
)

func handleError(ctx context.Context, err error) {
	ctx.Values().Set(middleware.CtxKeyHandlerError, err)
	ctx.Next()
}
