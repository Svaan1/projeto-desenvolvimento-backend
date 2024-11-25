package spotify

type ErrRecommendationsEmpty struct {
	Message string
}

func (e *ErrRecommendationsEmpty) Error() string {
	return e.Message
}
