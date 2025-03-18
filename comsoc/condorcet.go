package comsoc

//chaque alt 1 fois
func CondorcetWinner(p Profile) (bestAlts []Alternative, err error) {
	for _, candidatA := range p[0] {
		winner := true
		for _, candidatB := range p[0] {
			if candidatA == candidatB {
				continue
			}
			timeWin := 0
			timeLose := 0
			for _, pref := range p {
				if isPref(candidatA, candidatB, pref) {
					timeWin++
				} else {
					timeLose++
				}
			}
			if timeWin > timeLose {
				continue
			} else {
				winner = false
				break
			}
		}
		if winner {
			bestAlts = append(bestAlts, candidatA) // only one or nil
		}
	}
	return bestAlts, nil
}
