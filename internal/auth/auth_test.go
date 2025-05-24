package auth

import (
	"testing"
)

func TestComparePassword(t *testing.T) {
	toHash := "mysecretPass"
	hash, _ := HashPassword(toHash)
	if err := CheckPasswordHash(hash, toHash); err != nil {
		t.Errorf("CheckPasswordHash(): failed \n %e", err)
	}
}
