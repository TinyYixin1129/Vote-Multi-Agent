package multiagentvote

import (
	"time"
)

type Ballot_Req struct {
	Rule      string    `json:"rule"`
	Deadline  time.Time `json:"deadline"`
	Voter_ids []string  `json:"voter-ids"`
	Num_alts  int       `json:"#alts"`
	Tie_break []int     `json:"tie-break,omitempty"`
}

type Ballot_Res struct {
	Ballot_id string `json:"ballot-id"`
}

type Vote_Req struct {
	Agent_id  string `json:"agent-id"`
	Ballot_id string `json:"ballot-id"`
	Prefs     []int  `json:"prefs"`
	Options   []int  `json:"options,omitempty"`
}

type Result_Req struct {
	Ballot_id string `json:"ballot-id"`
}

type Result_Res struct {
	Winner  int   `json:"winner"`
	Ranking []int `json:"ranking"`
}
