package main

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
)

const (
	privKeyPath = "app.rsa"     // openssl genrsa -out app.rsa keysize
	pubKeyPath  = "app.rsa.pub" // openssl rsa -in app.rsa -pubout > app.rsa.pub
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

type CustomClaims struct {
	*jwt.StandardClaims
	UserData map[string]string
}

func (cc CustomClaims) Valid() error {
	return nil
}

func init() {
	signBytes, err := ioutil.ReadFile(privKeyPath)
	fatal(err)
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)
	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	fatal(err)
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)
}

func makeToken(userData map[string]string) (string, error) {
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	token.Claims = &CustomClaims{
		&jwt.StandardClaims{},
		userData,
	}
	return token.SignedString(signKey)
}

func readToken(tokenString string) map[string]string {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	fatal(err)
	claims := token.Claims.(*CustomClaims)
	return claims.UserData
}
