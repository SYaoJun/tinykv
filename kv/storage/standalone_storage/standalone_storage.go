package standalone_storage

import (
	"path"

	"github.com/Connor1996/badger"
	"github.com/pingcap-incubator/tinykv/kv/config"
	"github.com/pingcap-incubator/tinykv/kv/storage"
	"github.com/pingcap-incubator/tinykv/kv/util/engine_util"
	"github.com/pingcap-incubator/tinykv/proto/pkg/kvrpcpb"
)

// StandAloneStorage is an implementation of `Storage` for a single-node TinyKV instance. It does not
// communicate with other nodes and all data is stored locally.
type StandAloneStorage struct { /*实现了Storage接口中的所有方法*/
	// Your Data Here (1).
	engine *engine_util.Engines
	config *config.Config
}

func NewStandAloneStorage(conf *config.Config) *StandAloneStorage {
	// Your Code Here (1).
	dbPath := conf.DBPath
	kvPath := path.Join(dbPath, "kv")
	raftPath := path.Join(dbPath, "raft")

	kvEngine := engine_util.CreateDB(kvPath, false)
	raftEngine := engine_util.CreateDB(raftPath, true)

	engines := engine_util.NewEngines(kvEngine, raftEngine, kvPath, raftPath)
	return &StandAloneStorage{engines, conf}
}

/*为了实现StorageReader接口，必须创建一个结构体*/
type StandAloneReader struct {
	kvTxn *badger.Txn
}

/*reader为什么要主动记录事务？因为engine_util提供的函数需要txn参数*/
func NewStandAloneReader(kvTxn *badger.Txn) *StandAloneReader {
	return &StandAloneReader{
		kvTxn: kvTxn,
	}
}

func (s *StandAloneStorage) Start() error {
	// Your Code Here (1).
	return nil
}

func (s *StandAloneStorage) Stop() error {
	// Your Code Here (1).
	return s.engine.Close()
}

func (s *StandAloneStorage) Reader(ctx *kvrpcpb.Context) (storage.StorageReader, error) {
	// Your Code Here (1).
	kvTxn := s.engine.Kv.NewTransaction(false)
	return NewStandAloneReader(kvTxn), nil
}

func (s *StandAloneStorage) Write(ctx *kvrpcpb.Context, batch []storage.Modify) error {
	// Your Code Here (1).
	// 1. 利用engine_util中的badger进行存储，按列存储
	for _, b := range batch {
		switch v := b.Data.(type) {
		/*go语言的包管理？直接以文件名和结构体名称*/
		case storage.Put:
			key := v.Key
			value := v.Value
			cf := v.Cf
			err := engine_util.PutCF(s.engine.Kv, cf, key, value)
			if err != nil {
				return err
			}
		case storage.Delete:
			key := v.Key
			cf := v.Cf
			err := engine_util.DeleteCF(s.engine.Kv, cf, key)
			if err != nil {
				return nil
			}
		}
	}
	return nil
}

/*
*

	实现 StorageReader 接口
*/
func (s *StandAloneReader) GetCF(cf string, key []byte) ([]byte, error) {
	value, err := engine_util.GetCFFromTxn(s.kvTxn, cf, key)
	// key 不存在
	if err == badger.ErrKeyNotFound {
		return nil, nil // 测试要求 err 为 nil，而不是 KeyNotFound，否则没法过
	}
	return value, err
}

func (s *StandAloneReader) IterCF(cf string) engine_util.DBIterator {
	return engine_util.NewCFIterator(cf, s.kvTxn)
}

func (s *StandAloneReader) Close() {
	s.kvTxn.Discard()
}
