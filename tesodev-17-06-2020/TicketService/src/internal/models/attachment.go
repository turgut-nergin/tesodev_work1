package models

type Attachment struct {
	TicketId string `bson:"ticketId" json:"ticketId,omitempty"`
	FileName string `bson:"fileName" json:"fileName"`
	FilePath string `bson:"filePath" json:"filePath"`
}
