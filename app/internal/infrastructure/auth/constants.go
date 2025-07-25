package auth

type AuthMode string

const (
	AuthModeUser    AuthMode = "stackspot_user"
	AuthModeService AuthMode = "stackspot_service"
)