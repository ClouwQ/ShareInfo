package utils

import (
	"ShareInfo/internal/repository/database"
	"context"
	"crypto/rand"
	"math/big"
)

func GenerateLinkId(repo *database.LinkRepository) (int64, error) {
	newInt := big.NewInt(9999999999)
	n, err := rand.Int(rand.Reader, newInt)
	if err != nil {
		return 0, err
	}

	// Проверка, что такой ссылки в базе еще нет
	link, err := repo.GetByID(context.Background(), n.Int64())
	if link == nil {
		return n.Int64(), nil
	} else if err != nil {
		return 0, err
	} else {
		return GenerateLinkId(repo)
	}
}

func StartLinkCollection(stop <-chan interface{}) {

}
