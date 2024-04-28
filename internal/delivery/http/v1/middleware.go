package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	id, err := h.parseAuthHeader(c)
	if err != nil {
		newResponse(c, http.StatusUnauthorized, err.Error())
	}

	c.Set(userCtx, id)
}
func (h *Handler) parseAuthHeader(c *gin.Context) (string, error) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return h.tokenManager.Parse(headerParts[1])
}

func getUserId(c *gin.Context) (uint, error) {
	return getIdByContext(c, userCtx)
}

func getIdByContext(c *gin.Context, context string) (uint, error) {
	idFromCtx, ok := c.Get(context)
	if !ok {
		return 0, errors.New("user ID not found in context")
	}

	idStr, ok := idFromCtx.(string)
	if !ok {
		return 0, errors.New("user ID in context is of invalid type")
	}

	idUint, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, errors.New("failed to parse user ID as uint")
	}

	return uint(idUint), nil
}
