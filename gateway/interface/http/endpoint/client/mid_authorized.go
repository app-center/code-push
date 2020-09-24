package client

import (
	"github.com/funnyecho/code-push/gateway"
	"github.com/funnyecho/code-push/gateway/interface/http/endpoint"
	res "github.com/funnyecho/code-push/pkg/ginkit/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func MidAuthorized(c *gin.Context) {
	envId, authErr := authorized(c)
	if authErr != nil {
		res.ErrorWithStatusCode(c, http.StatusUnauthorized, authErr)
		return
	}

	WithEnvId(string(envId), c)
	c.Next()
}

func authorized(c *gin.Context) ([]byte, error) {
	var accessToken string

	accessTokenFromCookies, cookiesErr := c.Cookie("access-token")

	if cookiesErr == nil {
		accessToken = accessTokenFromCookies
	} else {
		accessToken = c.Query("access-token")
		if len(accessToken) == 0 {
			accessToken = c.GetHeader("X-Authentication")
		}
	}

	if len(accessToken) == 0 {
		return nil, gateway.ErrUnauthorized
	}

	envId, verifyErr := endpoint.UseUC(c).VerifyTokenForEnv(c.Request.Context(), []byte(accessToken))
	if verifyErr != nil {
		return nil, gateway.ErrInvalidToken
	}

	return envId, nil
}

