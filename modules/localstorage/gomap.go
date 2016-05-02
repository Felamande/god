package localstorage

import "errors"

type mapStorage struct {
	m map[string]string
}

func (s *mapStorage) Get(key []byte) ([]byte, error) {
	val, ok := s.m[string(key)]
	if !ok {
		return nil, errors.New("key not exist")
	}
	return []byte(val), nil
}

func (s *mapStorage) Put(key, val []byte) error {
	s.m[string(key)] = string(val)
	return nil
}
