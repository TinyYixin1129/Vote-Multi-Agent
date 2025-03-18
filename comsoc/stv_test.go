package comsoc

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSTV_SWF(t *testing.T) {
	fmt.Println("\nSTV_SWF")
	var prof Profile
	voter1 := []Alternative{1, 2, 3}
	voter2 := []Alternative{1, 2, 3}
	//voter3 := []Alternative{3, 1, 2}
	prof = append(prof, voter1, voter2)
	got, _ := STV_SWF(prof)
	want := map[Alternative]int{
		1: 2,
		2: 0,
		3: 0,
	}

	for alt, val := range got {
		if val != want[alt] {
			t.Errorf("On a %d alors que l'on veut %d etant donne, %v", got, want, prof)
		}
	}
}

func TestSTC_SCF(t *testing.T) {
	fmt.Println("\nSTV_SCF")
	var prof Profile
	voter1 := []Alternative{1, 2, 3}
	voter2 := []Alternative{1, 2, 3}
	voter3 := []Alternative{1, 2, 3}
	prof = append(prof, voter1, voter2, voter3)
	got, _ := STV_SCF(prof)
	want1 := []Alternative{1}
	want2 := []Alternative{1}
	if !(reflect.DeepEqual(got, want1) || reflect.DeepEqual(got, want2)) {
		t.Errorf("On a %d  alors que l'on veut %d ou %d etant donne, %v", got, want1, want2, prof)
	}
}
