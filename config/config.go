package config

type Config struct {
	HttpPort string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string
	SSLMode          string
	JWTsecretkey     []byte
}

func Cfg() Config {
	cfg := Config{}

	cfg.HttpPort = ":8080"

	cfg.PostgresHost = "localhost"
	cfg.PostgresUser = "postgres"
	cfg.PostgresDatabase = "postgres"
	cfg.PostgresPassword = "my_password_p"
	cfg.PostgresPort = "5434"
	cfg.SSLMode = "disable"
	cfg.JWTsecretkey = []byte("your_secret_key")

	return cfg
}
