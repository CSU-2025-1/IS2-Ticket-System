package entity

// Services types in registry
const (
	PublicAuthServiceType         = "public-auth-service"
	PublicUserManagerServiceType  = "public-user-manager"
	PublicTicketServiceType       = "public-ticket-service"
	PublicNotificationServiceType = "public-notification-service"
	PrivateAuthServiceType        = "private-auth-service"
	PrivateUserManagerServiceType = "private-user-manager"
	PrivateTicketServiceType      = "private-ticket-service"
)

var ServiceTypes = []string{
	PublicAuthServiceType,
	PublicUserManagerServiceType,
	PublicTicketServiceType,
	PublicNotificationServiceType,
	PrivateAuthServiceType,
	PrivateUserManagerServiceType,
	PrivateTicketServiceType,
}
