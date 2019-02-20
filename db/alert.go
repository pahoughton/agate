/* 2019-02-13 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package db


const (
	BNameFmt = "2006-01-02"  // buckets named by alert date
	queueBucket = "tqueue"
)


type AlertTicket struct {
	Tid		string
	Qid		uint64
}


func (adb *AgateDB) TicketCleanBuckets() {

	minDate := time.Now().AddDate(0,0,adb.maxDays * -1).Format(BNameFmt)

	fmt.Println("INFO cleaning buckets before ",minDate)

	var delList []string

	err := adb.db.View(func(tx *bolt.Tx) error {

		err := tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			if strings.Compare(string(name),minDate) < 0 {
				delList = append(delList,string(name))
			}
			return nil
		})
		return err
	})
	if err != nil {
		fmt.Println("FATAL reading buckets ",err.Error())
		return
	}
	err = adb.db.Update(func(tx *bolt.Tx) error {
		for _, bname := range delList {
			if err := tx.DeleteBucket([]byte(bname)); err != nil {
				fmt.Println("ERROR delete bucket ",bname," - ",err.Error())
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("FATAL deleting buckets ",err.Error())
		return
	}
}

FIXME - start time part of key

func (adb *AgateDB) TicketAdd(
	start	time.Time,
	fp		uint64,
	tid		string) error {

	bname := start.Format(BNameFmt)

	aKey := make([]byte, binary.MaxVarintLen64)

	binary.PutUvarint(aKey, fp)

	err := adb.db.Update(func(tx *bolt.Tx) error {

		bkt, err := tx.CreateBucketIfNotExists([]byte(bname))
		if err != nil {
			return err
		}
		return bkt.Put(aKey,[]byte(tid))
	})
	return err
}

func (adb *AgateDB) Ticket(start time.Time,fp uint64) (string, error) {

	bname := start.Format(BNameFmt)

	aKey := make([]byte, binary.MaxVarintLen64)

	binary.PutUvarint(aKey, fp)

	var tid string

	err := adb.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bname))
		if bkt == nil {
			return errors.New("bucket not found " + bname)
		}
		val := bkt.Get(aKey)
		if val == nil {
			return fmt.Errorf("alert not found: %u", aKey)
		}
		copy(tid,string(val))
		return nil
	})
	return tid, err
}

func (adb *AgateDB) TicketDelete(start time.Time,fp uint64) error {

	bname := start.Format(BNameFmt)

	aKey := make([]byte, binary.MaxVarintLen64)

	binary.PutUvarint(aKey, fp)

	err := adb.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bname))
		if bkt == nil {
			return errors.New("bucket not found " + bname)
		}
		return bkt.Delete(aKey)
	})
	return err
}
