package models

type Deck struct {
	MaxNewCards     int32
	MaxReviewsDaily int32
	Name            string
}

func newDeck(MaxNewCards int32, MaxReviewsDaily int32, Name string) Deck {
	return Deck{

		MaxNewCards:     MaxNewCards,
		MaxReviewsDaily: MaxReviewsDaily,
		Name:            Name}
}
