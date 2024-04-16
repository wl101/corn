package internal

import (
	"errors"
	"log"
	"os"
	"sync"
)

type DB struct {
	File   *os.File
	Offset int64
	Indexs map[string]int64
	sync.Mutex
}

// if name not exists,create a new db with name
func NewDB(name string) (*DB, error) {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Printf("open db file failed:%v", err)
		return nil, errors.New("open db file failed")
	}
	return &DB{
		File:   file,
		Indexs: make(map[string]int64),
	}, nil
}

func (db *DB) Write(b []byte) (int, error) {
	db.Lock()
	defer db.Unlock()
	return db.File.Write(b)
}

func (db *DB) Read(b []byte, off int64) (int, error) {
	return db.File.ReadAt(b, off)
}

func (db *DB) Set(key, val string) error {
	meta := meta{
		keySize: uint32(len(key)),
		valSize: uint32(len(val)),
		key:     []byte(key),
		val:     []byte(val),
	}
	b := meta.Encode()
	n, err := db.Write(b)
	if err != nil {
		return err
	}
	db.Indexs[key] = db.Offset
	db.Offset = db.Offset + int64(n)
	return nil
}

func (db *DB) Get(key string) (string, error) {
	off, ok := db.Indexs[key]
	if !ok {
		return "", errors.New("key is not exists")
	}
	header := make([]byte, metaHeaderSize)
	_, err := db.Read(header, off)
	if err != nil {
		return "'", err
	}
	meta := meta{}
	meta.Decode(header)
	buf := make([]byte, meta.keySize+meta.valSize)
	_, err = db.Read(buf, off+metaHeaderSize)
	if err != nil {
		return "", err
	}
	meta.Decode(buf)
	return string(meta.val), nil
}
