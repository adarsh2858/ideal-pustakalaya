package config

type Config struct {
	ProjectID  string `env:"PROJECT_ID" envDefault:"local"`
	ServerName string `env:"SERVER_NAME" envDefault:"{{.repoName}}"`
	LogLevel   string `env:"LOG_LEVEL" envDefault:"debug"`
	SentryDSN  string `env:"SENTRY_DSN" envDefault:"" yaml:"sentry_dsn" exportENV:"SENTRY_DSN"`

	GRPC   GRPCConfig
	PubSub PubSubConfig
}

type GRPCConfig struct {
	Host string `env:"GRPC_HOST"`
	Port int    `env:"GRPC_PORT"`
}

type PubSubConfig struct {
	MaxMessages int    `env:"PUBSUB_MAX_MESSAGES" envDefault:"3"`
	MaxRetries  int    `env:"PUBSUB_MAX_RETRIES" envDefault:"0"`
	AckTimeout  int    `env:"PUBSUB_ACK_TIMEOUT" envDefault:"600"`
	HelloSub    string `env:"HELLO_SUBSCRIPTION"`
	HelloTopic  string `env:"HELLO_TOPIC"`
}
