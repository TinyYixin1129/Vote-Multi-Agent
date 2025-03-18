package comsoc

import (
	"fmt"
	"testing"
)

func TestCondoecetWinner(t *testing.T) {
	fmt.Println("\nCondorcetWinner")
	var prof Profile
	voter1 := []Alternative{1, 2, 3}
	voter2 := []Alternative{2, 3, 1}
	voter3 := []Alternative{3, 1, 2}
	prof = append(prof, voter1, voter2, voter3)
	got, _ := CondorcetWinner(prof)
	want := []Alternative{}
	if got != nil {
		t.Errorf("On a %d alors que l'on veut %d etant donne, %v", got, want, prof)
	}
	/* for i, val := range got {
		if val != want[i] {
			t.Errorf("On a %d alors que l'on veut %d etant donne, %v", got, want, prof)
		}
	} */
}
