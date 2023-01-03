package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	bolt "go.etcd.io/bbolt"
)

var (
	read_path = "C:\\Users\\baltha\\Downloads\\Video"
	out_path  = ""
)

var (
	DBFile = "mvdup.bolt"
	Bucket = []byte("HashValues")
)

func main() {

	// Open database
	db, err := bolt.Open(DBFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) (err error) {
		if _, err = tx.CreateBucketIfNotExists(Bucket); err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	// Read files
	files, err := readDir(read_path)
	if err != nil {
		log.Fatal(err)
	}

	// Hash files
	for _, file := range files {
		hash := getHashOfFileAsStream(file)
		fmt.Println(file, hex.EncodeToString(hash))

		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(Bucket)
			v := b.Get(hash)
			fmt.Printf("The answer is: %s\n", v)
			return nil
		})
		return
	}
	fmt.Printf("found %d files\n", len(files))

}

func getHashOfFileAsStream(filep string) []byte {
	file, err := os.Open(filep)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// stat, _ := file.Stat()
	// stat.Size()

	sha256 := sha256.New()
	for {
		bytes := make([]byte, 4096)
		n, err := file.Read(bytes)
		if n == 0 || err != nil {
			break
		}
		sha256.Write(bytes[0:n])
	}

	return sha256.Sum(nil)
}
