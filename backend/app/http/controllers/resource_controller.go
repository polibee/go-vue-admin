package controllers

import (
	"errors"
	"net/url"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"goravel/app/core/permission"
	"goravel/app/core/resource"
	apierrors "goravel/app/core/shared/errors"
	"goravel/app/core/shared/pagination"
	"goravel/app/core/shared/query"
	"goravel/app/core/shared/response"
)

const demoResourcePermission = "dashboard.view"

type ResourceController struct {
	auth      *AuthController
	service   *resource.DemoResourceService
	initError error
}

func NewResourceController(auth *AuthController, service *resource.DemoResourceService) *ResourceController {
	return &ResourceController{auth: auth, service: service}
}

func NewConfiguredResourceController(auth *AuthController) *ResourceController {
	mode, err := resource.ParseResourceProviderMode(facades.Config().GetString("resource.provider", "memory"))
	if err != nil {
		return &ResourceController{auth: auth, initError: err}
	}

	var repository resource.DemoResourceRepository
	switch mode {
	case resource.ResourceProviderMemory:
		repository = resource.NewMemoryDemoResourceRepository(
			resource.DemoResource{ID: "demo-1", Name: "资源引擎示例", Status: "active", Owner: "Platform Admin"},
			resource.DemoResource{ID: "demo-2", Name: "可编辑记录", Status: "draft", Owner: "Platform Admin"},
		)
	case resource.ResourceProviderMySQL:
		repository = resource.NewMySQLDemoResourceRepository()
	default:
		return &ResourceController{auth: auth, initError: errors.New("resource provider is not configured")}
	}
	return NewResourceController(auth, resource.NewDemoResourceService(repository))
}

type demoResourceCreateRequest struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Owner  string `json:"owner"`
}

type demoResourceUpdateRequest struct {
	Name   *string `json:"name"`
	Status *string `json:"status"`
	Owner  *string `json:"owner"`
}

type demoResourceBulkDeleteRequest struct {
	IDs []string `json:"ids"`
}

func (c *ResourceController) Index(ctx http.Context) http.Response {
	if response := c.authorize(ctx); response != nil {
		return response
	}
	listQuery, err := c.listQuery(ctx)
	if err != nil {
		return c.error(ctx, 400, "INVALID_QUERY", "请求查询参数无效", map[string]string{"reason": err.Error()})
	}
	rows, total, err := c.service.List(ctx, listQuery)
	if err != nil {
		return c.storageError(ctx, err)
	}
	meta := response.Meta{RequestID: resourceRequestID(ctx)}
	pagination := responsePagination(listQuery, total)
	meta.Pagination = &pagination
	return ctx.Response().Success().Json(response.Success(rows, meta))
}

func (c *ResourceController) Show(ctx http.Context) http.Response {
	if response := c.authorize(ctx); response != nil {
		return response
	}
	row, err := c.service.Get(ctx, strings.TrimSpace(ctx.Request().Route("id")))
	if err != nil {
		return c.domainError(ctx, err)
	}
	return c.success(ctx, row)
}

func (c *ResourceController) Store(ctx http.Context) http.Response {
	if response := c.authorize(ctx); response != nil {
		return response
	}
	var input demoResourceCreateRequest
	if err := ctx.Request().Bind(&input); err != nil {
		return c.error(ctx, 400, "INVALID_REQUEST", "资源请求格式无效", nil)
	}
	row, err := c.service.Create(ctx, resource.DemoResource{
		ID: input.ID, Name: input.Name, Status: input.Status, Owner: input.Owner,
	})
	if err != nil {
		return c.domainError(ctx, err)
	}
	return ctx.Response().Status(201).Json(response.Success(row, response.Meta{RequestID: resourceRequestID(ctx)}))
}

func (c *ResourceController) Update(ctx http.Context) http.Response {
	if response := c.authorize(ctx); response != nil {
		return response
	}
	var input demoResourceUpdateRequest
	if err := ctx.Request().Bind(&input); err != nil {
		return c.error(ctx, 400, "INVALID_REQUEST", "资源请求格式无效", nil)
	}
	row, err := c.service.Update(ctx, strings.TrimSpace(ctx.Request().Route("id")), resource.DemoResourceUpdate{
		Name: input.Name, Status: input.Status, Owner: input.Owner,
	})
	if err != nil {
		return c.domainError(ctx, err)
	}
	return c.success(ctx, row)
}

func (c *ResourceController) Destroy(ctx http.Context) http.Response {
	if response := c.authorize(ctx); response != nil {
		return response
	}
	if err := c.service.Delete(ctx, strings.TrimSpace(ctx.Request().Route("id"))); err != nil {
		return c.domainError(ctx, err)
	}
	return c.success(ctx, map[string]bool{"deleted": true})
}

func (c *ResourceController) BulkDestroy(ctx http.Context) http.Response {
	if response := c.authorize(ctx); response != nil {
		return response
	}
	var input demoResourceBulkDeleteRequest
	if err := ctx.Request().Bind(&input); err != nil {
		return c.error(ctx, 400, "INVALID_REQUEST", "批量删除请求格式无效", nil)
	}
	if err := c.service.BulkDelete(ctx, input.IDs); err != nil {
		return c.domainError(ctx, err)
	}
	return c.success(ctx, map[string]bool{"deleted": true})
}

func (c *ResourceController) authorize(ctx http.Context) http.Response {
	if c.initError != nil || c.service == nil {
		return c.error(ctx, 503, "RESOURCE_PROVIDER_NOT_CONFIGURED", "资源 Provider 尚未配置", nil)
	}
	if c.auth == nil {
		return c.error(ctx, 401, "UNAUTHENTICATED", "请先登录", nil)
	}
	user, err := c.auth.CurrentUser(ctx)
	if err != nil {
		return c.error(ctx, 401, "UNAUTHENTICATED", "请先登录", nil)
	}
	if permission.NewAuthorizer().Require(user.Permissions, demoResourcePermission) != nil {
		return c.error(ctx, 403, "FORBIDDEN", "没有执行该操作的权限", map[string]string{"permission": demoResourcePermission})
	}
	return nil
}

func (c *ResourceController) listQuery(ctx http.Context) (resource.DemoResourceListQuery, error) {
	values := url.Values{}
	for key, value := range ctx.Request().Queries() {
		values.Set(key, value)
	}
	parsed, err := query.Parse(values)
	if err != nil {
		return resource.DemoResourceListQuery{}, err
	}
	return resource.DemoResourceListQuery{
		Page:      parsed.Pagination.Page,
		PerPage:   parsed.Pagination.PerPage,
		Search:    parsed.Search.Term,
		Filters:   parsed.Filter.Values,
		SortField: parsed.Sort.Field,
		SortDesc:  parsed.Sort.Desc,
	}, nil
}

func responsePagination(query resource.DemoResourceListQuery, total int) pagination.Meta {
	return pagination.NewMeta(pagination.Query{Page: query.Page, PerPage: query.PerPage}, total)
}

func (c *ResourceController) domainError(ctx http.Context, err error) http.Response {
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

func (c *ResourceController) storageError(ctx http.Context, err error) http.Response {
	return c.error(ctx, 503, "RESOURCE_STORAGE_UNAVAILABLE", "资源存储暂不可用", map[string]string{"reason": err.Error()})
}

func (c *ResourceController) success(ctx http.Context, data any) http.Response {
	return ctx.Response().Success().Json(response.Success(data, response.Meta{RequestID: resourceRequestID(ctx)}))
}

func (c *ResourceController) error(ctx http.Context, status int, code, message string, details any) http.Response {
	return ctx.Response().Status(status).Json(apierrors.New(code, message, details))
}

func resourceRequestID(ctx http.Context) string {
	requestID := ctx.Request().Header("X-Request-ID")
	if requestID == "" {
		requestID = "resource-request"
	}
	ctx.Response().Header("X-Request-ID", requestID)
	return requestID
}
