package comsoc

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMajoritySWF(t *testing.T) {
	fmt.Println("\nMajoritySWF")
	var prof Profile
	voter1 := []Alternative{1, 2, 4}
	voter2 := []Alternative{2, 3, 1}
	voter3 := []Alternative{1, 2, 3}
	prof = append(prof, voter1, voter2, voter3)
	got, err := MajoritySWF(prof)

	if err != nil {
		t.Errorf("error")
	}

	want := Count{
		1: 2,
		2: 1,
		3: 0,
	}

	for alt, val := range got {
		if val != want[alt] {
			t.Errorf("On a %d alors que l'on veut %d etant donne, %v", got, want, prof)
		}
	}
}

func TestMajoritySCF(t *testing.T) {
	fmt.Println("\nMajoritySCF")
	var prof Profile
	voter1 := []Alternative{1, 2, 3}
	voter2 := []Alternative{2, 3, 1}
	voter3 := []Alternative{1, 2, 3}
	prof = append(prof, voter1, voter2, voter3)
	got, _ := MajoritySCF(prof)
	want := []Alternative{1}
	if !(reflect.DeepEqual(got, want)) {
		t.Errorf("On a %d  alors que l'on veut %d etant donne, %v", got, want, prof)
	}
}
