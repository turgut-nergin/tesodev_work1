package config

type Config struct {
	Host             string
	Port             string
	DBName           string
	TicketCollection string
	AnswerCollection string
	UserName         string
	Password         string
}

var EnvConfig = map[string]Config{

	"local": {
		Host:             "localhost",
		Port:             "27017",
		DBName:           "Ticket",
		TicketCollection: "Tickets",
		AnswerCollection: "Answers",
		UserName:         "negrin",
		Password:         "Naciye.111",
	},
}
