package tokenz

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"regexp"
	"os"
	"io/ioutil"
	"strings"
	"io"
)

// Print a json object in accordance with the prophecy
func PrintJSON(j interface{}) error {
	out, err := json.MarshalIndent(j, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func GetToken(alg string, key string, tokData []byte) (*jwt.Token, error) {
	// get the token
	//tokData, err := loadData(*flagVerify)
	//if err != nil {
	//	return fmt.Errorf("Couldn't read token: %v", err)
	//}

	// trim possible whitespace from token
	tokData = regexp.MustCompile(`\s*$`).ReplaceAll(tokData, []byte{})

	//if *flagDebug
	//{
	//	fmt.Fprintf(os.Stderr, "Token len: %v bytes\n", len(tokData))
	//}

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		data, err := loadData(key)
		if err != nil {
			return nil, err
		}
		if isEs(alg) {
			return jwt.ParseECPublicKeyFromPEM(data)
		}
		return data, nil
	}
	// Parse the token.  Load the key from command line option
	return jwt.Parse(string(tokData), keyFunc)

	//return nil
}


// Create, sign, and output a token.  This is a great, simple example of
// how to use this library to create and sign a token.
func SignToken(alg string,flagKey string, tokData []byte) (string,error) {

	// parse the JSON of the claims
	var claims map[string]interface{}
	if err := json.Unmarshal(tokData, &claims); err != nil {
		return "", fmt.Errorf("Couldn't parse claims JSON: %v", err)
	}

	// get the key
	var key interface{}
	key, err := loadData(flagKey)
	if err != nil {
		return "", fmt.Errorf("Couldn't read key: %v", err)
	}

	// get the signing alg
	signMethod := jwt.GetSigningMethod(alg)
	if signMethod == nil {
		return "", fmt.Errorf("Couldn't find signing method: %v", alg)
	}

	// create a new token
	token := jwt.New(signMethod)
	token.Claims = claims

	if isEs(alg) {
		if k, ok := key.([]byte); !ok {
			return "", fmt.Errorf("Couldn't convert key data to key")
		} else {
			key, err = jwt.ParseECPrivateKeyFromPEM(k)
			if err != nil {
				return "", err
			}
		}
	}

	return token.SignedString(key)
}

// Helper func:  Read input from specified file or stdin
func loadData(p string) ([]byte, error) {
	if p == "" {
		return nil, fmt.Errorf("No path specified")
	}

	var rdr io.Reader
	if p == "-" {
		rdr = os.Stdin
	} else {
		if f, err := os.Open(p); err == nil {
			rdr = f
			defer f.Close()
		} else {
			return nil, err
		}
	}
	return ioutil.ReadAll(rdr)
}

//func writeData(path, content string) {
//	ioutil.WriteFile(path, []byte(content), 0666 )
//}

func isEs(alg string) bool {
	return strings.HasPrefix(alg, "ES")
}
