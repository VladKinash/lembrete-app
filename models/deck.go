package models

type Deck struct {
	ID              int32
	Name            string
	MaxNewCards     int32
	MaxReviewsDaily int32
}

func NewDeck(MaxNewCards int32, MaxReviewsDaily int32, Name string, ID int32) Deck {
	return Deck{

		MaxNewCards:     MaxNewCards,
		MaxReviewsDaily: MaxReviewsDaily,
		Name:            Name,
		ID:              ID}
}
