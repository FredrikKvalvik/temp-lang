package object

import "hash/fnv"

type HashKey struct {
	Type ObjectType
	Hash float64
}

type Hashable interface {
	HashKey() HashKey
}

func (s *StringObj) HashKey() HashKey {
	hash := fnv.New64a()
	hash.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Hash: float64(hash.Sum64())}
}

func (s *NumberObj) HashKey() HashKey {
	return HashKey{Type: s.Type(), Hash: s.Value}
}

func (s *BooleanObj) HashKey() HashKey {
	var hash uint64
	if s.Value {
		hash = 1
	} else {
		hash = 0
	}
	return HashKey{Type: s.Type(), Hash: float64(hash)}
}
