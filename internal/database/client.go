package database

import (
	"sync"

	bolt "go.etcd.io/bbolt"
)

var clientLock = &sync.Mutex{}
var clientInstance *client

type client struct {
	db *bolt.DB
}

func GetClient() (*client, error) {
	if clientInstance != nil {
		return clientInstance, nil
	}

	clientLock.Lock()
	defer clientLock.Unlock()

	var err error
	var client client
	client.db, err = bolt.Open(localDatabasePath, 0600, nil)
	if err != nil {
		return nil, err
	}

	clientInstance = &client
	return clientInstance, nil
}

func (c *client) Get(table, key string) ([]byte, error) {
	var value []byte
	return value, c.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(table))
		value = bucket.Get([]byte(key))
		if value == nil {
			return ErrNotFound
		}
		return nil
	})
}

func (c *client) GetEncoded(table, key string, target interface{}, decode func([]byte, interface{}) error) error {
	return c.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(table))
		value := bucket.Get([]byte(key))
		if value == nil {
			return ErrNotFound
		}

		return decode(value, target)
	})
}

func (c *client) Set(table, key string, data []byte) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(table))
		if err != nil {
			return err
		}

		return bucket.Put([]byte(key), data)
	})
}

func (c *client) List(table string) ([][]byte, error) {
	var values [][]byte
	return values, c.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(table))
		return bucket.ForEach(func(_, v []byte) error {
			values = append(values, v)
			return nil
		})
	})
}
