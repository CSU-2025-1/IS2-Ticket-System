package proxy

import "gateway/internal/entity"

var RoutePolicy = map[string]string{
	"/api/ticket": entity.PublicTicketServiceType,
	"/api/user":   entity.PublicUserManagerServiceType,
	"/api/auth":   entity.PublicAuthServiceType,
}

var AuthPolicy = map[string]bool{
	"/api/ticket": true,
	"/api/user":   true,
}
