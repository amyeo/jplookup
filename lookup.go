package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/fatih/color"
	bolt "go.etcd.io/bbolt"
)

func Lookup(second string) {
	db, err := bolt.Open(DBName, 0666, nil)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	defer db.Close()

	var final_index = new(Index)
	final_index.Items = make(map[rune]CharIndexMemberList)

	if second == "" {
		cyan := color.New(color.FgCyan).PrintfFunc()
		cyan("Lookup > ")
		fmt.Scanln(&second)
	}
	highlight_set := make(map[rune]bool)

	var results = make(map[string][]IndexItem)
	char_step := 1
	wildcard_skip := false

	for input_string_index, char := range second {
		var temp_results = make(map[string][]IndexItem)
		if len(final_index.Items[char].Members) == 0 {
			//load it. it means its not initialized yet
			db.View(func(tx *bolt.Tx) error {
				// Assume bucket exists and has keys
				b := tx.Bucket([]byte("INDEX3"))

				str_index := fmt.Sprintf("%d", char)
				jstring := b.Get([]byte(str_index))
				if jstring != nil {
					var temp_member_list CharIndexMemberList
					json.Unmarshal(jstring, &temp_member_list)
					final_index.Items[char] = temp_member_list
				}

				return nil
			})
		}

		if char == '?' || char == '？' {
			char_step = char_step + 1
		} else if char == '*' || char == '＊' {
			wildcard_skip = true
		} else {
			highlight_set[char] = true
			for _, member := range final_index.Items[char].Members {
				if input_string_index == 0 {
					temp_results[fmt.Sprintf("%d.%d", member.ReferenceSequence, member.SubEntryId)] = append(temp_results[fmt.Sprintf("%d.%d", member.ReferenceSequence, member.SubEntryId)], member)
				} else {
					for _, result_member := range results[fmt.Sprintf("%d.%d", member.ReferenceSequence, member.SubEntryId)] {
						if member.Position == (result_member.Position+char_step) || (wildcard_skip && (member.Position > result_member.Position)) {
							temp_results[fmt.Sprintf("%d.%d", member.ReferenceSequence, member.SubEntryId)] = append(temp_results[fmt.Sprintf("%d.%d", member.ReferenceSequence, member.SubEntryId)], member)
						}
					}
				}
			}
			results = temp_results
			char_step = 1
			wildcard_skip = false
		}

	}
	fmt.Println()
	fmt.Println("Results:")
	fmt.Println("===================\n")
	var final_results = []IndexItem{}
	for _, row := range results {
		for _, result_item := range row {
			final_results = append(final_results, result_item)
		}
	}

	sort.SliceStable(final_results, func(i, j int) bool {
		return final_results[i].RuneLength < final_results[j].RuneLength
	})

	for _, result_item := range final_results {
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("ENTRY"))
			v := b.Get([]byte(strconv.Itoa(result_item.ReferenceSequence)))
			var entry Entry
			json.Unmarshal(v, &entry)
			PrintEntry(&entry, highlight_set, result_item.SubEntryId)
			return nil
		})
		fmt.Println()
	}

}
