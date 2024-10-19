package models

type Flashcard struct {
	Front       string
	Back        string
	EaseFactor  float32
	Repetitions int
	Interval    float32
	NextReview  string
	DeckID      string
	Id          int32
}

func NewFlashcard(
	front string,
	back string,
	easeFactor float32,
	repetitions int,
	interval float32,
	nextReview string,
	DeckID string,
	Id int32,
) Flashcard {
	return Flashcard{
		Front:       front,
		Back:        back,
		EaseFactor:  easeFactor,
		Repetitions: repetitions,
		Interval:    interval,
		NextReview:  nextReview,
		DeckID:      DeckID,
		Id:          Id,
	}
}
