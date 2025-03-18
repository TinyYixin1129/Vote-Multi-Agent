package main

import (
	"fmt"
	client "gitlab.utc.fr/ia04_td1_group_j/ia04_td_vote/multiagentvote/clientag"
	"time"
)

// type RestClientAgent_NewBallot struct {
// 	id        string    // Identifiant ou unique représentant l'agent
// 	url       string    // URL du service REST que l'agent va consommer
// 	rule      string    // Règle de vote à appliquer
// 	deadline  time.Time // Date limite pour le vote
// 	voter_ids []string  // Identifiants des électeurs autorisés à participer au vote
// 	num_alts  int       // Nombre de candidats en lice
// 	tie_break []int     // Méthode pour départager en cas d'égalité
// }

func main() {
	ddl := "2023-12-30T23:05:08+04:00"
	t, _ := time.Parse(time.RFC3339, ddl)
	ag := client.NewRestClientAgent_NewBallot("id1", "http://localhost:8080", "borda", t, []string{"ag_id1", "ag_id2", "ag_id3"}, 2, []int{2, 1})
	ag.Start()
	fmt.Scanln()
}
