package apiserver

// Config is a interface for starting server
type Config struct {
	BindAddr    string `toml:"bind_addr"`
	LogLevel    string `toml:"log_level"`
	DatabaseURL string `toml:"database_url"`
	JwtKey      []byte `toml:"jwt_key"`
}

// NewConfig is a function for creating a new Config interface for starting server
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
		JwtKey:   []byte("secret"),
	}
}
