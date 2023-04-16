package hasher

type Hasher interface {
	GetHash(s string) ([]byte, error)
	Check(hashedStr string, checkStr string) bool
}
