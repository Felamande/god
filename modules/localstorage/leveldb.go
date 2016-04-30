package localstorage

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type lvdbStorage struct {
}

func (s *lvdbStorage) Put(key []byte, value []byte) error {
	db, err := leveldb.OpenFile("localstorage.lvdb", nil)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Put(key, value, nil)
}
func (s *lvdbStorage) Get(key []byte) (value []byte, err error) {
	db, err := leveldb.OpenFile("localstorage.lvdb", nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	return db.Get(key, nil)
}
