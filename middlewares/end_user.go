package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thejasmeetsingh/EHealth/handlers"
	"github.com/thejasmeetsingh/EHealth/internal/database"
)

// Middleware for checking if is it an end user or not
func EndUser(ctx *gin.Context) {
	user, exists := ctx.Get("user")

	if !exists {
		handlers.ErrorResponse(ctx, http.StatusForbidden, "Access Restricted")
		ctx.Abort()
		return
	}

	dbUser, ok := user.(database.User)

	if !ok {
		handlers.ErrorResponse(ctx, http.StatusInternalServerError, "Something went wrong")
		ctx.Abort()
		return
	}

	if !dbUser.IsEndUser {
		handlers.ErrorResponse(ctx, http.StatusForbidden, "You cannot access this resource")
		ctx.Abort()
		return
	}

	ctx.Next()
}
