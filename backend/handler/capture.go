package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/ready4god2513/requestbin/model"
)

func (h *Handler) Capture(c echo.Context) error {
	binID := c.Param("id")

	if _, err := h.store.GetBin(c.Request().Context(), binID); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "bin not found")
	}

	req := c.Request()

	body, err := io.ReadAll(io.LimitReader(req.Body, 1<<20)) // 1 MB limit
	if err != nil {
		return err
	}

	headers := make(map[string]string, len(req.Header))
	for k, vs := range req.Header {
		headers[k] = strings.Join(vs, ", ")
	}

	queryParams := make(map[string]string, len(req.URL.Query()))
	for k, vs := range req.URL.Query() {
		queryParams[k] = strings.Join(vs, ", ")
	}

	subPath := c.Param("*")
	if subPath == "" {
		subPath = "/"
	} else {
		subPath = "/" + subPath
	}
	if req.URL.RawQuery != "" {
		subPath += "?" + req.URL.RawQuery
	}

	contentLength := req.ContentLength
	if contentLength < 0 {
		contentLength = int64(len(body))
	}

	r := &model.Request{
		BinID:         binID,
		Method:        req.Method,
		Path:          subPath,
		Headers:       headers,
		QueryParams:   queryParams,
		Body:          string(body),
		RemoteAddr:    req.RemoteAddr,
		ContentType:   req.Header.Get("Content-Type"),
		ContentLength: contentLength,
	}

	saved, err := h.store.CreateRequest(c.Request().Context(), r)
	if err != nil {
		return err
	}

	if data, err := json.Marshal(saved); err == nil {
		h.hub.Broadcast(binID, data)
	}

	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}
