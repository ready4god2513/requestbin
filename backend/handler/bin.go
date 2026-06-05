package handler

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateBin(c echo.Context) error {
	bin, err := h.store.CreateBin(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, bin)
}

func (h *Handler) GetBin(c echo.Context) error {
	bin, err := h.store.GetBin(c.Request().Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "bin not found")
		}
		return err
	}
	return c.JSON(http.StatusOK, bin)
}

func (h *Handler) DeleteBin(c echo.Context) error {
	if err := h.store.DeleteBin(c.Request().Context(), c.Param("id")); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) ListRequests(c echo.Context) error {
	requests, err := h.store.ListRequests(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, requests)
}

func (h *Handler) ClearRequests(c echo.Context) error {
	if err := h.store.ClearRequests(c.Request().Context(), c.Param("id")); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) DeleteRequest(c echo.Context) error {
	if err := h.store.DeleteRequest(c.Request().Context(), c.Param("reqID")); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
