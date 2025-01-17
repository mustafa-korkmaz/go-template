package transport

import (
	"net/http"

	"github.com/mustafa-korkmaz/goapitemplate/pkg/enum"
	"github.com/mustafa-korkmaz/goapitemplate/pkg/viewmodel/request"
	"github.com/mustafa-korkmaz/goapitemplate/pkg/viewmodel/response"

	"github.com/labstack/echo"
	"github.com/mustafa-korkmaz/goapitemplate/pkg/api/healthcheck"
)

// HTTP represents healthcheck http transport service
type HTTP struct {
	svc healthcheck.Service
}

// New creates new healthcheck http service with valid api versions
func New(svc healthcheck.Service, groups ...*echo.Group) {
	h := HTTP{svc}
	v1 := groups[0].Group("/healthcheck")
	v2 := groups[1].Group("/healthcheck")

	//define /V1/healtcheck methods
	v1.GET("/:value", h.get)
	v1.POST("/paginationtest", h.getPagedList)
	v1.POST("", h.post)

	//define /V2/healtcheck methods
	v2.GET("/:value", h.getV2)
}

func (h *HTTP) get(c echo.Context) error {

	val := c.Param("value")

	if err := h.svc.Get(val); err != nil {
		return err
	}

	var resp = response.APIResponse{}
	resp.Result = enum.ResponseResult.Success

	resp.Data = struct {
		AwesomeCars []string `json:"awesome_cars"`
		Value       string   `json:"value"`
	}{
		[]string{"Wv", "Jaguar", "Tesla"},
		val,
	}

	return c.JSON(http.StatusOK, &resp)
}

func (h *HTTP) getPagedList(c echo.Context) error {

	req := request.PagedList{}
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	var pagedListResp, err = h.svc.GetPagedList(req)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &pagedListResp)
}

func (h *HTTP) post(c echo.Context) error {

	req := HealthCheckReq{}
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	if err := h.svc.Post(req.Value); err != nil {
		return err
	}

	var resp = response.APIResponse{}

	resp.Result = enum.ResponseResult.Success

	resp.Data = struct {
		AwesomePhones []string `json:"awesome_phones"`
		Value         string   `json:"value"`
	}{
		[]string{"iPhone", "Samsung", "Huawei"},
		req.Value,
	}

	return c.JSON(http.StatusOK, &resp)
}

func (h *HTTP) getV2(c echo.Context) error {

	val := c.Param("value")
	val = "Your V2 value is: " + val

	if err := h.svc.Get(val); err != nil {
		return err
	}

	var resp = response.APIResponse{}
	resp.Result = enum.ResponseResult.Success

	resp.Data = struct {
		AwesomeCars []string `json:"awesome_cars"`
		Value       string   `json:"value"`
	}{
		[]string{"VW", "Jaguar", "Tesla"},
		val,
	}

	return c.JSON(http.StatusOK, &resp)
}

// HealthCheckReq represents body of HealthCheck request.
type HealthCheckReq struct {
	Value string `json:"value" validate:"required,min=8"`
}
