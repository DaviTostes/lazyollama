package model

type Embedding struct {
	Id     int
	Text   string
	Emb    []float64
	ChatId int64
}
