package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/thinhlu123/shortener/config"
	"time"
)

var serviceToken = map[string]string{
	"test": "aaaa",
}

// Credentials Create a struct to read the username and password from the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Claims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateAuthToken Create token
func GenerateAuthToken(user, pwd string) (string, error) {
	creds := Credentials{
		Username: user,
		Password: pwd,
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString([]byte(config.Conf.Server.JwtSecretKey))
	if err != nil {
		// If there is an error in creating the JWT return an internal server error

		return "", err
	}

	return tokenString, nil
}

func RefreshToken(token string) (string, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Conf.Server.JwtSecretKey), nil
	})
	if err != nil {
		return "", err
	}
	if !tkn.Valid {
		return "", errors.New("invalid token")
	}
	// (END) The code up-till this point is the same as the first part of the `Welcome` route

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < 30*time.Second {
		return "", errors.New("expiry token")
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := newToken.SignedString([]byte(config.Conf.Server.JwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUsernameFromToken(token string) (string, error) {
	if user, ok := serviceToken[token]; ok {
		return user, nil
	}

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Conf.Server.JwtSecretKey), nil
	})
	if err != nil {
		return "", err
	}
	if !tkn.Valid {
		return "", errors.New("invalid token")
	}

	return claims.Username, nil
}
