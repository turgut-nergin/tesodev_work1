package config

type Config struct {
	Host           string
	Port           string
	DBName         string
	CollectionName string
	UserName       string
	Password       string
}

var EnvConfig = map[string]Config{

	"local": {
		Host:           "localhost",
		Port:           "27017",
		DBName:         "User",
		CollectionName: "users",
		UserName:       "negrin",
		Password:       "Naciye.111",
	},
}
