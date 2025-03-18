package main

import (
	"fmt"
	server "gitlab.utc.fr/ia04_td1_group_j/ia04_td_vote/multiagentvote/serverag"
)

func main() {
	server := server.NewRestServerAgent(":8080")
	server.Start()
	fmt.Scanln()
}
