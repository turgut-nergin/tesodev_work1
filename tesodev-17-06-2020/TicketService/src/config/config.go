package config

type Config struct {
	Host                 string
	Port                 string
	DBName               string
	TicketCollection     string
	AttachmentCollection string
	AnswerCollection     string
	UserName             string
	Password             string
}

var EnvConfig = map[string]Config{

	"local": {
		Host:                 "localhost",
		Port:                 "27017",
		DBName:               "Ticket",
		TicketCollection:     "Tickets",
		AttachmentCollection: "Attachments",
		AnswerCollection:     "Answers",
		UserName:             "negrin",
		Password:             "Naciye.111",
	},
}
