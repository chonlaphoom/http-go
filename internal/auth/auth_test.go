package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

// TODO: refactor
func TestComparePassword(t *testing.T) {
	toHash := "mysecretPass"
	hash, _ := HashPassword(toHash)
	if err := CheckPasswordHash(hash, toHash); err != nil {
		t.Errorf("CheckPasswordHash(): failed \n %e", err)
		return
	}
}

func TestValidateJWT(t *testing.T) {
	var secretToken = "abc"
	_uuid, _ := uuid.NewRandom()
	jwt, err := MakeJWT(_uuid, secretToken, time.Hour)
	if err != nil {
		t.Errorf("make jwt token error \n %e", err)
		return
	}

	ok, err := ValidateJWT(jwt, secretToken)
	if err != nil {
		t.Errorf("Token is not valid\n %e", err)
		return
	}
	t.Logf("Token is valid: %v\n", ok)
}

func TestExpiredJWT(t *testing.T) {
	var secretToken = "abc"
	_uuid, _ := uuid.NewRandom()
	jwt, err := MakeJWT(_uuid, secretToken, time.Second)
	if err != nil {
		t.Errorf("make jwt token error \n %e", err)
		return
	}
	time.Sleep(time.Second * 3)

	_, err_va := ValidateJWT(jwt, secretToken)
	if err_va != nil {
		t.Logf("Token is expired and it's expected\n %e", err_va)
		return
	}
	t.Error("expect token to expired")
}

func TestGetBearerToken(t *testing.T) {
	header := map[string][]string{"Authorization": {"Bearer 123"}}
	token, err := GetBearerToken(header)
	if err != nil || token != "123" {
		t.Error("expect getting bearer 123")
	}
}

func TestMakeRefreshTokens(t *testing.T) {
	_, err := MakeRefreshToken()

	if err != nil {
		t.Error("expect making refresh token successfully")
		return
	}
}
