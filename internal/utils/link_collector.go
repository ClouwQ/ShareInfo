package utils

import (
	"ShareInfo/internal/repository/database"
	"crypto/rand"
	"math/big"
)

func GenerateLinkId(repo *database.LinkRepository) (int64, error) {
	newInt := big.NewInt(9999999999)
	n, err := rand.Int(rand.Reader, newInt)
	if err != nil {
		return 0, err
	}
	return n.Int64(), nil
}

func StartLinkCollection(stop <-chan interface{}) {

}
