package link

import (
	"ShareInfo/internal/domain"
	"ShareInfo/internal/repository/database"
	"crypto/rand"
	"math/big"
)

func cacheMetaLink(link domain.Link) error {
	// Кэшируем только что созданную ссылку
	return nil
}

func generateLinkId(repo *database.LinkRepository) (int64, error) {
	newInt := big.NewInt(9999999999)
	n, err := rand.Int(rand.Reader, newInt)
	if err != nil {
		return 0, err
	}
	return n.Int64(), nil
}
