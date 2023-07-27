package db
//fix
import (
	"log"

	"github.com/boltdb/bolt"
)

var Db *bolt.DB

func init() {
	var err error
	Db, err = bolt.Open("./lm.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
}

//参数读取
func Read(bucketname string, key string) (bs []byte) {
	Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketname))
		bs = b.Get([]byte(key))
		return nil
	})
	return bs
}

//参数设置
func Update(bucketname string, key string, value []byte) {
	Db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketname))
		if err != nil {
			return err
		}
		return b.Put([]byte(key), value)
	})
}
