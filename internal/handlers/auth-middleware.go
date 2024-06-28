package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	incorrectData    = "Некорректные данные"
	userUnauthorized = "Пользователь не авторизован"
	errorAuthorize   = "ошибка авторизации"
	errorParseToken  = "ошибка парсинга токена"
)

func (h *Handlers) authMiddleware(c *gin.Context) {
	const op = "handlers.authMiddleware"

	header := c.GetHeader("Authorization")
	if header == "" {
		h.logger.Error("Error get header (Authorization): ", op)
		newErrorResponse(c, http.StatusUnauthorized, userUnauthorized)
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		h.logger.Error("Error format header (Authorization): ", op)
		newErrorResponse(c, http.StatusUnauthorized, errorAuthorize)
		return
	}

	email, err := h.service.Auth.ParseToken(headerParts[1])
	if err != nil {
		h.logger.Error("Error parse token: ", op)
		newErrorResponse(c, http.StatusUnauthorized, errorParseToken)
		return
	}

	c.Set("email", email)
	c.Next()
}
