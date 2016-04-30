package localstorage

import "errors"

type mockStorage struct {
	m map[string]string
}

func (s *mockStorage) Get(key []byte) ([]byte, error) {
	val, ok := s.m[string(key)]
	if !ok {
		return nil, errors.New("key not exist")
	}
	return []byte(val), nil
}

func (s *mockStorage) Put(key, val []byte) error {
	s.m[string(key)] = string(val)
	return nil
}
