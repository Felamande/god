package localstorage

import (
	"os"

	"github.com/steveyen/gkvlite"
)

type kvdbStorage struct {
}

func (s *kvdbStorage) Put(key []byte, value []byte) error {
	f, err := os.OpenFile("localstorage.db", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	store, err := gkvlite.NewStore(f)
	if err != nil {
		return err
	}
	defer store.Close()
	c := store.SetCollection("localstorage", nil)
	err = c.Set(key, value)
	store.Flush()
	f.Sync()
	return err

}
func (s *kvdbStorage) Get(key []byte) (value []byte, err error) {
	f, err := os.OpenFile("localstorage.db", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	store, err := gkvlite.NewStore(f)
	if err != nil {
		return nil, err
	}
	defer store.Close()
	c := store.SetCollection("localstorage", nil)
	value, err = c.Get(key)
	store.Flush()
	f.Sync()
	return
}
