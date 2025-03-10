package entity

// Services types in registry
const (
	AuthServiceType         = "auth-service"
	UserManagerServiceType  = "user-manager"
	TicketServiceType       = "ticket-service"
	NotificationServiceType = "notification-service"
)

var ServiceTypes = []string{
	AuthServiceType,
	UserManagerServiceType,
	TicketServiceType,
	NotificationServiceType,
}
