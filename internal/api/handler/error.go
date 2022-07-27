package handler

import "github.com/kataras/iris/v12/context"

var NotFound = func(ctx context.Context) {
	_, _ = ctx.WriteString("404 NOT FOUND")
}

var InternalError = func(ctx context.Context) {
	_, _ = ctx.WriteString("500 Error Occurred")
}

var BadRequest = func(ctx context.Context) {
	_, _ = ctx.WriteString("400 Bad Request")
}
