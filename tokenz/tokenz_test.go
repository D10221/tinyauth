package tokenz_test

import (
	"testing"
	"fmt"
	"github.com/D10221/tinyauth/tokenz"
	"encoding/json"
)

type User struct {
	Email string `json:"email"`
}

type TokenData struct {
	Id string `json:"id"`
	User *User `json:"user"`
}

func (data *TokenData) ToJson() ([]byte, error) {
	return json.Marshal(data)
}

func Test_Tokenz(t *testing.T){

	key:= "../testdata/sample_key"
	tokData, err:= (
	& TokenData {
		User: &User{
			Email: "me@no.mail",
		},
		Id : "ABC"}).ToJson()

	if err != nil { t.Error(err) ; return }


	alg:= "RS256"
	out, err := tokenz.SignToken(alg, key, tokData)
	if err != nil { t.Error(err) ; return }

	// to Token
	token, err := tokenz.GetToken(alg,key+".pub", []byte(out))
	if  token == nil {
		t.Error("Error: No Token")
		return
	}

	// Print an error if we can't parse for some reason
	if err != nil {
		t.Error(fmt.Errorf("Couldn't parse token: %v", err))
		return
	}

	// Is token invalid?
	if !token.Valid {
		t.Error(fmt.Errorf("Token is invalid"))
		return
	}

	jsonBytes, err := json.Marshal(token.Claims)
	if err!= nil {
		t.Error(err)
		return
	}
	data:= &TokenData{}
	e:= json.Unmarshal(jsonBytes, data)
	if e!=nil {
		t.Error(e)
		return;
	}
	if data.Id != "ABC" {
		t.Error("?")
	}
	if data.User.Email != "me@no.mail" {
		t.Error("?")
	}

	// Print the token details
	if err := tokenz.PrintJSON(token.Claims); err != nil {
		t.Error(fmt.Errorf("Failed to output claims: %v", err))
	}

}
