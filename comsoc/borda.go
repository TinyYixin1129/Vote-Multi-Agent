package comsoc

func BordaSWF(p Profile) (count Count, err error) {
	count = make(Count)
	for _, pref := range p {
		for i, alt := range pref {
			if _, ok := count[alt]; !ok {
				count[alt] = len(pref) - i-1
			} else {
				count[alt] += len(pref) - i-1
			}
		}
	}
	return count, nil
}

func BordaSCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := BordaSWF(p)
	bestAlts = maxCount(count)
	return bestAlts, err
}
