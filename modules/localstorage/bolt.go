package localstorage

import "github.com/boltdb/bolt"

type boltStorage struct {
}

func (bs *boltStorage) Put(key []byte, value []byte) error {
	db, err := bolt.Open("localstorage.bolt", 0777, nil)
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Batch(func(tx *bolt.Tx) error {
		bk, err := tx.CreateBucketIfNotExists([]byte("localstorage"))
		if err != nil {
			return err
		}
		return bk.Put(key, value)
	})
	return err
}

func (bs *boltStorage) Get(key []byte) (value []byte, err error) {
	// fmt.Println(runtime.Caller(0))
	// fmt.Println("hello get")
	db, err := bolt.Open("localstorage.bolt", 0777, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	// fmt.Println(runtime.Caller(0))

	err = db.Batch(func(tx *bolt.Tx) error {
		bk, err := tx.CreateBucketIfNotExists([]byte("localstorage"))
		if err != nil {
			return err
		}
		value = bk.Get(key)
		return nil
	})
	// fmt.Println(runtime.Caller(0))

	return value, err
}

func (bs *boltStorage) Clone() Storage {
	return &boltStorage{}
}
