package comsoc

import (
	"errors"
)

// renvoie toujours le premier
func TieBreak(alts []Alternative) (Alternative, error) {
	if len(alts) == 0 {
		return -1, errors.New("alternative est vide")
	}
	return alts[0], nil
}

func TieBreakFactory(orderedAlts []Alternative) func([]Alternative) (Alternative, error) {
	tiebreak := func(TieAlts []Alternative) (Alternative, error) {
		if len(TieAlts) == 0 {
			return -1, errors.New("bestAlts vide")
		}
		result := TieAlts[0]
		for _, a := range TieAlts {
			if rank(a, orderedAlts) < rank(result, orderedAlts) && rank(a, orderedAlts) > -1 {
				result = a
			}
		}
		return result, nil
	}
	return tiebreak
}

func SWFFactory(swf func(p Profile) (Count, error), tiebreak func([]Alternative) (Alternative, error)) func(Profile) ([]Alternative, error) {
	SWFfunc := func(pp Profile) ([]Alternative, error) {
		count, err := swf(pp)
		if err != nil {
			return nil, err
		}
		var order_result []Alternative
		for {
			bestAlts := maxCount(count)
			if len(bestAlts) == 1 {
				order_result = append(order_result, bestAlts[0])
				delete(count, bestAlts[0])
			} else if len(bestAlts) > 1 {
				for len(bestAlts) > 0 {
					alt, err := tiebreak(bestAlts)
					if err != nil {
						return nil, err
					}
					order_result = append(order_result, alt)
					for i := range bestAlts {
						if bestAlts[i] == alt {
							bestAlts = append(bestAlts[:i], bestAlts[i+1:]...)
							break
						}
					}
					delete(count, alt)
				}
			} else {
				//all done
				break
			}
		}
		return order_result, nil
	}
	return SWFfunc
}

func SCFFactory(scf func(p Profile) ([]Alternative, error), tiebreak func([]Alternative) (Alternative, error)) func(Profile) (Alternative, error) {
	SCFfunc := func(pp Profile) (Alternative, error) {
		bestAlts, err := scf(pp)
		if err != nil {
			return -1, err
		}
		alt_result, err := tiebreak(bestAlts)
		if err != nil {
			return -1, err
		}
		return alt_result, nil
	}
	return SCFfunc
}

func TieBreak4Worst(orderedAlts []Alternative) func([]Alternative) ([]Alternative, error) { // 把没有票数的特殊情况也加入到最终的ranking中
	tiebreak := func(notWorstAlts []Alternative) ([]Alternative, error) {
		cp := make([]Alternative, len(orderedAlts))
		copy(cp, orderedAlts)
		for _, alt := range notWorstAlts {
			cp = remove(cp, alt)
		}
		result := append(notWorstAlts, cp[0:]...)
		return result, nil
	}
	return tiebreak
}
