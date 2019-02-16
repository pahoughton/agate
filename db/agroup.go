/* 2019-02-13 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package go

const (
	agJsonBucket	"ag-json"
)

type AGroupRcvd struct {
	Json	[]byte
	Resolve	bool
}

func (adb *AgateDB) AgroupAdd(json []byte,resolve bool) {

	err := adb.db.Update(func(tx *bolt.Tx) error {

		bkt, err := tx.CreateBucketIfNotExists([]byte(agJsonBucket))
		if err != nil {
			panic(err)
		}
		key, err := bkt.NextSequence()
		if err != nil {
			panic(err)
		}
		keyBuf := make([]byte,binary.MaxVarintLen64)
		kn := binary.PutUvarint(keyBuf,key)

		rbyte byte
		if resolve {
			rbyte = 1
		} else {
			rbyte = 0
		}

		return bkt.Put(keyBuf[:kn],append(json,rbyte))

	})
    if err != nil {
		panic(err)
	}
}

func (adb *AgateDB) AGroupNext() *AGroupRcvd {

	var agrcv *AGroupRcvd
	agrcv = nil

	err := adb.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(agJsonBucket))
		if bkt == nil {
			return errors.New(agJsonBucket + " db bucket nil")
		}

		key, val := bkt.Cursor().First()
		if key == nil {
			return nil
		}
		if err := bkt.Delete(key); err != nil {
			return err
		}
		if bkt.Stats().KeyN == 0 {
			bkt.SetSequence(0)
		}
		agrcv = &AGroupRcvd{}

		rbyte = val[len(val)-1:1]
		if rbyte == 0 {
			agrcv.Resolve = false
		} else {
			agrcv.Resolve = true
		}
		copy(agrcv.Json,val[:len(val)-1])
		return nil
	})
	if err != nil {
		panic(err)
	}
	return agrcv
}
