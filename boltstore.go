package quoteapi

import (
	"encoding/json"
	"errors"
	"strings"

	bolt "go.etcd.io/bbolt"
)

func getKey(primaryKeys []string) string {
	return strings.Join(primaryKeys, " :: ")
}

func (store *boltStore) LoadRecord(tableName string, record interface{}, primaryKeys ...string) error {
	if store.database == nil {
		return errors.New("Database is nil")
	}

	key := getKey(primaryKeys)

	if err := store.database.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(tableName))
		if bucket == nil {
			return errors.New("Missing bucket")
		}
		bytes := bucket.Get([]byte(key))

		if bytes == nil {
			return errors.New("Empty record")
		}

		if err := json.Unmarshal(bytes, record); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (store *boltStore) SaveRecord(tableName string, record interface{}, primaryKeys ...string) error {
	if store.database == nil {
		return errors.New("Database is nil")
	}

	bytes, err := json.Marshal(record)
	if err != nil {
		return err
	}

	key := getKey(primaryKeys)

	if err := store.database.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(tableName))

		if err != nil {
			return err
		}
		if err := b.Put([]byte(key), bytes); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (store *boltStore) DeleteRecord(tableName string, primaryKeys ...string) error {
	if store.database == nil {
		return errors.New("Database is nil")
	}

	key := getKey(primaryKeys)

	if err := store.database.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tableName))

		if b == nil {
			return errors.New("Bucket missing")
		}
		if err := b.Delete([]byte(key)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (store *boltStore) Close() error {
	if store.database == nil {
		return errors.New("Database is nil")
	}
	return store.database.Close()
}

type boltStore struct {
	database *bolt.DB
}

func GetBoltStore(path string) (Datastore, error) {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}
	return &boltStore{database: db}, nil
}
