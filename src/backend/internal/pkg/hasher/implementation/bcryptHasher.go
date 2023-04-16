package implementation

import (
	"backend/internal/pkg/hasher"
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
}

func NewBcryptHasher() hasher.Hasher {
	return &BcryptHasher{}
}

func (b *BcryptHasher) GetHash(s string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
}

func (b *BcryptHasher) Check(hashedStr string, checkStr string) bool {
	res := bcrypt.CompareHashAndPassword([]byte(hashedStr), []byte(checkStr))
	return res == nil
}
