package algorithm

import (
	//"ftm"
	f "Lembrete/models"
	"errors"
	"math"
)

func algorithm(card f.Flashcard, quality float32) (f.Flashcard, error) {

	if quality < 0 || quality > 5 {
		return card, errors.New("quality cannot be larger than 5 and smaller than 0")
	}

	if quality >= 3 {
		if card.Repetitions == 0 {
			card.Interval = 1
		} else if card.Repetitions == 1 {
			card.Interval = 6
		} else if card.Repetitions > 1 {
			card.Interval = roundToTwoDecimals((card.Interval * card.EaseFactor))
			card.Repetitions++
			updatedCard := setEaseFactor(card, quality)
			return updatedCard, nil
		}
	}

	if quality < 3 {
		card.Repetitions = 0
		card.Interval = 1
		if card.EaseFactor < 1.3 {
			card.EaseFactor = 1.3
		}
	}

	return card, nil
}

func roundToTwoDecimals(num float32) float32 {
	return float32(math.Round(float64(num)*100) / 100)
}

func setEaseFactor(card f.Flashcard, quality float32) f.Flashcard {
	card.EaseFactor = card.EaseFactor + (0.1 - (5-quality)*(0.08+(5-quality)*0.02))
	return card
}
