package redis

import "time"

type GenerateImagePayload struct {
	MatchDate time.Time `json:"matchDate"`
}

type SendToDiscordPayload struct {
	JournalPath string `json:"journalPath"`
}
