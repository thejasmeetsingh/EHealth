// A common struct for all the handlers. So that handler can query the DB without any external module import

package handlers

import "github.com/thejasmeetsingh/EHealth/internal/database"

type ApiCfg struct {
	DB *database.Queries
}
