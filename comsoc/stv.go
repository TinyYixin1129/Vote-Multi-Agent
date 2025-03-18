package comsoc

func STV_SWF(p Profile) (Count, error) {
	count := make(Count)
	score := make(Count)
	pp := make(Profile, len(p))
	copy(pp, p)
	for i := 0; i < len(p[0]); i++ {
		for j := 0; j < len(pp[0]); j++ {
			score[pp[0][j]] = 0
		}
		for _, pref := range pp {
			score[pref[0]]++
		}
		minalts := minCount(score)
		for _, elimate := range minalts {
			count[elimate] = i
			for j := range pp {
				pp[j] = remove(pp[j], elimate)
			}
			delete(score, elimate)
		}
		i += (len(minalts) - 1)
	}
	return count, nil
}

func STV_SCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := STV_SWF(p)
	if err != nil {
		return nil, err
	}
	return maxCount(count), nil
}
