package secrets

type Secret struct {
	Value *string
	ContentType *string
}

func (s Secret) String() string {
	return *s.Value
}
