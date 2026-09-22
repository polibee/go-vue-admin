package middleware

import (
	"strconv"

	httpcontract "github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	auditservices "goravel/app/services/audit"
)

type httpAuditMiddleware struct{}

func (httpAuditMiddleware) Signature() string {
	return "admin:http-audit"
}

func (httpAuditMiddleware) Handle(ctx httpcontract.Context) {
	if !auditservices.ShouldAuditHTTP(ctx.Request().Method(), ctx.Request().Path()) {
		ctx.Request().Next()
		return
	}

	input := auditservices.HTTPAuditInput{
		Method:      ctx.Request().Method(),
		Path:        ctx.Request().Path(),
		Query:       ctx.Request().Queries(),
		Route:       routeParams(ctx),
		RequestBody: ctx.Request().All(),
	}
	ctx.Request().Next()

	origin := ctx.Response().Origin()
	input.Status = origin.Status()
	input.ContentType = origin.Header().Get("Content-Type")
	if origin.Body() != nil {
		body := origin.Body().Bytes()
		input.ResponseBodyBytes = len(body)
		input.ResponseBody, input.ResponseTruncated = auditservices.BoundedBytes(body)
	}

	var userID *uint
	if identity, err := facades.Auth(ctx).ID(); err == nil {
		if parsed, parseErr := strconv.ParseUint(identity, 10, 32); parseErr == nil && parsed > 0 {
			value := uint(parsed)
			userID = &value
		}
	}
	if err := auditservices.NewAuditService().RecordHTTP(userID, input); err != nil {
		facades.Log().Errorf("http audit record failed method=%s path=%s error=%v", input.Method, input.Path, err)
	}
}

func routeParams(ctx httpcontract.Context) map[string]string {
	params := make(map[string]string)
	for _, name := range []string{"resource", "id", "action", "relation", "key"} {
		if value := ctx.Request().Route(name); value != "" {
			params[name] = value
		}
	}
	return params
}

func HTTPAudit() httpcontract.Middleware {
	return httpAuditMiddleware{}
}
