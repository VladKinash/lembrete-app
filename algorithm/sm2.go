package algorithm

import (
	f "Lembrete/models"
	"errors"
	"math"
)

func SM2Algorithm(card *f.Flashcard, quality float32) (*f.Flashcard, error) {
	if quality < 0 || quality > 5 {
		return nil, errors.New("quality must be between 0 and 5") 
	}

	if quality >= 3 {
		switch card.Repetitions {
		case 0:
			card.Repetitions = 1
			card.Interval = 1
		case 1:
			card.Repetitions = 2
			card.Interval = 6
		default:
			card.Repetitions++
			card.Interval = roundToTwoDecimals(card.Interval * card.EaseFactor)
		}
	} else {
		card.Repetitions = 0
		card.Interval = 1
	}

	card.EaseFactor = card.EaseFactor + (0.1 - (5-quality)*(0.08+(5-quality)*0.02))
	if card.EaseFactor < 1.3 {
		card.EaseFactor = 1.3
	}

	return card, nil // Return the pointer
}

func roundToTwoDecimals(num float32) float32 {
	return float32(math.Round(float64(num)*100) / 100)
}
