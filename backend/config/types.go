package config

type Config struct {
	Server   ServerConfig   `json:"server" yaml:"server"`
	WebAuthn WebAuthnConfig `json:"webauthn" yaml:"webauthn"`
	Database DatabaseConfig `json:"database" yaml:"database"`
}

type ServerConfig struct {
	Host string `json:"host" yaml:"host"`
	Port int    `json:"port" yaml:"port"`
}

type WebAuthnConfig struct {
	Id          string `json:"id" yaml:"id"`
	DisplayName string `json:"displayName" yaml:"displayName"`
	Origin      string `json:"origin" yaml:"origin"`
}

type DatabaseConfig struct {
	Host         string `json:"host" yaml:"host"`
	Port         int    `json:"port" yaml:"port"`
	Username     string `json:"username" yaml:"username"`
	Password     string `json:"password" yaml:"password"`
	DatabaseName string `json:"dbname" yaml:"dbname"`
}
