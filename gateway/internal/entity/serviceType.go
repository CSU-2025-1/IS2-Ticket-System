package entity

// Services types in registry
const (
	PublicAuthServiceType         = "public-auth-service"
	PublicUserManagerServiceType  = "public-user-manager"
	PublicTicketServiceType       = "public-ticket-service"
	PublicNotificationServiceType = "public-notification-service"
)

var ServiceTypes = []string{
	PublicAuthServiceType,
	PublicUserManagerServiceType,
	PublicTicketServiceType,
	PublicNotificationServiceType,
}
