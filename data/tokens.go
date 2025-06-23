package data

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func jwtTokenCreator(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": email,
			"exp":   time.Now().Add(time.Minute * 15).Unix(),
		})
	tokenval, err := token.SignedString(jwtkey)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return tokenval, nil
}
func verifyToken(a string) (string, time.Duration, error) {
	claim := jwt.MapClaims{}
	token, er := jwt.ParseWithClaims(a, claim, func(tok *jwt.Token) (interface{}, error) { return jwtkey, nil })
	if er != nil {
		return "", 0, errors.New("unauthorized")
	}

	if !token.Valid {
		return "", 0, errors.New("unauthorized")
	}
	exptim, _ := claim.GetExpirationTime()
	remtim := time.Until(exptim.Time)
	email := fmt.Sprintf("%v", claim["email"])
	return email, remtim, nil

}

func authhandler(a string) (string, string, error) {
	if len(a) < 8 || a == "" {
		return "", "", errors.New("invalid")
	}
	a = a[len("Bearer "):]
	email, tim, er := verifyToken(a)
	if er != nil {
		return "", "", er
	}
	if tim < time.Minute*5 && tim != 0 {
		tok, _ := jwtTokenCreator(email)
		return email, tok, nil
	} else {
		return email, a, nil
	}
}
