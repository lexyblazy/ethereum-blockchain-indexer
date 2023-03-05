package db

import (
	gorocksdb "github.com/linxGnu/grocksdb"
	"log"
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

func NewConn(path string, dbCacheSize int) *RocksDB {

	opts := createDBOptions(dbCacheSize)

	db, err := gorocksdb.OpenDb(opts, path)

	if err != nil {
		log.Fatal("Failed to establish database connection")
	}

	return &RocksDB{
		path: path,
		db:   db,
		wo:   gorocksdb.NewDefaultWriteOptions(),
		ro:   gorocksdb.NewDefaultReadOptions(),
	}

}
