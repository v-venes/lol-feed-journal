package imagegenerator

import "time"

type GenerateImagePayload struct {
	MatchDate time.Time `json:"matchDate"`
}
