package models

type UpSertResult struct {
	ModifiedCount int64
	ID            string
	Err           error
	ErrCode       int
}
