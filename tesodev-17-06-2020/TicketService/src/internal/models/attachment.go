package models

type Attachment struct {
	Id       string `bson:"_id,omitempty" json:"_id,omitempty"`
	TicketId string `bson:"ticketId" json:"ticketId,omitempty"`
	FileName string `bson:"fileName" json:"fileName"`
	FilePath string `bson:"filePath" json:"filePath"`
}
