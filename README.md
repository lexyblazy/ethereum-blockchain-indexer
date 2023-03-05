### A (barebones) Ethereum blockchain indexer

This is still a work in progress. I'm working on this project to learn and understand the Go programming language and blockchain technology in general.

I'm also using [Trezor's Blockbook](https://github.com/trezor/blockbook) codebase as reference point.

[Ganache](https://trufflesuite.com/ganache/) can be used to run a node locally

To run this project.

```sh
 chmod +x ./start.sh    # grant permissions to the start script
./start.sh
```


### RocksDB Setup

It's very important to use the correct versions. There seems to be a lot of breaking changes and compatibility issues between the versions. So use the below versions

- [RocksDB ](https://github.com/facebook/rocksdb) - `v7.9.2`
- [GoRocksDB](https://github.com/linxGnu/grocksdb)  - `v1.7.15`

When properly setup a `data` directory is created in the root folder when the server starts up.



#### Learning Resources

[Go Ethereum Book](https://goethereumbook.org/)

[Go by example](https://gobyexample.com/)

Other Similar Indexer projects worth studying:

- [Tezos Indexer](https://github.com/blockwatch-cc/tzindex) written in Go
- [BlockStream Electrs And Esplora ](https://github.com/Blockstream/electrs) written in Rust
