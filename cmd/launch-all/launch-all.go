package main

import (
	"fmt"
	"log"
	"math/rand"
	//multiagentvote "gitlab.utc.fr/ia04_td1_group_j/ia04_td_vote/multiagentvote"
	client "gitlab.utc.fr/ia04_td1_group_j/ia04_td_vote/multiagentvote/clientag"
	server "gitlab.utc.fr/ia04_td1_group_j/ia04_td_vote/multiagentvote/serverag"
	"time"
)

func main() {
	const n = 10
	const url1 = ":8080"
	const url2 = "http://localhost:8080"
	rule := [...]string{"majority", "borda", "approval", "condorcet", "copeland", "stv"}

	clAgts_NewBallot := make([]client.RestClientAgent_NewBallot, 0, n)
	clAgts_Vote := make([]client.RestClientAgent_Vote, 0, n)
	clAgts_Result := make([]client.RestClientAgent_Result, 0, n)
	servAgt := server.NewRestServerAgent(url1)

	log.Println("démarrage du serveur...")

	go servAgt.Start()

	log.Println("démarrage des clients...")

	//creat clAgts_NewBallot
	ddl := "2023-12-30T23:05:08+04:00"
	t, _ := time.Parse(time.RFC3339, ddl)
	for i := 0; i < n; i++ {
		id := fmt.Sprintf("clAgts_NewBallot:id%02d", i)
		rule := rule[rand.Intn(6)]
		ag := client.NewRestClientAgent_NewBallot(id, url2, rule, t, []string{"ag_id1", "ag_id2", "ag_id3"}, 4, []int{2, 1, 4, 3})
		clAgts_NewBallot = append(clAgts_NewBallot, *ag)
	}

	//creat clAgts_Vote
	agt_list := [...]string{"ag_id1", "ag_id2", "ag_id3"}
	slice := []int{1, 2, 3, 4}
	for i := 0; i < n; i++ {
		rand.Shuffle(len(slice), func(i, j int) {
			slice[i], slice[j] = slice[j], slice[i]
		})
		id := fmt.Sprintf("clAgts_Vote:id%02d", i)
		agt := agt_list[rand.Intn(3)]
		sr_num := rand.Intn(11)
		scruntin := fmt.Sprintf("scrutin%d", sr_num)
		ag := client.NewRestClientAgent_Vote(id, url2, agt, scruntin, slice, []int{1})
		clAgts_Vote = append(clAgts_Vote, *ag)
	}
	//creat clAgts_Result
	for i := 0; i < n; i++ {
		id := fmt.Sprintf("clAgts_Result:id%02d", i)
		sr_num := rand.Intn(11)
		scruntin := fmt.Sprintf("scrutin%d", sr_num)
		ag := client.NewRestClientAgent_Result(id, url2, scruntin)
		clAgts_Result = append(clAgts_Result, *ag)
	}

	for _, agt := range clAgts_NewBallot {
		// attention, obligation de passer par cette lambda pour faire capturer la valeur de l'itération par la goroutine
		func(agt client.RestClientAgent_NewBallot) {
			go agt.Start()
		}(agt)
	}
	time.Sleep(2 * time.Second)
	for _, agt := range clAgts_Vote {
		// attention, obligation de passer par cette lambda pour faire capturer la valeur de l'itération par la goroutine
		func(agt client.RestClientAgent_Vote) {
			go agt.Start()
		}(agt)
	}
	time.Sleep(2 * time.Second)
	for _, agt := range clAgts_Result {
		// attention, obligation de passer par cette lambda pour faire capturer la valeur de l'itération par la goroutine
		func(agt client.RestClientAgent_Result) {
			go agt.Start()
		}(agt)
	}

	fmt.Scanln()
}
