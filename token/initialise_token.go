package token

import (
	"os"
)

func InitializePasetoToken() (Maker, error) {
	TOKEN_SYMETRIC_KEY := os.Getenv("TOKEN_SYMETRIC_KEY")
	maker, err := NewPasetoMaker(TOKEN_SYMETRIC_KEY)
	return maker, err
}
