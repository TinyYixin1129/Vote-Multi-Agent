package comsoc

import "errors"

/*
func ApprovalSWF(p Profile, thresholds []int) (count Count, err error) {
	count = make(Count)
	for j, pref := range p {
		for i, alt := range pref {
			if i < thresholds[j]-1 {
				if _, ok := count[alt]; !ok {
					count[alt] = 1
				} else {
					count[alt] += 1
				}
			}
		}
	}
	return count, nil
}

func ApprovalSCF(p Profile, thresholds []int) (bestAlts []Alternative, err error) {
	count, err := ApprovalSWF(p, thresholds)
	bestAlts = maxCount(count)
	return bestAlts, err
}
*/

func ApprovalSWF(thresholds []int) func(Profile) (Count, error) {
	approval := func(p Profile) (Count, error) {
		count := make(Count)
		if thresholds == nil {
			return nil, errors.New("thresholds vide")
		}

		for j, pref := range p {
			for i, alt := range pref {
				if i < thresholds[j] {
					if _, ok := count[alt]; !ok {
						count[alt] = 1
					} else {
						count[alt] += 1
					}
				}
			}
		}
		return count, nil
	}
	return approval
}

func ApprovalSCF(thresholds []int) func(Profile) ([]Alternative, error) {
	approval := func(p Profile) (bestAlts []Alternative, err error) {
		count, err := ApprovalSWF(thresholds)(p)
		bestAlts = maxCount(count)
		return bestAlts, err
	}
	return approval
}
