package comsoc

import (
	"errors"
	"sort"
)

// renvoie l'indice ou se trouve alt dans prefs
func rank(alt Alternative, prefs []Alternative) int {
	for i, a := range prefs {
		if alt == a {
			return i
		}
	}
	return -1
}

// renvoie vrai ssi alt1 est préférée à alt2
func isPref(alt1, alt2 Alternative, prefs []Alternative) bool {
	p1 := -1
	p2 := -1
	for i, a := range prefs {
		if alt1 == a {
			p1 = i
		}
		if alt2 == a {
			p2 = i
		}
	}
	if p1 == -1 {
		return false
	} else if p2 == -1 {
		return true
	} else {
		return p1 < p2
	}
}

// renvoie les meilleures alternatives pour un décomtpe donné
func maxCountList(count Count) (bestAlts []Alternative) {
	var keyValuePaires []struct {
		key Alternative
		val int
	}
	for k, v := range count {
		keyValuePaires = append(keyValuePaires, struct {
			key Alternative
			val int
		}{k, v})
	}
	sort.Slice(keyValuePaires, func(i, j int) bool {
		return keyValuePaires[i].val > keyValuePaires[j].val
	})
	for _, a := range keyValuePaires {
		bestAlts = append(bestAlts, a.key)
	}
	return bestAlts
}

func maxCount(counts map[Alternative]int) (bestAlts []Alternative) {
	var BestA []Alternative
	count := make(map[Alternative]int)
	for alt, sco := range counts {
		count[alt] = sco
	}
	max := 0
	for len(count) != 0 {
		for alternat, score := range count {
			if score > max {
				max = score
				BestA = []Alternative{alternat}
				delete(count, alternat)
			} else if score == max {
				BestA = append(BestA, alternat)
				delete(count, alternat)
			} else {
				delete(count, alternat)
			}
		}
	}
	return BestA
}

func minCount(counts map[Alternative]int) (bestAlts []Alternative) {
	var BestA []Alternative
	count := make(map[Alternative]int)
	for alt, sco := range counts {
		count[alt] = sco
	}
	var min int
	for _, v := range counts {
		min = v
		break
	}
	for len(count) != 0 {
		for alternat, score := range count {
			if score < min {
				min = score
				BestA = []Alternative{alternat}
				delete(count, alternat)
			} else if score == min {
				BestA = append(BestA, alternat)
				delete(count, alternat)
			} else {
				delete(count, alternat)
			}
		}
	}
	return BestA
}

func unique(sl []Alternative) bool {
	numCount := make(map[Alternative]int)
	for _, num := range sl {
		if numCount[num] == 1 {
			return false
		}
		numCount[num] = 1
	}
	return true
}

// vérifie le profil donné, par ex. qu'ils sont tous complets et que chaque alternative n'apparaît qu'une seule fois par préférences
func checkProfile(prefs []Alternative, alts []Alternative) error {
	/*
		if len(prefs) != len(alts) {
			return errors.New("preference size not equal")
		}
	*/
	if !(unique(prefs)) {
		return errors.New("not unique")
	}

	for _, p := range prefs {
		ok := false
		for _, q := range alts {
			if p == q {
				ok = true
			}
		}
		if !ok {
			return errors.New("alt not found")
		}
	}
	return nil
}

// vérifie le profil donné, par ex. qu'ils sont tous complets et que chaque alternative de alts apparaît exactement une fois par préférences
func checkProfileAlternative(prefs Profile, alts []Alternative) error {
	expectedLen := len(alts)
	for _, p := range prefs {
		if len(p) != expectedLen {
			return errors.New("preference size not equal")
		}
		if !(unique(p)) {
			return errors.New("not unique")
		}
		for _, pp := range p {
			ok := false
			for _, q := range alts {
				if pp == q {
					ok = true
				}
			}
			if !ok {
				return errors.New("alt not found")
			}

		}
	}

	return nil
}

/* func remove(sl []Alternative, alt Alternative) []Alternative {
	var i int
	for i = range sl {
		if sl[i] == alt {
			break
		}
	}
	if i == len(sl) {
		return sl
	}

	return append(sl[:i], sl[i+1:]...)
}
*/

// DeleteSlice2 删除指定元素。
func remove(a []Alternative, elem Alternative) []Alternative {
	tmp := make([]Alternative, 0, len(a))
	for _, v := range a {
		if v != elem {
			tmp = append(tmp, v)
		}
	}
	return tmp
}
