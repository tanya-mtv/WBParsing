package config

type ConfigMSSQL struct {
	Driver   string `json:"driver"`
	Server   string `json:"server"`
	User     string `json:"user"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	DSN      string `json:"dsn"`
}
