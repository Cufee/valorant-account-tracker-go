package database

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/Cufee/valorant-account-tracker-go/internal/projectpath"
	"github.com/stretchr/testify/assert"
	bolt "go.etcd.io/bbolt"
)

func makeTestClient() (*client, error) {
	var err error
	var client client
	client.db, err = bolt.Open(filepath.Join(projectpath.Root, "tmp", "client_test.db"), 0600, nil)
	return &client, err
}

func TestClientSetGet(t *testing.T) {
	client, err := makeTestClient()
	assert.Nil(t, err)

	key := "key_1"
	table := "test_table_1"
	value := []byte("Test Value")

	err = client.Set(table, key, value)
	assert.Nil(t, err)

	v, err := client.Get(table, key)
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, value, v)
}

func TestClientGetEncoded(t *testing.T) {
	client, err := makeTestClient()
	assert.Nil(t, err)

	key := "key_1"
	table := "test_table_1"
	value := map[string]string{"value_key": "value_value"}
	data, _ := json.Marshal(value)

	err = client.Set(table, key, data)
	assert.Nil(t, err)

	var decoded map[string]string
	err = client.GetEncoded(table, key, &decoded, json.Unmarshal)
	assert.Nil(t, err)
	assert.NotNil(t, decoded)
	assert.Equal(t, value["value_key"], decoded["value_key"])
}

func TestClientList(t *testing.T) {
	client, err := makeTestClient()
	assert.Nil(t, err)

	table := "test_table_2"
	values := map[string][]byte{}
	values["key_1"] = []byte("value_1")
	values["key_2"] = []byte("value_2")
	values["key_3"] = []byte("value_3")

	for k, v := range values {
		err = client.Set(table, k, v)
		assert.Nil(t, err)
	}

	list, err := client.List(table)
	assert.Nil(t, err)
	assert.Equal(t, len(values), len(list))

	for _, value := range list {
		var found bool
		for _, initialValue := range values {
			if string(value) == string(initialValue) {
				found = true
				break
			}
		}
		assert.True(t, found)
	}
}
