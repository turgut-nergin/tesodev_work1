package config

type Config struct {
	Host           string
	Port           string
	DBName         string
	CollectionName string
	UserName       string
	Password       string
	MaxPageLimit   int64
}

var EnvConfig = map[string]Config{

	"local": {
		Host:           "mongo-db",
		Port:           "27017",
		DBName:         "User",
		CollectionName: "users",
		UserName:       "negrin",
		Password:       "Naciye.111",
		MaxPageLimit:   100,
	},
}
