package controllers

import (
	"errors"
	"net/url"
	"strings"

	"github.com/goravel/framework/contracts/http"

	"goravel/app/core/permission"
	"goravel/app/core/resource"
	apierrors "goravel/app/core/shared/errors"
	"goravel/app/core/shared/pagination"
	"goravel/app/core/shared/query"
	"goravel/app/core/shared/response"
)

type CoreResourceController struct {
	auth       *AuthController
	service    *resource.CoreResourceService
	permission string
}

func NewCoreResourceController(auth *AuthController, service *resource.CoreResourceService, permissionPrefix string) *CoreResourceController {
	return &CoreResourceController{auth: auth, service: service, permission: strings.TrimSpace(permissionPrefix)}
}

func (c *CoreResourceController) Index(ctx http.Context) http.Response {
	if response := c.authorize(ctx, "view"); response != nil {
		return response
	}
	listQuery, err := coreResourceListQuery(ctx)
	if err != nil {
		return c.error(ctx, 400, "INVALID_QUERY", "请求查询参数无效", map[string]string{"reason": err.Error()})
	}
	rows, total, err := c.service.List(ctx, listQuery)
	if err != nil {
		return c.storageError(ctx, err)
	}
	meta := response.Meta{RequestID: resourceRequestID(ctx)}
	paginationMeta := pagination.NewMeta(pagination.Query{Page: listQuery.Page, PerPage: listQuery.PerPage}, total)
	meta.Pagination = &paginationMeta
	return ctx.Response().Success().Json(response.Success(rows, meta))
}

func (c *CoreResourceController) Show(ctx http.Context) http.Response {
	if response := c.authorize(ctx, "view"); response != nil {
		return response
	}
	row, err := c.service.Get(ctx, strings.TrimSpace(ctx.Request().Route("id")))
	if err != nil {
		return c.domainError(ctx, err)
	}
	return c.success(ctx, row)
}

func (c *CoreResourceController) Store(ctx http.Context) http.Response {
	if response := c.authorize(ctx, "create"); response != nil {
		return response
	}
	var input resource.CoreResourceRecord
	if err := ctx.Request().Bind(&input); err != nil {
		return c.error(ctx, 400, "INVALID_REQUEST", "资源请求格式无效", nil)
	}
	row, err := c.service.Create(ctx, input)
	if err != nil {
		return c.domainError(ctx, err)
	}
	return ctx.Response().Status(201).Json(response.Success(row, response.Meta{RequestID: resourceRequestID(ctx)}))
}

func (c *CoreResourceController) Update(ctx http.Context) http.Response {
	if response := c.authorize(ctx, "update"); response != nil {
		return response
	}
	var input resource.CoreResourceRecord
	if err := ctx.Request().Bind(&input); err != nil {
		return c.error(ctx, 400, "INVALID_REQUEST", "资源请求格式无效", nil)
	}
	row, err := c.service.Update(ctx, strings.TrimSpace(ctx.Request().Route("id")), input)
	if err != nil {
		return c.domainError(ctx, err)
	}
	return c.success(ctx, row)
}

func (c *CoreResourceController) Destroy(ctx http.Context) http.Response {
	if response := c.authorize(ctx, "delete"); response != nil {
		return response
	}
	if err := c.service.Delete(ctx, strings.TrimSpace(ctx.Request().Route("id"))); err != nil {
		return c.domainError(ctx, err)
	}
	return c.success(ctx, map[string]bool{"deleted": true})
}

func (c *CoreResourceController) BulkDestroy(ctx http.Context) http.Response {
	if response := c.authorize(ctx, "delete"); response != nil {
		return response
	}
	var input struct {
		IDs []string `json:"ids"`
	}
	if err := ctx.Request().Bind(&input); err != nil {
		return c.error(ctx, 400, "INVALID_REQUEST", "资源请求格式无效", nil)
	}
	if err := c.service.BulkDelete(ctx, input.IDs); err != nil {
		return c.domainError(ctx, err)
	}
	return c.success(ctx, map[string]bool{"deleted": true})
}

func (c *CoreResourceController) authorize(ctx http.Context, action string) http.Response {
	if c.service == nil || c.permission == "" {
		return c.error(ctx, 503, "RESOURCE_NOT_CONFIGURED", "核心资源尚未配置", nil)
	}
	if c.auth == nil {
		return c.error(ctx, 401, "UNAUTHENTICATED", "请先登录", nil)
	}
	currentUser, err := c.auth.CurrentUser(ctx)
	if err != nil {
		return c.error(ctx, 401, "UNAUTHENTICATED", "请先登录", nil)
	}
	if permission.NewAuthorizer().Require(currentUser.Permissions, c.permission+"."+action) != nil {
		return c.error(ctx, 403, "FORBIDDEN", "没有执行该操作的权限", map[string]string{"permission": c.permission + "." + action})
	}
	return nil
}

func coreResourceListQuery(ctx http.Context) (resource.DemoResourceListQuery, error) {
	values := url.Values{}
	for key, value := range ctx.Request().Queries() {
		values.Set(key, value)
	}
	parsed, err := query.Parse(values)
	if err != nil {
		return resource.DemoResourceListQuery{}, err
	}
	return resource.DemoResourceListQuery{
		Page: parsed.Pagination.Page, PerPage: parsed.Pagination.PerPage,
		Search: parsed.Search.Term, Filters: parsed.Filter.Values,
		SortField: parsed.Sort.Field, SortDesc: parsed.Sort.Desc,
	}, nil
}

func (c *CoreResourceController) domainError(ctx http.Context, err error) http.Response {
	switch {
	case errors.Is(err, resource.ErrInvalidResource):
		return c.error(ctx, 400, "INVALID_RESOURCE", "资源数据无效", nil)
	case errors.Is(err, resource.ErrResourceNotFound):
		return c.error(ctx, 404, "RESOURCE_NOT_FOUND", "资源不存在", nil)
	case errors.Is(err, resource.ErrResourceExists):
		return c.error(ctx, 409, "RESOURCE_EXISTS", "资源已存在", nil)
	default:
		return c.storageError(ctx, err)
	}
}

func (c *CoreResourceController) storageError(ctx http.Context, err error) http.Response {
	return c.error(ctx, 503, "RESOURCE_STORAGE_UNAVAILABLE", "资源存储暂不可用", map[string]string{"reason": err.Error()})
}

func (c *CoreResourceController) success(ctx http.Context, data any) http.Response {
	return ctx.Response().Success().Json(response.Success(data, response.Meta{RequestID: resourceRequestID(ctx)}))
}

func (c *CoreResourceController) error(ctx http.Context, status int, code, message string, details any) http.Response {
	return ctx.Response().Status(status).Json(apierrors.New(code, message, details))
}
