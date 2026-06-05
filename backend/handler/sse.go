package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func (h *Handler) SSE(c echo.Context) error {
	binID := c.Param("id")

	if _, err := h.store.GetBin(c.Request().Context(), binID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "bin not found")
		}
		return err
	}

	w := c.Response().Writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "streaming not supported")
	}

	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("X-Accel-Buffering", "no")
	c.Response().WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "event: connected\ndata: {}\n\n")
	flusher.Flush()

	ch := h.hub.Subscribe(binID)
	defer h.hub.Unsubscribe(binID, ch)

	ctx := c.Request().Context()
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-ch:
			if !ok {
				return nil
			}
			fmt.Fprintf(w, "event: request\ndata: %s\n\n", data)
			flusher.Flush()
		}
	}
}
