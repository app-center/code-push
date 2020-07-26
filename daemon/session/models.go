package session

type AccessTokenIssuer int32

const (
	AccessTokenIssuerSYS    AccessTokenIssuer = 0
	AccessTokenIssuerPORTAL AccessTokenIssuer = 1
	AccessTokenIssuerCLIENT AccessTokenIssuer = 2
)

type AccessTokenClaims struct {
	Issuer   AccessTokenIssuer
	Subject  string
	Audience []byte
}
