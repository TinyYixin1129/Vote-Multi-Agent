package clientag

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	multiagentvote "gitlab.utc.fr/ia04_td1_group_j/ia04_td_vote/multiagentvote"
	"net/http"
)

type RestClientAgent_Result struct {
	id        string //Identifiant ou unique de l'agent
	url       string //URL du service REST que l'agent va interroger
	ballot_id string // Identifiant du bulletin de vote
}

func NewRestClientAgent_Result(id string, url string, ballot_id string) *RestClientAgent_Result {
	return &RestClientAgent_Result{id, url, ballot_id}
}

func (rca *RestClientAgent_Result) treatResponse(r *http.Response) (int, []int) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	var resp multiagentvote.Result_Res
	json.Unmarshal(buf.Bytes(), &resp)

	return resp.Winner, resp.Ranking
}

func (rca *RestClientAgent_Result) doRequest() (winner int, ranking []int, err error) {
	req := multiagentvote.Result_Req{
		Ballot_id: rca.ballot_id,
	}

	// sérialisation de la requête
	url := rca.url + "/result"
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
	winner, ranking = rca.treatResponse(resp)

	return
}

func (rca *RestClientAgent_Result) Start() {
	log.Printf("démarrage de %s", rca.id)
	winner, ranking, err := rca.doRequest()
	if err != nil {
		log.Printf(rca.id, "error:", err.Error())
	} else {
		log.Printf("[%s] ballot_id: %s , Winner: %d , Ranking: %d \n", rca.id, rca.ballot_id, winner, ranking)
	}
}
