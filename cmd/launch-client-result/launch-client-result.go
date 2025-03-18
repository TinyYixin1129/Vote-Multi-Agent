package main

import (
	"fmt"
	client "gitlab.utc.fr/ia04_td1_group_j/ia04_td_vote/multiagentvote/clientag"
)

// type RestClientAgent_Result struct {
// 	id        string //Identifiant ou unique de l'agent
// 	url       string //URL du service REST que l'agent va interroger
// 	ballot_id string // Identifiant du bulletin de vote
// }

func main() {
	ag := client.NewRestClientAgent_Result("id2", "http://localhost:8080", "scrutin1")
	ag.Start()
	fmt.Scanln()
}
