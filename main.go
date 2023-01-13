package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"os"
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
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if createTable(db) != nil {
		log.Fatal(err)
	}

	// Read files
	files, err := readDir(read_path)
	if err != nil {
		log.Fatal(err)
	}

	// Hash files
	for _, file := range files {
		hash := getHashOfFileAsStream(file)
		fmt.Println(file, hex.EncodeToString(hash))

		entry := FileEntry{
			name: file,
			size: 0, // FIXME:
			hash: hex.EncodeToString(hash),
		}

		if insert(db, entry) != nil {
			log.Fatal(err)
		}
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
