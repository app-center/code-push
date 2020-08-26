package middleware

import (
	"github.com/funnyecho/code-push/gateway/portal"
	"github.com/funnyecho/code-push/gateway/portal/interface/http/constants"
	"github.com/funnyecho/code-push/gateway/portal/usecase"
	res "github.com/funnyecho/code-push/pkg/gin-response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (m *Middleware) Authorized(c *gin.Context) {
	branchId, authErr := AuthorizedWithReturns(m.uc, c)
	if authErr != nil {
		res.ErrorWithStatusCode(c, http.StatusUnauthorized, authErr)
		return
	}

	c.Set(constants.CtxBranchId, branchId)
	c.Next()
}

func AuthorizedWithReturns(uc usecase.UseCase, c *gin.Context) ([]byte, error) {
	var accessToken string

	accessTokenFromCookies, cookiesErr := c.Cookie("access-token")

	if cookiesErr == nil {
		accessToken = accessTokenFromCookies
	} else {
		accessToken = c.Query("access-token")
		if accessToken == "" {
			accessToken = c.GetHeader("Portal-Access-Token")
		}
	}

	if len(accessToken) == 0 {
		return nil, portal.ErrUnauthorized
	}

	branchId, verifyErr := uc.VerifyToken(c.Request.Context(), []byte(accessToken))
	if verifyErr != nil {
		return nil, portal.ErrInvalidToken
	}

	return branchId, nil
}
