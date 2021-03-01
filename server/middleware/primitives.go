package middleware

//RoleHolder user interface with role
type RoleHolder interface {
	IsInRole(role string) bool
}
type FromJWT func(jwt string) (uinfo RoleHolder, err error)
