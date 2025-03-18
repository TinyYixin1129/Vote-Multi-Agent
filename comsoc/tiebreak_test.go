package comsoc

import (
	"fmt"
	"testing"
)

func TestSWFFactory(t *testing.T) {
	fmt.Println("\nSWFFactroy")
	var prof Profile
	voter1 := []Alternative{1, 2, 3, 4}
	voter2 := []Alternative{1, 3, 2, 4}
	//voter3 := []Alternative{1, 2, 3}
	//voter4 := []Alternative{2, 1, 5}
	orderedAlts := []Alternative{4, 2, 1, 3}
	prof = append(prof, voter1, voter2)
	get, err1 := SWFFactory(MajoritySWF, TieBreakFactory(orderedAlts))(prof)
	got, err2 := TieBreak4Worst(orderedAlts)(get)

	if err1 != nil {
		t.Errorf("error")
	}
	if err2 != nil {
		t.Errorf("error")
	}

	want := []Alternative{1, 4, 2, 3}

	for i, val := range got {
		if val != want[i] {
			t.Errorf("On a %d alors que l'on veut %d etant donne, %v", got, want, prof)
		}
	}
}

func TestSCFFactory(t *testing.T) {
	fmt.Println("\nSCFFactory")
	var prof Profile
	voter1 := []Alternative{1, 2, 4}
	voter2 := []Alternative{2, 3, 1}
	voter3 := []Alternative{1, 2, 3}
	voter4 := []Alternative{2, 1, 5}
	orderedAlts := []Alternative{2, 1, 4, 3, 5}
	prof = append(prof, voter1, voter2, voter3, voter4)
	got, err := SCFFactory(MajoritySCF, TieBreakFactory(orderedAlts))(prof)

	if err != nil {
		t.Errorf("error")
	}

	var want Alternative = 2

	if got != want {
		t.Errorf("On a %d alors que l'on veut %d etant donne, %v", got, want, prof)
	}

}
