package server

import (
	"context"

	"github.com/pingcap-incubator/tinykv/kv/storage"
	"github.com/pingcap-incubator/tinykv/proto/pkg/kvrpcpb"
)

// The functions below are Server's Raw API. (implements TinyKvServer).
// Some helper methods can be found in sever.go in the current directory

// RawGet return the corresponding Get response based on RawGetRequest's CF and Key fields
func (server *Server) RawGet(_ context.Context, req *kvrpcpb.RawGetRequest) (*kvrpcpb.RawGetResponse, error) {
	// Your Code Here (1).
	reader, err := server.storage.Reader(req.Context)
	if err != nil {
		return &kvrpcpb.RawGetResponse{}, err
	}

	// 获取 value
	cf := req.Cf
	key := req.Key
	value, err := reader.GetCF(cf, key)
	if err != nil {
		return &kvrpcpb.RawGetResponse{}, err
	}

	// 返回响应
	notFound := false
	if value == nil {
		notFound = true
	}
	resp := &kvrpcpb.RawGetResponse{
		Value:    value,
		NotFound: notFound,
	}

	return resp, nil
}

// RawPut puts the target data into storage and returns the corresponding response
func (server *Server) RawPut(_ context.Context, req *kvrpcpb.RawPutRequest) (*kvrpcpb.RawPutResponse, error) {
	// Your Code Here (1).
	// Hint: Consider using Storage.Modify to store data to be modified
	cf := req.Cf
	key := req.Key
	value := req.Value
	err := server.storage.Write(req.Context, []storage.Modify{
		{
			Data: storage.Put{
				Cf:    cf,
				Key:   key,
				Value: value,
			},
		},
	})

	return &kvrpcpb.RawPutResponse{}, err
}

// RawDelete delete the target data from storage and returns the corresponding response
func (server *Server) RawDelete(_ context.Context, req *kvrpcpb.RawDeleteRequest) (*kvrpcpb.RawDeleteResponse, error) {
	// Your Code Here (1).
	// Hint: Consider using Storage.Modify to store data to be deleted
	cf := req.Cf
	key := req.Key
	err := server.storage.Write(req.Context, []storage.Modify{
		{
			Data: storage.Delete{
				Cf:  cf,
				Key: key,
			},
		},
	})

	return &kvrpcpb.RawDeleteResponse{}, err
}

// RawScan scan the data starting from the start key up to limit. and return the corresponding result
func (server *Server) RawScan(_ context.Context, req *kvrpcpb.RawScanRequest) (*kvrpcpb.RawScanResponse, error) {
	// Your Code Here (1).
	// Hint: Consider using reader.IterCF
	reader, err := server.storage.Reader(req.Context)
	if err != nil {
		return &kvrpcpb.RawScanResponse{}, err
	}

	cf := req.Cf
	iter := reader.IterCF(cf)
	var kvs []*kvrpcpb.KvPair
	cnt := uint32(0)
	for iter.Seek(req.StartKey); iter.Valid(); iter.Next() {
		item := iter.Item()
		/*如果value无效，是不是应该忽略呢？*/
		val, _ := item.Value()
		kvs = append(kvs, &kvrpcpb.KvPair{
			Key:   item.Key(),
			Value: val,
		})
		cnt++
		if cnt == req.Limit {
			break
		}
	}

	return &kvrpcpb.RawScanResponse{Kvs: kvs}, nil
}
