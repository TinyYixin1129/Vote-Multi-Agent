package comsoc

//only count the first
func MajoritySWF(p Profile) (count Count, err error) {
	count = make(Count)
	for _, preferences := range p {
		if _, ok := count[preferences[0]]; !ok {
			count[preferences[0]] = 1
		} else {
			count[preferences[0]] = count[preferences[0]] + 1
		}
	}
	return count, err
}

func MajoritySCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := MajoritySWF(p)
	bestAlts = maxCount(count)
	return bestAlts, err
}
