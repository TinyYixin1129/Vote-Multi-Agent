package main

import (
	"fmt"
	client "gitlab.utc.fr/ia04_td1_group_j/ia04_td_vote/multiagentvote/clientag"
)

// type RestClientAgent_Vote struct {
// 	id        string // Identifiant ou unique de l'agent
// 	url       string // URL du service REST que l'agent va utiliser
// 	agend_id  string // Nom de l'électeur
// 	ballot_id string // Identifiant du bulletin de vote
// 	prefs     []int  // Préférences de vote
// 	options   []int  // Options de vote
// }

func main() {
	ag := client.NewRestClientAgent_Vote("id3", "http://localhost:8080", "ag_id1", "scrutin1", []int{4, 2, 3, 1}, []int{1})
	ag.Start()
	fmt.Scanln()
}
