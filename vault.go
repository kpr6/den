package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"

	"github.com/kpr6/den/encrypter"
)

func config(encodingKey, filePath string) *vault {
	return &vault{
		encodingKey: encodingKey,
		filePath:    filePath,
	}
}

type vault struct {
	encodingKey string
	filePath    string
	mutex       sync.Mutex
	keyValues   map[string]string
}

func (v *vault) load() error {
	f, err := os.OpenFile(v.filePath, os.O_RDWR|os.O_CREATE, 0640)
	defer f.Close()
	fileinfo, err := f.Stat()
	if err != nil {
		return err
	}
	// return empty vault struct if file is empty, day 0 edgecase
	if fileinfo.Size() == 0 {
		v.keyValues = make(map[string]string)
		return nil
	}

	r, err := encrypter.DecryptReader(v.encodingKey, f)
	if err != nil {
		return err
	}
	return v.readKeyValues(r)

}

func (v *vault) readKeyValues(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(&v.keyValues)
}

func (v *vault) get(key string) (string, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	err := v.load()
	if err != nil {
		return "", err
	}
	value, ok := v.keyValues[key]
	if !ok {
		return "", errors.New("secret: no value for that key")
	}
	return value, nil
}

func (v *vault) save() error {
	f, err := os.OpenFile(v.filePath, os.O_RDWR|os.O_CREATE, 0640)
	if err != nil {
		return err
	}
	defer f.Close()
	w, err := encrypter.EncryptWriter(v.encodingKey, f)
	if err != nil {
		return err
	}
	return v.writeKeyValues(w)
}

func (v *vault) writeKeyValues(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(v.keyValues)
}

func (v *vault) set(key, value string) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	err := v.load()
	if err != nil {
		return err
	}
	v.keyValues[key] = value
	return v.save()
}

func (v *vault) list() error {
	err := v.load()
	if err != nil {
		return err
	}
	enc := json.NewEncoder(os.Stdout)
	enc.Encode(v.keyValues)
	return nil
}

func (v *vault) del(key string) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	err := v.load()
	if err != nil {
		return err
	}
	delete(v.keyValues, key)
	return v.save()
}
