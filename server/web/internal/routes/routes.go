package routes

import (
	"github.com/jhamill34/prophet-security-takehome/server/api/pkg/api"
	"github.com/jhamill34/prophet-security-takehome/server/database/pkg/database"
)

type ServerRoutes struct {
	queries *database.Queries
}

func NewServerRoutes(queries *database.Queries) *ServerRoutes {
	return &ServerRoutes{
		queries,
	}
}

var _ api.StrictServerInterface = (*ServerRoutes)(nil)
