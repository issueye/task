package bdb

import (
	"encoding/binary"
	"errors"
	"time"

	bolt "go.etcd.io/bbolt"
)

var boltDB *Bdb

func InitBdb() error {
	var err error
	boltDB, err = New()
	return err
}

func GetBdb() *Bdb {
	return boltDB
}

var (
	ErrBucketNotFound = errors.New("bucket not found")
	ErrRecordNotFound = errors.New("record not found")
)

type Bdb struct {
	db *bolt.DB
}

func New() (*Bdb, error) {
	db, err := bolt.Open("task.bdb.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	return &Bdb{db: db}, nil
}

func (b *Bdb) Close() error {
	return b.db.Close()
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
