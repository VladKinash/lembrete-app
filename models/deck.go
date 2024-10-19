package models

type Deck struct {
	MaxNewCards     int32
	MaxReviewsDaily int32
	Name            string
	Id              int32
}

func newDeck(MaxNewCards int32, MaxReviewsDaily int32, Name string, Id int32) Deck {
	return Deck{

		MaxNewCards:     MaxNewCards,
		MaxReviewsDaily: MaxReviewsDaily,
		Name:            Name,
		Id:              Id}
}
