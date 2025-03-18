package comsoc

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCopelandSWF(t *testing.T) {
	fmt.Println("\nCopelandSWF")
	var prof Profile
	voter1 := []Alternative{1, 2, 3}
	voter2 := []Alternative{2, 3, 1}
	voter3 := []Alternative{3, 1, 2}
	prof = append(prof, voter1, voter2, voter3)
	got, _ := CopelandSWF(prof)
	want := map[Alternative]int{
		1: 0,
		2: 0,
		3: 0,
	}

	for alt, val := range got {
		if val != want[alt] {
			t.Errorf("On a %d alors que l'on veut %d etant donne, %v", got, want, prof)
		}
	}
}

func TestCopelandSCF(t *testing.T) {
	fmt.Println("\nCopelandSCF")
	var prof Profile
	voter1 := []Alternative{1, 2, 3}
	voter2 := []Alternative{1, 2, 3}
	voter3 := []Alternative{1, 2, 3}
	prof = append(prof, voter1, voter2, voter3)
	got, _ := CopelandSCF(prof)
	want1 := []Alternative{1}
	want2 := []Alternative{1}
	if !(reflect.DeepEqual(got, want1) || reflect.DeepEqual(got, want2)) {
		t.Errorf("On a %d  alors que l'on veut %d ou %d etant donne, %v", got, want1, want2, prof)
	}
}
