package entity

// Services types in registry
const (
	PublicAuthServiceType          = "auth-service-test"
	PublicUserManagerServiceType   = "user-manager-rest"
	PublicTicketServiceType        = "ticket-service-rest"
	PublicNotificationServiceType  = "notification-service-rest"
	PrivateAuthServiceType         = "auth-service-grpc"
	PrivateUserManagerServiceType  = "user-manager-grpc"
	PrivateTicketServiceType       = "ticket-service-grpc"
	PrivateNotificationServiceType = "notification-service-grpc"
)

var ServiceTypes = []string{
	PublicAuthServiceType,
	PublicUserManagerServiceType,
	PublicTicketServiceType,
	PublicNotificationServiceType,
	PrivateAuthServiceType,
	PrivateUserManagerServiceType,
	PrivateTicketServiceType,
	PrivateNotificationServiceType,
}
