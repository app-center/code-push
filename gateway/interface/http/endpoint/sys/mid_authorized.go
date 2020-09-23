package sys

import (
	"github.com/funnyecho/code-push/gateway"
	"github.com/funnyecho/code-push/gateway/interface/http/endpoint"
	res "github.com/funnyecho/code-push/pkg/ginkit/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func MidAuthorized(c *gin.Context) {
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
		res.ErrorWithStatusCode(c, http.StatusUnauthorized, gateway.ErrUnauthorized)
		return
	}

	verifyErr := endpoint.UseUC(c).VerifyTokenForRootUser(c.Request.Context(), []byte(accessToken))
	if verifyErr != nil {
		res.ErrorWithStatusCode(c, http.StatusUnauthorized, gateway.ErrInvalidToken)
		return
	}

	c.Next()
}
