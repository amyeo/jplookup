package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"

	bolt "go.etcd.io/bbolt"
)

func initDB() {
	err := os.Remove(DBName)
	if err != nil {
		fmt.Printf(" (warn) error: %v\n", err)
	}
	db, err := bolt.Open(DBName, 0666, nil)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("ENTRY"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("INDEX3"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	file, err := os.Open("JMdict")
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	jmd := JMdict{}

	d := xml.NewDecoder(bytes.NewReader([]byte(data)))
	d.Strict = false //needed or else error at lines with "&unc;""

	err = d.Decode(&jmd)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	indexMap := new(CharIndex)
	indexMap.Char = make(map[rune]CharIndexMemberList)

	jmdEntryIndex := make(map[int]Entry)

	for _, entry := range jmd.Entry {
		sub_entry := 0
		fmt.Printf("E: %d\n", entry.Sequence)
		jmdEntryIndex[entry.Sequence] = entry
		for _, kanji := range entry.Kanji {
			kanji_runes := []rune(kanji.Kanji)
			for i, kanjiChar := range kanji_runes {
				//fmt.Printf("%c\n", kanjiChar)
				char_index := indexMap.Char[kanjiChar]
				char_index.AddItem(entry.Sequence, i, sub_entry, len(kanji_runes))
				indexMap.Char[kanjiChar] = char_index
			}
			sub_entry++
		}
		for _, kana := range entry.Kana {
			kana_runes := []rune(kana.Kana)
			for i, kanaChar := range kana_runes {
				char_index := indexMap.Char[kanaChar]
				char_index.AddItem(entry.Sequence, i, sub_entry, len(kana_runes))
				indexMap.Char[kanaChar] = char_index
			}
			sub_entry++
		}
	}

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ENTRY"))
		for _, entry := range jmd.Entry {
			jsonStr, err := json.Marshal(entry)
			if err != nil {
				fmt.Printf("error: %v", err)
				return err
			}
			b.Put([]byte(strconv.Itoa(entry.Sequence)), []byte(jsonStr))
		}
		return nil
	})

	//new new indexing method
	for key, val := range indexMap.Char {
		index := fmt.Sprintf("%d", key)
		fmt.Printf("(index v3) Key %s\n", index)

		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("INDEX3"))
			jsonStr, err := json.Marshal(val)
			if err != nil {
				fmt.Printf("error: %v", err)
				return err
			}
			return b.Put([]byte(index), jsonStr)
		})
	}

	fmt.Println("All indexing finished.")
}
