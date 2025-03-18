package comsoc

import (
	"fmt"
	"testing"
)

func TestVote(t *testing.T) {
	fmt.Println("\ntools")

	alternatives := []Alternative{2, 6, 5, 71, 156, 42, 17, 26}

	counts := make(Count)

	counts[2] = 0
	counts[156] = 4
	counts[5] = 15
	counts[17] = 15
	counts[26] = 4

	if rank(26, alternatives) != 7 {
		t.Errorf("error")
	}
	//fmt.Println(maxCount(counts))
}
