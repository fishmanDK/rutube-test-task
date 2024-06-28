package handlers

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func (h *Handlers) subscribe(c *gin.Context) {
	const op = "handlers.subscribe"

	rootEmail := c.Query("root-email")
	subscriberEmail := c.Query("subs-email")

	if rootEmail == "" || subscriberEmail == "" {
		h.logger.Error(op, slog.String("err", incorrectData))
		newErrorResponse(c, http.StatusBadRequest, incorrectData)
		return
	}

	if err := h.service.Api.Subscribe(rootEmail, subscriberEmail); err != nil {
		h.logger.Error(op, slog.String("err", err.Error()))
		newErrorResponse(c, http.StatusBadRequest, incorrectData)
	}

	c.JSON(http.StatusOK, nil)
}

func (h *Handlers) unsubscribe(c *gin.Context) {
	const op = "handlers.unsubscribe"

	rootEmail := c.Query("root-email")
	subscriberEmail := c.Query("subs-email")

	if rootEmail == "" || subscriberEmail == "" {
		h.logger.Error(op, slog.String("err", incorrectData))
		newErrorResponse(c, http.StatusBadRequest, incorrectData)
		return
	}

	if err := h.service.Api.Unsubscribe(rootEmail, subscriberEmail); err != nil {
		h.logger.Error(op, slog.String("err", err.Error()))
		newErrorResponse(c, http.StatusBadRequest, incorrectData)
	}

	c.JSON(http.StatusOK, nil)
}

func (h *Handlers) allSubs(c *gin.Context) {
	userEmail, _ := c.Get("email")

	subs := h.service.Api.GetAllSubs(userEmail.(string))

	c.JSON(http.StatusOK, subs)
}
