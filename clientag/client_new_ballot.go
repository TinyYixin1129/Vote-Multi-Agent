package clientag

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	multiagentvote "gitlab.utc.fr/ia04_td1_group_j/ia04_td_vote/multiagentvote"
	"net/http"
	"time"
)

type RestClientAgent_NewBallot struct {
	id        string    // Identifiant ou unique représentant l'agent
	url       string    // URL du service REST que l'agent va consommer
	rule      string    // Règle de vote à appliquer
	deadline  time.Time // Date limite pour le vote
	voter_ids []string  // Identifiants des électeurs autorisés à participer au vote
	num_alts  int       // Nombre de candidats en lice
	tie_break []int     // Méthode pour départager en cas d'égalité
}

func NewRestClientAgent_NewBallot(id string, url string, rule string, deadline time.Time, voter_ids []string, num_alts int, tie_break []int) *RestClientAgent_NewBallot {
	return &RestClientAgent_NewBallot{id, url, rule, deadline, voter_ids, num_alts, tie_break}
}

func (rca *RestClientAgent_NewBallot) treatResponse(r *http.Response) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	var resp multiagentvote.Ballot_Res
	json.Unmarshal(buf.Bytes(), &resp)

	return resp.Ballot_id + " rule: " + rca.rule
}

func (rca *RestClientAgent_NewBallot) doRequest() (res string, err error) {
	req := multiagentvote.Ballot_Req{
		Rule:      rca.rule,
		Deadline:  rca.deadline,
		Voter_ids: rca.voter_ids,
		Num_alts:  rca.num_alts,
		Tie_break: rca.tie_break,
	}

	// sérialisation de la requête
	url := rca.url + "/new_ballot"
	data, _ := json.Marshal(req)

	// envoi de la requête
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	// traitement de la réponse
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		return
	}
	res = rca.treatResponse(resp)

	return
}

func (rca *RestClientAgent_NewBallot) Start() {
	log.Printf("démarrage de %s", rca.id)
	Ballot_id, err := rca.doRequest()

	if err != nil {
		log.Printf(rca.id, "error:", err.Error())
	} else {
		log.Printf("Creat %s successfully\n", Ballot_id)
	}
}
