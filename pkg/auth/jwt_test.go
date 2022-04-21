package auth

import (
	"github.com/stretchr/testify/assert"
	"github.com/thinhlu123/shortener/config"
	"log"
	"testing"
)

func beforeTest() {
	if err := config.GetConfig("../../config/config_local"); err != nil {
		log.Fatalf("Loading config: %v", err)
	}
}

func TestGenerateAuthToken(t *testing.T) {
	beforeTest()

	user := map[string]string{
		"a": "ad213@!",
		"b": "avagcdq2412",
		"c": "ACNJSndjnqw",
		"d": "ANCSJ123sdn",
		"e": "AAAAA123abc!@",
		"f": "!@#$%",
		"g": "$!@ssfacsaADMJWEQ'",
		"h": "bưdjbajkdsbjkaaskdnakndkandkasndkmasndksnjkdnjkasndjkansdjkbajghvdhwqbjdjkabskjd",
	}

	for usr, pwd := range user {
		token, err := GenerateAuthToken(usr, pwd)
		if err != nil {
			t.Error(err)
		}
		t.Log(usr, token)
	}
}

func TestRefreshToken(t *testing.T) {
	var (
		usr = "test"
		pwd = "admin"
	)
	beforeTest()

	token, err := GenerateAuthToken(usr, pwd)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	newToken, err := RefreshToken(token)
	if err != nil {
		t.Error(err)
	}

	t.Log(newToken)
}

func TestGetUsernameFromToken(t *testing.T) {
	var (
		usr = "test"
		pwd = "admin"
	)
	beforeTest()

	as := assert.New(t)

	token, err := GenerateAuthToken(usr, pwd)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	username, err := GetUsernameFromToken(token)
	if err != nil {
		t.Error(err)
	}

	if as.Equal(username, usr) {
		t.Errorf("User get from token is %v when expect %v", username, usr)
	}
}
