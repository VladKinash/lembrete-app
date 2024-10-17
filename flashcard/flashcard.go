package flashcard

type Flashcard struct {
	Front       string
	Back        string
	EaseFactor  float32
	Repetitions int 
	Interval    float32
	NextReview  string
}

func NewFlashcard(front string, back string, easeFactor float32, repetitions int, interval float32, nextReview string) Flashcard {
	return Flashcard{
		Front:       front,
		Back:        back,
		EaseFactor:  easeFactor,
		Repetitions: repetitions,
		Interval:    interval,
		NextReview:  nextReview,
	}
}
