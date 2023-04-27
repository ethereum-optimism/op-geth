package vm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type KVDB interface {
	Get(key []byte) ([]byte, error)
	Create(key, value []byte) error
	Update(key, value []byte) error
	Delete(key []byte) error
	Close() error
}

type TopiaDB struct {
	Url string
}

func NewTopiaDB(url string) *TopiaDB {
	return &TopiaDB{Url: url}
}

func (db *TopiaDB) Get(key []byte) ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf("%v/get/%s", db.Url, key))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get request failed with status: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]string
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return []byte(result["value"]), nil
}

func (db *TopiaDB) Create(key, value []byte) error {
	return db.put(key, value, "create")
}

func (db *TopiaDB) Update(key, value []byte) error {
	return db.put(key, value, "update")
}

func (db *TopiaDB) put(key, value []byte, action string) error {
	data := fmt.Sprintf(`{"key":"%s", "value":"%s"}`, string(key), string(value))
	resp, err := http.Post(fmt.Sprintf("%v/%v", db.Url, action), "application/json",
		strings.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%v request failed with status: %d", action, resp.StatusCode)
	}
	return nil
}

func (db *TopiaDB) Delete(key []byte) error {
	req, err := http.NewRequest(http.MethodDelete,
		fmt.Sprintf("%s/delete/%s", db.Url, key), nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete request failed with status: %d", resp.StatusCode)
	}
	return nil
}

func (db *TopiaDB) Lock() error {
	resp, err := http.Get(fmt.Sprintf("%v/lock_tx", db.Url))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("lock request failed with status: %d", resp.StatusCode)
	}

	return nil
}

func (db *TopiaDB) Unlock() error {
	resp, err := http.Get(fmt.Sprintf("%v/unlock_tx", db.Url))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unlock request failed with status: %d", resp.StatusCode)
	}

	return nil
}

func (db *TopiaDB) Close() error {
	return nil
}
