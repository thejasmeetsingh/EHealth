// A common struct for all the handlers. So that handler can query the DB without any external module import

package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/thejasmeetsingh/EHealth/internal/database"
)

type ApiCfg struct {
	DB *database.Queries
}

// Common function for getting user object from the context, For all handlers
func getDBUser(c *gin.Context) (database.User, error) {
	user, exists := c.Get("user")

	if !exists {
		return database.User{}, fmt.Errorf("authentication required")
	}

	dbUser, ok := user.(database.User)

	if !ok {
		return database.User{}, fmt.Errorf("invalid user provied")
	}

	return dbUser, nil
}
