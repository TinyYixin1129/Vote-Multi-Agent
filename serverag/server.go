package serverag

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	multiagentvote "gitlab.utc.fr/ia04_td1_group_j/ia04_td_vote/multiagentvote"
	comsoc "gitlab.utc.fr/ia04_td1_group_j/ia04_td_vote/multiagentvote/comsoc"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Stocker les informations complètes d'un ballot
type Ballot struct {
	Ballot_id string    //"scrutin12"
	Rule      string    //"majority","borda", "approval", "stv", "kemeny"
	Deadline  time.Time //"2023-10-09T23:05:08+02:00"  (format RFC 3339)
	Voter_ids []string  //["ag_id1", "ag_id2", "ag_id3"]
	Num_alts  int       //Nombre de Alternative
	Tie_break []int     //[4, 2, 3, 5, 9, 8, 7, 1, 6, 11, 12, 10]
	Options   []int
	pro       comsoc.Profile //[12][0] Alternative_1
	winner    int            //4 SCF
	Ranking   []int          //[2, 1, 4, 3]  SWF
}

type RestServerAgent struct {
	sync.Mutex
	id          string
	reqCount    int //Le nombre de requêtes au serveur
	addr        string
	ballot_list []Ballot //Tableau qui enregistre tous les bulletins de vote
}

func NewRestServerAgent(addr string) *RestServerAgent {
	return &RestServerAgent{id: addr, addr: addr}
}

// Test de la méthode
func (rsa *RestServerAgent) checkMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method { // method : POST GET
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "method %q not allowed", r.Method)
		return false
	}
	return true
}

func (*RestServerAgent) decodeBallot_Req(r *http.Request) (req multiagentvote.Ballot_Req, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

func (*RestServerAgent) decodeVote_Req(r *http.Request) (req multiagentvote.Vote_Req, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

func (*RestServerAgent) decodeResult_Req(r *http.Request) (req multiagentvote.Result_Req, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

func (rsa *RestServerAgent) doCreatballot(w http.ResponseWriter, r *http.Request) {
	rsa.Lock()
	defer rsa.Unlock()
	rsa.reqCount++

	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := rsa.decodeBallot_Req(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) //400 bad request
		fmt.Fprint(w, err.Error())
		return
	}

	var resp multiagentvote.Ballot_Res
	ballot_num := len(rsa.ballot_list) + 1
	id := "scrutin" + strconv.Itoa(ballot_num)

	switch req.Rule {
	case "majority", "borda", "approval", "condorcet", "copeland", "stv":
		if len(req.Tie_break) != 0 {
			if req.Num_alts != len(req.Tie_break) { //Lorsque le nombre de candidats est différent du nombre de départagements
				w.WriteHeader(http.StatusNotImplemented) //501 not implemented
				msg := "Le nombre d'alts est mauvais"
				w.Write([]byte(msg))
				return
			}
		}
		currentTime := time.Now()
		if currentTime.After(req.Deadline) { //Déterminer si la date limite de saisie correspond à la fois précédente
			w.WriteHeader(http.StatusNotImplemented) //501 not implemented
			msg := "Le valeur de Deadline est mauvais"
			w.Write([]byte(msg))
			return
		}
		if len(req.Voter_ids) == 0 { //Voter_ids vide
			w.WriteHeader(http.StatusNotImplemented) //501 not implemented
			msg := "Le list de Voter_ids  est vide"
			w.Write([]byte(msg))
			return
		}
		new_ballot := Ballot{
			Ballot_id: id,
			Rule:      req.Rule,
			Deadline:  req.Deadline,
			Voter_ids: req.Voter_ids,
			Num_alts:  req.Num_alts,
			Tie_break: req.Tie_break,
		}
		rsa.ballot_list = append(rsa.ballot_list, new_ballot) //Enregistrez la bulle nouvellement créée dans la liste
		resp.Ballot_id = id
	default:
		w.WriteHeader(http.StatusNotImplemented) //501 not implemented
		msg := fmt.Sprintf("Unkonwn command '%s'", req.Rule)
		w.Write([]byte(msg))
		return
	}

	w.WriteHeader(http.StatusOK)
	serial, _ := json.Marshal(resp)
	//Convertir une structure de données ou en données d'octets au format JSON
	w.Write(serial)
}

func (rsa *RestServerAgent) doVote(w http.ResponseWriter, r *http.Request) {
	// mise à jour du nombre de requêtes
	rsa.Lock()
	defer rsa.Unlock()
	rsa.reqCount++

	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := rsa.decodeVote_Req(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) //400 bad request
		fmt.Fprint(w, err.Error())
		return
	}

	if len(req.Prefs) == 0 {
		w.WriteHeader(http.StatusBadRequest) //400 bad request
		fmt.Fprint(w, err.Error())
		return
	}
	// traitement de la requête
	bal_find := 0
	ag_find := 0
	for i, ballot := range rsa.ballot_list { //Parcourez la liste pour trouver les votes requis
		if req.Ballot_id == ballot.Ballot_id {
			bal_find = 1
			currentTime := time.Now()
			//ddl, _ := time.Parse(time.RFC3339, rsa.ballot_list[i].Deadline)
			ddl := rsa.ballot_list[i].Deadline
			if currentTime.After(ddl) { // Date limite dépassée
				w.WriteHeader(http.StatusServiceUnavailable) //503 Service Unavailable
				msg := "la deadline est dépassée"
				w.Write([]byte(msg))
				return
			}

			for j, ag_id := range ballot.Voter_ids {
				if req.Agent_id == ag_id {
					ag_find = 1
					var voter []comsoc.Alternative
					for _, value := range req.Prefs {
						//rsa.ballot_list[i].pro[j][index]=Alternative(value)
						voter = append(voter, comsoc.Alternative(value))
					}
					rsa.ballot_list[i].pro = append(rsa.ballot_list[i].pro, voter)                  //传入prefs
					rsa.ballot_list[i].Options = append(rsa.ballot_list[i].Options, req.Options...) //传入options

					var orderedAlts []comsoc.Alternative
					for _, value := range ballot.Tie_break {

						orderedAlts = append(orderedAlts, comsoc.Alternative(value))
					}

					switch ballot.Rule {
					case "majority":
						ranking, err1 := comsoc.SWFFactory(comsoc.MajoritySWF, comsoc.TieBreakFactory(orderedAlts))(rsa.ballot_list[i].pro)
						if err1 != nil {
							w.WriteHeader(http.StatusNotImplemented) //501 not implemented
							msg := "MajoritySWF error "
							w.Write([]byte(msg))
							//删除本次voter放入的数据
							rsa.ballot_list[i].pro = rsa.ballot_list[i].pro[:len(rsa.ballot_list[i].pro)-1]
							rsa.ballot_list[i].Options = rsa.ballot_list[i].Options[:len(rsa.ballot_list[i].Options)-1]
							return
						}
						ranking, err3 := comsoc.TieBreak4Worst(orderedAlts)(ranking)
						if err3 != nil {

							w.WriteHeader(http.StatusNotImplemented) //501 not implemented
							msg := "MajoritySWF_TieBreak4Worst error "
							w.Write([]byte(msg))
							//删除本次voter放入的数据
							rsa.ballot_list[i].pro = rsa.ballot_list[i].pro[:len(rsa.ballot_list[i].pro)-1]
							rsa.ballot_list[i].Options = rsa.ballot_list[i].Options[:len(rsa.ballot_list[i].Options)-1]
							return
						}
						best, err2 := comsoc.SCFFactory(comsoc.MajoritySCF, comsoc.TieBreakFactory(orderedAlts))(rsa.ballot_list[i].pro)
						if err2 != nil {
							w.WriteHeader(http.StatusNotImplemented) //501 not implemented
							msg := "MajoritySCF error "
							w.Write([]byte(msg))
							//删除本次voter放入的数据
							rsa.ballot_list[i].pro = rsa.ballot_list[i].pro[:len(rsa.ballot_list[i].pro)-1]
							rsa.ballot_list[i].Options = rsa.ballot_list[i].Options[:len(rsa.ballot_list[i].Options)-1]
							return
						}
						rsa.ballot_list[i].winner = 0
						rsa.ballot_list[i].winner = int(best)
						rsa.ballot_list[i].Ranking = nil
						for _, value := range ranking {

							rsa.ballot_list[i].Ranking = append(rsa.ballot_list[i].Ranking, int(value))
						}
					case "borda":
						ranking, err1 := comsoc.SWFFactory(comsoc.BordaSWF, comsoc.TieBreakFactory(orderedAlts))(rsa.ballot_list[i].pro)
						if err1 != nil {
							w.WriteHeader(http.StatusNotImplemented) //501 not implemented
							msg := "BordaSWF error "
							w.Write([]byte(msg))
							//删除本次voter放入的数据
							rsa.ballot_list[i].pro = rsa.ballot_list[i].pro[:len(rsa.ballot_list[i].pro)-1]
							rsa.ballot_list[i].Options = rsa.ballot_list[i].Options[:len(rsa.ballot_list[i].Options)-1]
							return
						}
						ranking, err3 := comsoc.TieBreak4Worst(orderedAlts)(ranking)
						if err3 != nil {

							w.WriteHeader(http.StatusNotImplemented) //501 not implemented
							msg := "BordaSWF_TieBreak4Worst error "
							w.Write([]byte(msg))
							//删除本次voter放入的数据
							rsa.ballot_list[i].pro = rsa.ballot_list[i].pro[:len(rsa.ballot_list[i].pro)-1]
							rsa.ballot_list[i].Options = rsa.ballot_list[i].Options[:len(rsa.ballot_list[i].Options)-1]
							return
						}
						best, err2 := comsoc.SCFFactory(comsoc.BordaSCF, comsoc.TieBreakFactory(orderedAlts))(rsa.ballot_list[i].pro)
						if err2 != nil {
							w.WriteHeader(http.StatusNotImplemented) //501 not implemented
							msg := "BordaSCF error "
							w.Write([]byte(msg))
							//删除本次voter放入的数据
							rsa.ballot_list[i].pro = rsa.ballot_list[i].pro[:len(rsa.ballot_list[i].pro)-1]
							rsa.ballot_list[i].Options = rsa.ballot_list[i].Options[:len(rsa.ballot_list[i].Options)-1]
							return
						}
						rsa.ballot_list[i].winner = 0
						rsa.ballot_list[i].winner = int(best)
						rsa.ballot_list[i].Ranking = nil
						for _, value := range ranking {

							rsa.ballot_list[i].Ranking = append(rsa.ballot_list[i].Ranking, int(value))
						}
					case "approval":
						ranking, err1 := comsoc.SWFFactory(comsoc.ApprovalSWF(rsa.ballot_list[i].Options), comsoc.TieBreakFactory(orderedAlts))(rsa.ballot_list[i].pro)
						if err1 != nil {
							w.WriteHeader(http.StatusNotImplemented) //501 not implemented
							msg := "ApprovalSWF error "
							w.Write([]byte(msg))
							//删除本次voter放入的数据
							rsa.ballot_list[i].pro = rsa.ballot_list[i].pro[:len(rsa.ballot_list[i].pro)-1]
							rsa.ballot_list[i].Options = rsa.ballot_list[i].Options[:len(rsa.ballot_list[i].Options)-1]
							return
						}
						ranking, err3 := comsoc.TieBreak4Worst(orderedAlts)(ranking)
						if err3 != nil {

							w.WriteHeader(http.StatusNotImplemented) //501 not implemented
							msg := "ApprovalSWF_TieBreak4Worst error "
							w.Write([]byte(msg))
							//删除本次voter放入的数据
							rsa.ballot_list[i].pro = rsa.ballot_list[i].pro[:len(rsa.ballot_list[i].pro)-1]
							rsa.ballot_list[i].Options = rsa.ballot_list[i].Options[:len(rsa.ballot_list[i].Options)-1]
							return
						}
						best, err2 := comsoc.SCFFactory(comsoc.ApprovalSCF(rsa.ballot_list[i].Options), comsoc.TieBreakFactory(orderedAlts))(rsa.ballot_list[i].pro)
						if err2 != nil {
							w.WriteHeader(http.StatusNotImplemented) //501 not implemented
							msg := "ApprovalSCF error "
							w.Write([]byte(msg))
							//删除本次voter放入的数据
							rsa.ballot_list[i].pro = rsa.ballot_list[i].pro[:len(rsa.ballot_list[i].pro)-1]
							rsa.ballot_list[i].Options = rsa.ballot_list[i].Options[:len(rsa.ballot_list[i].Options)-1]
							return
						}
						rsa.ballot_list[i].winner = 0
						rsa.ballot_list[i].winner = int(best)
						rsa.ballot_list[i].Ranking = nil
						for _, value := range ranking {

							rsa.ballot_list[i].Ranking = append(rsa.ballot_list[i].Ranking, int(value))
						}
					case "copeland":
						ranking, err1 := comsoc.SWFFactory(comsoc.CopelandSWF, comsoc.TieBreakFactory(orderedAlts))(rsa.ballot_list[i].pro)
						if err1 != nil {
							w.WriteHeader(http.StatusNotImplemented) //501 not implemented
							msg := "CopelandSWF error "
							w.Write([]byte(msg))
							//删除本次voter放入的数据
							rsa.ballot_list[i].pro = rsa.ballot_list[i].pro[:len(rsa.ballot_list[i].pro)-1]
							rsa.ballot_list[i].Options = rsa.ballot_list[i].Options[:len(rsa.ballot_list[i].Options)-1]
							return
						}
						ranking, err3 := comsoc.TieBreak4Worst(orderedAlts)(ranking)
						if err3 != nil {

							w.WriteHeader(http.StatusNotImplemented) //501 not implemented
							msg := "CopelandSWF_TieBreak4Worst error "
							w.Write([]byte(msg))
							//删除本次voter放入的数据
							rsa.ballot_list[i].pro = rsa.ballot_list[i].pro[:len(rsa.ballot_list[i].pro)-1]
							rsa.ballot_list[i].Options = rsa.ballot_list[i].Options[:len(rsa.ballot_list[i].Options)-1]
							return
						}
						best, err2 := comsoc.SCFFactory(comsoc.CopelandSCF, comsoc.TieBreakFactory(orderedAlts))(rsa.ballot_list[i].pro)
						if err2 != nil {
							w.WriteHeader(http.StatusNotImplemented) //501 not implemented
							msg := "CopelandSCF error "
							w.Write([]byte(msg))
							//删除本次voter放入的数据
							rsa.ballot_list[i].pro = rsa.ballot_list[i].pro[:len(rsa.ballot_list[i].pro)-1]
							rsa.ballot_list[i].Options = rsa.ballot_list[i].Options[:len(rsa.ballot_list[i].Options)-1]
							return
						}
						rsa.ballot_list[i].winner = 0
						rsa.ballot_list[i].winner = int(best)
						rsa.ballot_list[i].Ranking = nil
						for _, value := range ranking {

							rsa.ballot_list[i].Ranking = append(rsa.ballot_list[i].Ranking, int(value))
						}
					case "stv":
						ranking, err1 := comsoc.SWFFactory(comsoc.STV_SWF, comsoc.TieBreakFactory(orderedAlts))(rsa.ballot_list[i].pro)
						if err1 != nil {
							w.WriteHeader(http.StatusNotImplemented) //501 not implemented
							msg := "STV_SWF error "
							w.Write([]byte(msg))
							//删除本次voter放入的数据
							rsa.ballot_list[i].pro = rsa.ballot_list[i].pro[:len(rsa.ballot_list[i].pro)-1]
							rsa.ballot_list[i].Options = rsa.ballot_list[i].Options[:len(rsa.ballot_list[i].Options)-1]
							return
						}
						ranking, err3 := comsoc.TieBreak4Worst(orderedAlts)(ranking)
						if err3 != nil {

							w.WriteHeader(http.StatusNotImplemented) //501 not implemented
							msg := "STV_SWF_TieBreak4Worst error "
							w.Write([]byte(msg))
							//删除本次voter放入的数据
							rsa.ballot_list[i].pro = rsa.ballot_list[i].pro[:len(rsa.ballot_list[i].pro)-1]
							rsa.ballot_list[i].Options = rsa.ballot_list[i].Options[:len(rsa.ballot_list[i].Options)-1]
							return
						}
						best, err2 := comsoc.SCFFactory(comsoc.STV_SCF, comsoc.TieBreakFactory(orderedAlts))(rsa.ballot_list[i].pro)
						if err2 != nil {
							w.WriteHeader(http.StatusNotImplemented) //501 not implemented
							msg := "STV_SCF error "
							w.Write([]byte(msg))
							//删除本次voter放入的数据
							rsa.ballot_list[i].pro = rsa.ballot_list[i].pro[:len(rsa.ballot_list[i].pro)-1]
							rsa.ballot_list[i].Options = rsa.ballot_list[i].Options[:len(rsa.ballot_list[i].Options)-1]
							return
						}
						rsa.ballot_list[i].winner = 0
						rsa.ballot_list[i].winner = int(best)
						rsa.ballot_list[i].Ranking = nil
						for _, value := range ranking {

							rsa.ballot_list[i].Ranking = append(rsa.ballot_list[i].Ranking, int(value))
						}
					case "condorcet":
						best, err := comsoc.CondorcetWinner(rsa.ballot_list[i].pro) // return only one alt or nil
						if err != nil {
							w.WriteHeader(http.StatusNotImplemented) //501 not implemented
							msg := "CondorcetWinner error "
							w.Write([]byte(msg))
							//删除本次voter放入的数据
							rsa.ballot_list[i].pro = rsa.ballot_list[i].pro[:len(rsa.ballot_list[i].pro)-1]
							rsa.ballot_list[i].Options = rsa.ballot_list[i].Options[:len(rsa.ballot_list[i].Options)-1]
							return
						}
						rsa.ballot_list[i].winner = 0
						if best != nil {
							rsa.ballot_list[i].winner = int(best[0])
						}
					default:
						w.WriteHeader(http.StatusNotImplemented)
						msg := fmt.Sprintf("Il y a un problème avec la méthode de vote enregistré '%s'", ballot.Rule)
						w.Write([]byte(msg))
						return
					}

					//Les électeurs ayant déjà voté sont rayés de la liste
					rsa.ballot_list[i].Voter_ids = append(rsa.ballot_list[i].Voter_ids[:j], rsa.ballot_list[i].Voter_ids[j+1:]...)
					break
				}
			}
			break
		}
	}

	if bal_find == 0 { //L'identifiant du bulletin de vote n'existe pas
		w.WriteHeader(http.StatusNotImplemented) //501 not implemented
		msg := fmt.Sprintf("Le Ballot '%s' n'existe pas ", req.Ballot_id)
		w.Write([]byte(msg))
		return
	}
	if ag_find == 0 { // l'identifiant de l'électeur n'existe pas
		w.WriteHeader(http.StatusForbidden) //403 Forbidden
		msg := fmt.Sprintf(" '%s' vote déjà effectué ou Vous n'êtes pas sur la liste de vote", req.Agent_id)
		w.Write([]byte(msg))
		return
	}

	w.WriteHeader(http.StatusOK)
	msg := "vote pris en compte "
	w.Write([]byte(msg))
}

func (rsa *RestServerAgent) doGetresult(w http.ResponseWriter, r *http.Request) {

	rsa.Lock()
	defer rsa.Unlock()
	rsa.reqCount++

	if !rsa.checkMethod("POST", w, r) {
		return
	}
	// décodage de la requête
	req, err := rsa.decodeResult_Req(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) //400 bad request
		fmt.Fprint(w, err.Error())
		return
	}
	bal_find := 0
	var resp multiagentvote.Result_Res
	for _, ballot := range rsa.ballot_list {
		if req.Ballot_id == ballot.Ballot_id { //Trouvez le vote que vous souhaitez interroger
			bal_find = 1
			//Déterminer si le vote est terminé
			currentTime := time.Now()
			//ddl, _ := time.Parse(time.RFC3339, ballot.Deadline)
			if currentTime.Before(ballot.Deadline) {
				w.WriteHeader(http.StatusPreconditionRequired) //425 Too early
				msg := "Too early"
				w.Write([]byte(msg))
				return
			}
			resp.Winner = ballot.winner
			resp.Ranking = ballot.Ranking
			break
		}
	}

	if bal_find == 0 {
		w.WriteHeader(http.StatusNotFound) //501 not implemented
		msg := fmt.Sprintf("Not Found: Le Ballot '%s' n'existe pas ", req.Ballot_id)
		w.Write([]byte(msg))
		return
	}

	w.WriteHeader(http.StatusOK)
	if resp.Winner == 0 {
		msg := fmt.Sprintf("No Winner for Le Ballot '%s' ", req.Ballot_id)
		w.Write([]byte(msg))
		return
	} else if len(resp.Ranking) == 0 { // CondorcetWinner , no ranking
		msg := fmt.Sprintf("CondorcetWinner: '%d' ", resp.Winner)
		w.Write([]byte(msg))
		return
	} else {
		serial, _ := json.Marshal(resp)
		w.Write(serial)
	}

}

func (rsa *RestServerAgent) Start() {
	// création du multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/new_ballot", rsa.doCreatballot)
	mux.HandleFunc("/vote", rsa.doVote)
	mux.HandleFunc("/result", rsa.doGetresult)

	// création du serveur http
	s := &http.Server{
		Addr:           rsa.addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second, //Délai d'expiration de lecture et d'écriture
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20} //Nombre maximum d'octets d'en-tête

	// lancement du serveur
	log.Println("Listening on", rsa.addr)
	go log.Fatal(s.ListenAndServe())
}
