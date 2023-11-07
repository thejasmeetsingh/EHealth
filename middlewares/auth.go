package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thejasmeetsingh/EHealth/handlers"
	"github.com/thejasmeetsingh/EHealth/utils"
)

// Validate the request by checking wheather or not they have the valid JWT access token or not
//
// Token format: Bearer <TOKEN>
func JWTAuth(apiCfg handlers.ApiCfg) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		headerAuthToken := ctx.GetHeader("Authorization")

		if headerAuthToken == "" {
			handlers.ErrorResponse(ctx, http.StatusForbidden, "Authentication required")
			ctx.Abort()
			return
		}

		// Split the token string
		authToken := strings.Split(headerAuthToken, " ")

		// Validate the token string
		if len(authToken) != 2 || authToken[0] != "Bearer" {
			handlers.ErrorResponse(ctx, http.StatusForbidden, "Invalid authentication format")
			ctx.Abort()
			return
		}

		// Verify the token and get the encoded payload which is the userID string
		claims, err := utils.VerifyToken(authToken[1])
		if err != nil {
			handlers.ErrorResponse(ctx, http.StatusForbidden, fmt.Sprintf("Error caught while verifying the token: %v", err))
			ctx.Abort()
			return
		}

		// Check the validity of the token
		if !time.Unix(claims.ExpiresAt.Unix(), 0).After(time.Now()) {
			handlers.ErrorResponse(ctx, http.StatusForbidden, "Invalid authentication token")
			ctx.Abort()
			return
		}

		// Conver the userID string to UUID
		userID, err := uuid.Parse(claims.Data)
		if err != nil {
			handlers.ErrorResponse(ctx, http.StatusForbidden, "Invalid authentication token")
			ctx.Abort()
			return
		}

		// Fetch the user by the ID
		dbUser, err := apiCfg.DB.GetUserById(ctx, userID)
		if err != nil {
			handlers.ErrorResponse(ctx, http.StatusForbidden, "Invalid authentication token")
			ctx.Abort()
			return
		}

		ctx.Set("user", dbUser)

		// Further call the given handler and send the user instance as well
		ctx.Next()
	}
}
