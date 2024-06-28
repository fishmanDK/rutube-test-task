package handlers

import (
	"github.com/fishmanDK/rutube-test-task/models"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func (h *Handlers) signIn(c *gin.Context) {
	const op = "handlers.signIn"

	var input models.User
	if err := c.BindJSON(&input); err != nil {
		h.logger.Error("Failed to bind JSON: ", err)
		newErrorResponse(c, http.StatusBadRequest, incorrectData)
		return
	}

	tokens, err := h.service.Auth.Authentication(input)
	if err != nil {
		h.logger.Error(op, slog.String("err", err.Error()))
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"email": input.Email, "tokens": tokens})
}

func (h *Handlers) signUp(c *gin.Context) {
	const op = "handlers.signUp"

	var input models.NewUser
	if err := c.BindJSON(&input); err != nil {
		h.logger.Error(op, slog.String("err", err.Error()))
		newErrorResponse(c, http.StatusBadRequest, incorrectData)
		return
	}

	err := h.service.Auth.CreateUser(input)
	if err != nil {
		h.logger.Error(op, slog.String("err", err.Error()))
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func (h *Handlers) apdateAccessToken(c *gin.Context) {
	const op = "handlers.apdateAccessToken"

	var input models.UpdateRefreshTokenRequest
	if err := c.BindJSON(&input); err != nil {
		h.logger.Error(op, slog.String("err", err.Error()))
		newErrorResponse(c, http.StatusBadRequest, incorrectData)
		return
	}

	tokens, err := h.service.Auth.ApdateAccessToken(input)
	if err != nil {
		h.logger.Error(op, slog.String("err", err.Error()))
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//fmt.Println(tokens)

	c.JSON(http.StatusCreated, tokens)
}
