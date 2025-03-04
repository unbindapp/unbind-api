package utils

import (
	"crypto/rand"
	"math/big"
)

func GenerateRandomSimpleID(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	charsetLength := big.NewInt(int64(len(charset)))

	result := make([]byte, length)

	// Fill the slice with random characters from the charset
	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}

		result[i] = charset[randomIndex.Int64()]
	}

	return string(result), nil
}
