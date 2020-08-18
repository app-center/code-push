package middleware

import (
	"github.com/funnyecho/code-push/gateway/client"
	"github.com/funnyecho/code-push/gateway/client/interface/http/constants"
	"github.com/funnyecho/code-push/gateway/client/usecase"
	res "github.com/funnyecho/code-push/pkg/gin_res"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (m *Middleware) Authorized(c *gin.Context) {
	branchId, authErr := AuthorizedWithReturns(m.uc, c)
	if authErr != nil {
		res.ErrorWithStatusCode(c, http.StatusUnauthorized, authErr)
		return
	}

	c.Set(constants.CtxEnvId, branchId)
	c.Next()
}

func AuthorizedWithReturns(uc usecase.UseCase, c *gin.Context) ([]byte, error) {
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
		return nil, client.ErrUnauthorized
	}

	envId, verifyErr := uc.VerifyToken([]byte(accessToken))
	if verifyErr != nil {
		return nil, client.ErrInvalidToken
	}

	return envId, nil
}
