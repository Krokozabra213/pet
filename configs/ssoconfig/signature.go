package ssoconfig

var (
	DSN                 = "SSO_DSN"
	HOST                = "SSO_AUTH_HOST"
	PORT                = "SSO_AUTH_PORT"
	PrivateKey          = "SSO_RSA_PRIVATE_KEY_PATH"
	Secret              = "SSO_APP_SECRET"
	DevAccessTokenTTL   = "SSO_DEV_ACCESS_TOKEN_TTL"
	ProdAccessTokenTTL  = "SSO_PROD_ACCESS_TOKEN_TTL"
	DevRefreshTokenTTL  = "SSO_DEV_REFRESH_TOKEN_TTL"
	ProdRefreshTokenTTL = "SSO_PROD_REFRESH_TOKEN_TTL"
	ContextTimeout      = "SSO_CTX_TIMEOUT"
	RedisAddr           = "SSO_REDIS_ADDR"
	RedisPass           = "SSO_REDIS_PASS"
	RedisCache          = "SSO_REDIS_CACHE"
)
