package bdb

import (
	"encoding/binary"
	"errors"
	"path/filepath"
	"task/internal/global"
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
	path := filepath.Join(global.RuntimePath, "data", "task.bdb.db")
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	rtnData := &Bdb{db: db}
	err = rtnData.InitTable()
	if err != nil {
		return nil, err
	}

	return rtnData, nil
}

func (b *Bdb) InitTable() error {
	return b.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(TaskName)
		return err
	})
}

func (b *Bdb) Close() error {
	return b.db.Close()
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
