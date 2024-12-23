package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/google/uuid"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

func ExtractUserIDAndEmail(c *gin.Context) (uuid.UUID, string, error) {
	session := sessions.Default(c)
	useridStr, ok := session.Get("userid").(string)
	if !ok {
		return uuid.Nil, "", internalErrors.NewAPIError(http.StatusUnauthorized, "failed to get userid from session")
	}

	userid, err := uuid.Parse(useridStr)
	if err != nil {
		return uuid.Nil, "", internalErrors.NewAPIError(http.StatusBadRequest, "invalid userid format")
	}

	email, ok := c.Get("email")
	if !ok {
		return uuid.Nil, "", internalErrors.NewAPIError(http.StatusUnauthorized, "failed to get email from context")
	}

	emailStr, ok := email.(string)
	if !ok {
		return uuid.Nil, "", internalErrors.NewAPIError(http.StatusBadRequest, "invalid email format")
	}

	return userid, emailStr, nil
}