package models
	
type ReviewQueue struct {
	newCards []*Flashcard
	dueCards []*Flashcard
	current int
}

func NewReviewQueue(newCards, dueCards []*Flashcard) *ReviewQueue {
	return &ReviewQueue{
		newCards: newCards,
		dueCards: dueCards,
		current: 0,
	}
}


func (rq *ReviewQueue) Next() *Flashcard{
	if rq.current < len(rq.newCards){
		card := rq.newCards[rq.current]
		rq.current++
		return card
	} else if rq.current-len(rq.newCards) < len(rq.newCards){
		card := rq.dueCards[rq.current-len(rq.newCards)]
		rq.current++
		return card
	}
	return nil
}