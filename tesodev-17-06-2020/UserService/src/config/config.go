package config

type Config struct {
	Host           string
	Port           string
	DBName         string
	CollectionName string
}

var EnvConfig = map[string]Config{
	"local": {
		Host:           "mongodb://localhost",
		Port:           "27017",
		DBName:         "Tickets",
		CollectionName: "users",
	},
}
