package config

type MongoConfig struct {
	Host             string
	Port             string
	DBName           string
	TicketCollection string
	AnswerCollection string
	UserName         string
	Password         string
	MaxPageLimit     int64
}

var MongoEnvConfig = map[string]MongoConfig{

	"local": {
		Host:             "mongo-db",
		Port:             "27017",
		DBName:           "Ticket",
		TicketCollection: "Tickets",
		AnswerCollection: "Answers",
		UserName:         "negrin",
		Password:         "Naciye.111",
		MaxPageLimit:     100,
	},
}
