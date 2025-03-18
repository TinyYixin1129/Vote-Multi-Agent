package clientag

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	multiagentvote "gitlab.utc.fr/ia04_td1_group_j/ia04_td_vote/multiagentvote"
	"net/http"
)

type RestClientAgent_Vote struct {
	id        string // Identifiant ou unique de l'agent
	url       string // URL du service REST que l'agent va utiliser
	agend_id  string // Nom de l'électeur
	ballot_id string // Identifiant du bulletin de vote
	prefs     []int  // Préférences de vote
	options   []int  // Options de vote
}

func NewRestClientAgent_Vote(id string, url string, agend_id string, ballot_id string, prefs []int, options []int) *RestClientAgent_Vote {
	return &RestClientAgent_Vote{id, url, agend_id, ballot_id, prefs, options}
}

// func (rca *RestClientAgent_Vote) treatResponse(r *http.Response) {
// 	buf := new(bytes.Buffer)
// 	buf.ReadFrom(r.Body)
// 	// ???
// 	var resp rad.Response
// 	json.Unmarshal(buf.Bytes(), &resp)
// 	return resp.Result
// }

func (rca *RestClientAgent_Vote) doRequest() (err error) {
	req := multiagentvote.Vote_Req{
		Agent_id:  rca.agend_id,
		Ballot_id: rca.ballot_id,
		Prefs:     rca.prefs,
		Options:   rca.options,
	}

	// sérialisation de la requête
	url := rca.url + "/vote"
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
	//res = rca.treatResponse(resp) //resp是http响应的类

	return
}

func (rca *RestClientAgent_Vote) Start() {
	log.Printf("démarrage de %s", rca.id)
	err := rca.doRequest()

	if err != nil {
		log.Printf(rca.id, "error:", err.Error()) //log.Fatal
	} else {
		log.Printf("[%s] Vote successfully, scrutin_id: %s\n", rca.id, rca.ballot_id)
	}
}
