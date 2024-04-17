package db

import (
	"encoding/json"
	"log"

	gorocksdb "github.com/linxGnu/grocksdb"
)

const (
	cfDefault = iota
	cfHeight
	cfAddress
	cfBlockTxs
	cfTransactions
)

var cfNames = []string{"default", "height", "addresses", "blockTxs", "transactions"}

type RocksDB struct {
	path  string
	db    *gorocksdb.DB
	wo    *gorocksdb.WriteOptions
	ro    *gorocksdb.ReadOptions
	cfh   []*gorocksdb.ColumnFamilyHandle
	cache *gorocksdb.Cache
}

func createDBOptions(cache *gorocksdb.Cache, maxOpenFiles int) *gorocksdb.Options {
	bbto := gorocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(cache)
	bbto.SetBlockSize(32 << 10) // 32kB
	bbto.SetFilterPolicy(gorocksdb.NewBloomFilter(float64(10)))


	opts := gorocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)
	opts.SetCreateIfMissingColumnFamilies(true)
	opts.SetMaxOpenFiles(maxOpenFiles)

	return opts
}

func openDb(path string, cache *gorocksdb.Cache) (*gorocksdb.DB, []*gorocksdb.ColumnFamilyHandle, error) {
	opts := createDBOptions(cache, 1<<14)

	cfOptions := []*gorocksdb.Options{opts, opts, opts, opts, opts}

	db, cfh, err := gorocksdb.OpenDbColumnFamilies(opts, path, cfNames, cfOptions)

	if err != nil {
		log.Fatal("Failed to establish database connection")
		return nil, nil, err
	}

	return db, cfh, nil
}

func NewConn(path string, dbCacheSize int) (*RocksDB, error) {

	cache := gorocksdb.NewLRUCache(uint64(dbCacheSize))
	db, cfh, err := openDb(path, cache)

	if err != nil {
		return nil, err
	}

	return &RocksDB{
		path:  path,
		db:    db,
		wo:    gorocksdb.NewDefaultWriteOptions(),
		ro:    gorocksdb.NewDefaultReadOptions(),
		cfh:   cfh,
		cache: cache,
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
