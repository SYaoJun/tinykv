package main

import (
	"fmt"
	"log"

	"github.com/Connor1996/badger"
)

/*
1. 打开数据库
2. 插入数据
3. 点查
4. 遍历
5. 前缀遍历
*/
func main() {
	opts := badger.DefaultOptions
	opts.Dir = "./qq"
	opts.ValueDir = opts.Dir
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}

	updates := make(map[string]string)
	updates["hello"] = "world"
	updates["yaojun"] = "1995"
	updates["yaoming"] = "1980"
	// 写入数据
	fmt.Println("insert========")
	for k, v := range updates {
		err := db.Update(func(txn *badger.Txn) error {
			err := txn.Set([]byte([]byte(k)), []byte(v))
			return err
		})
		if err != nil {
			fmt.Println("write failed")
		}
	}
	fmt.Println("point query========")
	// badger读出数据
	_ = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("yaojun"))
		if err != nil {
			return err
		}
		val, err := item.Value()
		if err != nil {
			return err
		}
		fmt.Printf("The answer is: %s\n", val)
		return nil
	})
	fmt.Println("iterator========")
	// 遍历
	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			v, err := item.Value()
			if err != nil {
				return err
			}
			fmt.Printf("key=%s, value=%s\n", k, v)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
	fmt.Println("prefix iterator========")
	// 前缀遍历
	db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte("yao")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			v, err := item.Value()
			if err != nil {
				return err
			}
			fmt.Printf("key=%s, value=%s\n", k, v)
		}
		return nil
	})
	// 只遍历key
	fmt.Println("key-only iterator========")
	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			fmt.Printf("key=%s\n", k)
		}
		return nil
	})

	defer db.Close()
}
