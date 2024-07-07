package database

import(
	"fmt"
	"errors"

	badger "github.com/dgraph-io/badger/v4"
)

type Alert struct {
    Name     string  `json:"name"`
    Instance  string  `json:"instance"`
}

func GetAll() ([]Alert, error) {
	var alerts []Alert
	err := db.View(func(txn *badger.Txn) error {
	  opts := badger.DefaultIteratorOptions
	  it := txn.NewIterator(opts)
	  defer it.Close()
	  for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		k := item.Key()
		err := item.Value(func(v []byte) error {
			alerts = append(alerts, Alert{Name: string(k), Instance: string(v)})
		  // fmt.Printf("key=%s, value=%s\n", k, v)
		  return nil
		})
		if err != nil {
		  return err
		}
	  }
	  return nil
	})
	if err != nil {
		return []Alert{}, err
	}
	return alerts, nil
}

var db, err = badger.Open(badger.DefaultOptions("/tmp/badger"))


func Close() error {
 return db.Close()
}

// nolint:wrapcheck
func exists(key string) (bool, error) {
 var exists bool
 err := db.View(
  func(tx *badger.Txn) error {
   if val, err := tx.Get([]byte(key)); err != nil {
    return err
   } else if val != nil {
    exists = true
   }
   return nil
  })
 if errors.Is(err, badger.ErrKeyNotFound) {
  err = nil
 }
 return exists, err
}

func Get(key string) (string, error) {
 var value string
 return value, db.View(
  func(tx *badger.Txn) error {
   item, err := tx.Get([]byte(key))
   if err != nil {
    return fmt.Errorf("getting value: %w", err)
   }
   valCopy, err := item.ValueCopy(nil)
   if err != nil {
    return fmt.Errorf("copying value: %w", err)
   }
   value = string(valCopy)
   return nil
  })
}

func Set(key, value string) error {
 return db.Update(
  func(txn *badger.Txn) error {
   return txn.Set([]byte(key), []byte(value))
  })
}

func Delete(key string) error {
 return db.Update(
  func(txn *badger.Txn) error {
   return txn.Delete([]byte(key))
  })
}
