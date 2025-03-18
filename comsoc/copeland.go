package comsoc

//chaque alt 1 fois
func CopelandSWF(p Profile) (Count, error) {
	count := make(Count)
	for i, candidatA := range p[0] {
		for j := i + 1; j < len(p[0]); j++ {
			candidatB := p[0][j]
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
				count[candidatA]++
				count[candidatB]--
			} else if timeWin < timeLose {
				count[candidatA]--
				count[candidatB]++
			}
		}
	}
	return count, nil
}

func CopelandSCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := CopelandSWF(p)
	if err != nil {
		return nil, err
	}
	return maxCount(count), nil
}
