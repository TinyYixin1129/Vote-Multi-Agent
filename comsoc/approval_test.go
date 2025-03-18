package comsoc

import (
	"fmt"
	"reflect"
	"testing"
)

func TestApprovalSWF(t *testing.T) {
	fmt.Println("\nApprovalSWF")
	var prof Profile
	tre := []int{2, 2, 2}
	voter1 := []Alternative{1, 2, 3}
	voter2 := []Alternative{2, 1}
	voter3 := []Alternative{1}
	prof = append(prof, voter1, voter2, voter3)
	got, _ := ApprovalSWF(tre)(prof)
	want := map[Alternative]int{
		1: 3,
		2: 2,
		3: 0,
	}

	for alt, val := range got {
		if val != want[alt] {
			t.Errorf("On a %d alors que l'on veut %d etant donne, %v", got, want, prof)
		}
	}
}

func TestApprovalSCF(t *testing.T) {
	fmt.Println("\nApprovalSCF")
	var prof Profile
	tre := []int{4, 4, 4}
	voter1 := []Alternative{1, 2, 3}
	voter2 := []Alternative{2, 1}
	voter3 := []Alternative{1}
	prof = append(prof, voter1, voter2, voter3)
	got, _ := ApprovalSCF(tre)(prof)
	want := []Alternative{1}
	if !(reflect.DeepEqual(got, want)) {
		t.Errorf("On a %d  alors que l'on veut %d etant donne, %v", got, want, prof)
	}
}
