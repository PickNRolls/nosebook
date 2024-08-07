package dto

type QueryResultDTO struct {
	Err            error      `json:"error"`
	RemainingCount int        `json:"remainingCount"`
	Data           []*PostDTO `json:"data"`
	Next           string     `json:"next"`
}
