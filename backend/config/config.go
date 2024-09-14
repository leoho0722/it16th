package config

func GetServerConfiguration() ServerConfig {
	return cfg.Server
}

func GetWebAuthnConfiguration() WebAuthnConfig {
	return cfg.WebAuthn
}

func GetDatabaseConfiguration() DatabaseConfig {
	return cfg.Database
}
