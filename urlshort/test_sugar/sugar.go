package test_sugar

import bolt "go.etcd.io/bbolt"

func FillBucket(dbFile, bucketName string, data map[string]string) error {
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return err
	}

	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		for k, v := range data {
			err = b.Put([]byte(k), []byte(v))
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
