package comsoc

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBordaSWF(t *testing.T) {
	fmt.Println("\nBordaSWF")
	var prof Profile
	voter1 := []Alternative{1, 2, 3}
	voter2 := []Alternative{2, 3, 1}
	voter3 := []Alternative{1, 2, 3}
	prof = append(prof, voter1, voter2, voter3)
	got, _ := BordaSWF(prof)
	want := map[Alternative]int{
		1: 4,
		2: 4,
		3: 1,
	}

	for alt, val := range got {
		if val != want[alt] {
			t.Errorf("On a %d alors que l'on veut %d etant donne, %v", got, want, prof)
		}
	}
}

func TestBordaSCF(t *testing.T) {
	fmt.Println("\nBordaSCF")
	var prof Profile
	voter1 := []Alternative{1, 2, 3}
	voter2 := []Alternative{2, 3, 1}
	voter3 := []Alternative{1, 2, 3}
	prof = append(prof, voter1, voter2, voter3)
	got, _ := BordaSCF(prof)
	want1 := []Alternative{1, 2}
	want2 := []Alternative{2, 1}
	if !(reflect.DeepEqual(got, want1) || reflect.DeepEqual(got, want2)) {
		t.Errorf("On a %d  alors que l'on veut %d ou %d etant donne, %v", got, want1, want2, prof)
	}
}
