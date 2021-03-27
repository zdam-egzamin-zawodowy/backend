package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type Credentials struct {
	Email    string
	Password string
}

type Metadata struct {
	StaySignedIn bool
	Credentials  Credentials
}

func (opts Metadata) ToMapClaims() jwt.MapClaims {
	mClaims := jwt.MapClaims{}
	mClaims["email"] = opts.Credentials.Email
	mClaims["password"] = opts.Credentials.Password
	mClaims["stay_signed_in"] = opts.StaySignedIn
	return mClaims
}

type TokenGenerator interface {
	Generate(metadata Metadata) (string, error)
	ExtractAccessTokenMetadata(token string) (*Metadata, error)
}

type tokenGenerator struct {
	accessSecret string
}

func NewTokenGenerator(accessSecret string) TokenGenerator {
	return &tokenGenerator{
		accessSecret: accessSecret,
	}
}

func (g *tokenGenerator) Generate(metadata Metadata) (string, error) {
	atClaims := metadata.ToMapClaims()
	if !metadata.StaySignedIn {
		atClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err := at.SignedString([]byte(g.accessSecret))
	if err != nil {
		return "", errors.Wrap(err, "TokenGenerator.Generate")
	}

	return accessToken, nil
}

func (g *tokenGenerator) ExtractAccessTokenMetadata(token string) (*Metadata, error) {
	return extractTokenMetadata(g.accessSecret, token)
}

func verifyToken(secret string, tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "verifyToken")
	}
	if !token.Valid {
		return nil, errors.New("Token is invalid")
	}
	return token, nil
}

func extractTokenMetadata(secret, tokenString string) (*Metadata, error) {
	token, err := verifyToken(secret, tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		staySignedIn, ok := claims["stay_signed_in"].(bool)
		if !ok {
			return nil, errors.New("Invalid token payload (staySignedIn should be a boolean)")
		}
		email, ok := claims["email"].(string)
		if !ok {
			return nil, errors.New("Invalid token payload (email should be a string)")
		}
		password, ok := claims["password"].(string)
		if !ok {
			return nil, errors.New("Invalid token payload (password should be a string)")
		}
		return &Metadata{
			StaySignedIn: staySignedIn,
			Credentials: Credentials{
				Email:    email,
				Password: password,
			},
		}, nil
	}

	return nil, errors.New("Cannot extract token metadata")
}
