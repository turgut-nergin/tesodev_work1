package user

type User struct {
	ID        string `bson:",omitempty"`
	UserName  string `bson:"userName"`
	Password  string `bson:"password"`
	Email     string `bson:"email"`
	Type      string `bson:"type"`
	CreatedAt int64  `bson:",omitempty"`
	UpdatedAt int64  `bson:",omitempty"`
}
