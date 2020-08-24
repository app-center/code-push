package middleware

import (
	"github.com/funnyecho/code-push/gateway/sys"
	res "github.com/funnyecho/code-push/pkg/ginResponse"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (m *Middleware) Authorized(c *gin.Context) {
	var accessToken string

	accessTokenFromCookies, cookiesErr := c.Cookie("access-token")

	if cookiesErr == nil {
		accessToken = accessTokenFromCookies
	} else {
		accessToken = c.Query("access-token")
		if len(accessToken) == 0 {
			accessToken = c.GetHeader("Sys-Access-Token")
		}
	}

	if len(accessToken) == 0 {
		res.ErrorWithStatusCode(c, http.StatusUnauthorized, sys.ErrUnauthorized)
		return
	}

	verifyErr := m.uc.VerifyToken(c.Request.Context(), []byte(accessToken))
	if verifyErr != nil {
		res.ErrorWithStatusCode(c, http.StatusUnauthorized, sys.ErrInvalidToken)
		return
	}

	c.Next()
}
