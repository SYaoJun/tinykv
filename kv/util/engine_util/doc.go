package engine_util

/*
An engine is a low-level system for storing key/value pairs locally (without distribution or any transaction support,
etc.). This package contains code for interacting with such engines.

CF means 'column family'. A good description of column families is given in https://github.com/facebook/rocksdb/wiki/Column-Families
(specifically for RocksDB, but the general concepts are universal). In short, a column family is a key namespace.
Multiple column families are usually implemented as almost separate databases. Importantly each column family can be
configured separately. Writes can be made atomic across column families, which cannot be done for separate databases.

CF 代表 '列族'。关于列族的一个好的描述可以在 [RocksDB 的 GitHub 页面](https://github.com/facebook/rocksdb/wiki/Column-Families) 上找到（特别是针对 RocksDB，但其一般概念是通用的）。简而言之，列族是一个键命名空间。多个列族通常实现为几乎独立的数据库。重要的是，每个列族可以单独配置。写操作可以跨列族原子化完成，而这在独立数据库之间是无法做到的。
engine_util includes the following packages:

* engines: a data structure for keeping engines required by unistore.
* write_batch: code to batch writes into a single, atomic 'transaction'.
* cf_iterator: code to iterate over a whole column family in badger.
*/
