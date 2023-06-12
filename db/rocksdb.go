package db

import (
	"encoding/json"
	"log"

	gorocksdb "github.com/linxGnu/grocksdb"
)

type RocksDB struct {
	path string
	db   *gorocksdb.DB
	wo   *gorocksdb.WriteOptions
	ro   *gorocksdb.ReadOptions
}

func createDBOptions(dbCacheSize int) *gorocksdb.Options {
	bbto := gorocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(gorocksdb.NewLRUCache(uint64(dbCacheSize)))
	bbto.SetBlockSize(32 << 10) // 32kB

	opts := gorocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)

	return opts
}

func openDb(path string, dbCacheSize int) (*gorocksdb.DB, error) {
	opts := createDBOptions(dbCacheSize)

	db, err := gorocksdb.OpenDb(opts, path)

	if err != nil {
		log.Fatal("Failed to establish database connection")
		return nil, err
	}

	return db, nil
}

func NewConn(path string, dbCacheSize int) (*RocksDB, error) {

	db, err := openDb(path, dbCacheSize)

	if err != nil {
		return nil, err
	}

	return &RocksDB{
		path: path,
		db:   db,
		wo:   gorocksdb.NewDefaultWriteOptions(),
		ro:   gorocksdb.NewDefaultReadOptions(),
	}, nil

}

func (d *RocksDB) SaveRecord(key string, record struct{}) {

	serializedRecord, err := json.Marshal(record)

	if err != nil {
		log.Println("Failed to serialize record to byte array")
		return
	}

	err2 := d.db.Put(d.wo, []byte(key), serializedRecord)

	if err2 != nil {
		log.Println("failed to put key", key)
	}
}
