package config

type RabbitMQConfig struct {
	Host     string
	Port     string
	UserName string
	Password string
}

var RabbitMQEnvConfig = map[string]RabbitMQConfig{

	"local": {
		Host:     "rabbitmq",
		Port:     "5672",
		UserName: "user",
		Password: "password",
	},
}
