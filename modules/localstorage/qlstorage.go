package localstorage

import (
	"errors"
	"fmt"

	"github.com/cznic/ql"
)

type qlStorage struct {
}

func (s *qlStorage) Put(key []byte, value []byte) error {
	db, err := ql.OpenFile("localstorage.qldb", &ql.Options{CanCreate: true})
	if err != nil {
		return err
	}
	defer db.Close()

	l, _ := ql.Compile(`begin transaction;create table localstorage(key string,val string);commit`)
	db.Execute(ql.NewRWCtx(), l)

	l, err = ql.Compile(fmt.Sprintf(`begin transaction;insert into localstorage values("%s","%s");commit`, string(key), string(value)))
	if err != nil {
		return err
	}
	_, _, err = db.Execute(ql.NewRWCtx(), l)

	return err

}
func (s *qlStorage) Get(key []byte) (value []byte, err error) {
	db, err := ql.OpenFile("localstorage.qldb", &ql.Options{CanCreate: true})
	if err != nil {
		return nil, err
	}
	defer db.Close()

	l, _ := ql.Compile(`begin transaction;create table localstorage(key string,val string);commit`)
	db.Execute(ql.NewRWCtx(), l)
	l, err = ql.Compile(fmt.Sprintf(`begin transaction;select val from localstorage where key == "%s";commit`, string(key)))
	if err != nil {
		// fmt.Print("45", err)
		return nil, err
	}

	rs, _, err := db.Execute(ql.NewRWCtx(), l)
	if err != nil {
		// fmt.Print("50", err)

		return nil, err
	}
	for _, r := range rs {
		r.Do(false, func(data []interface{}) (bool, error) {
			if len(data) == 0 {
				return false, errors.New("cannot find key " + string(key))
			}
			value = []byte(data[0].(string))
			return true, nil
		})
	}
	return

}
