package proxy

import "gateway/internal/entity"

/*
	Policies for Proxy - for requests with URL prefixes like this

	It can be added in config, but there is not enough time for it
*/

// RoutePolicy is a list of rules for proxying request to each service
var RoutePolicy = map[string]string{
	"/api/ticket":       entity.PublicTicketServiceType,
	"/api/user":         entity.PublicUserManagerServiceType,
	"/api/auth":         entity.PublicAuthServiceType,
	"/api/notification": entity.PublicNotificationServiceType,
}

// AuthPolicy is a policy for auth at routes
var AuthPolicy = map[string]bool{
	"/api/ticket": true,
	"/api/user":   true,
}
