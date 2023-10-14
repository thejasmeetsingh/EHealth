package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thejasmeetsingh/EHealth/handlers"
	"github.com/thejasmeetsingh/EHealth/internal/database"
	"github.com/thejasmeetsingh/EHealth/utils"
)

type authHandler func(*gin.Context, database.User)

func JWTAuth(apiCfg handlers.ApiCfg, handler authHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		headerAuthToken := ctx.GetHeader("Authorization")

		if headerAuthToken == "" {
			handlers.ErrorResponse(ctx, http.StatusForbidden, "Authentication required")
			ctx.Abort()
			return
		}

		authToken := strings.Split(headerAuthToken, " ")

		if len(authToken) != 2 || authToken[0] != "Bearer" {
			handlers.ErrorResponse(ctx, http.StatusForbidden, "Invalid authentication format")
			ctx.Abort()
			return
		}

		claims, err := utils.VerifyAccessToken(authToken[1])
		if err != nil {
			handlers.ErrorResponse(ctx, http.StatusForbidden, fmt.Sprintf("Error caught while verifying the token: %v", err))
			ctx.Abort()
			return
		}

		dbUser, err := apiCfg.DB.GetUserById(ctx, claims.UserID)
		if err != nil {
			handlers.ErrorResponse(ctx, http.StatusForbidden, "Invalid authentication token")
			ctx.Abort()
			return
		}

		handler(ctx, dbUser)
	}
}
