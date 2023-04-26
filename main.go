package main

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/fatih/color"
)

const DBName = "jisho.db"

type JMdict struct {
	XMLName xml.Name `xml:"JMdict"`
	Entry   []Entry  `xml:"entry"`
}

type Entry struct {
	XMLName  xml.Name `xml:"entry" json:"entry"`
	Sequence int      `xml:"ent_seq" json:"sequence"`
	Kanji    []Kanji  `xml:"k_ele" json:"kanji"`
	Kana     []Kana   `xml:"r_ele" json:"kana"`
	Sense    Sense    `xml:"sense" json:"sense"`
}

type Kanji struct {
	XMLName xml.Name `xml:"k_ele" json:""`
	Kanji   string   `xml:"keb" json:"kanji"`
}

type Kana struct {
	XMLName xml.Name `xml:"r_ele" json:""`
	Kana    string   `xml:"reb" json:"kana"`
}

type Sense struct {
	XMLName xml.Name `xml:"sense" json:"sense"`
	Gloss   []Gloss  `xml:"gloss" json:"gloss"`
}
type Gloss struct {
	Lang  string `xml:"lang,attr" json:"lang"`
	Gloss string `xml:",chardata" json:"gloss"`
}

type CharIndex struct {
	Char map[rune]CharIndexMemberList
}

type CharIndexMemberList struct {
	Members []IndexItem
}

type IndexItem struct {
	ReferenceSequence int
	Position          int
	SubEntryId        int
	RuneLength        int
}

type Index struct {
	Items map[rune]CharIndexMemberList
}

func (charIndexMemberList *CharIndexMemberList) AddItem(item int, position int, sub_entry int, rune_length int) {
	var temp_item IndexItem
	temp_item.ReferenceSequence = item
	temp_item.Position = position
	temp_item.SubEntryId = sub_entry
	temp_item.RuneLength = rune_length
	charIndexMemberList.Members = append(charIndexMemberList.Members, temp_item)
}

func (e Entry) String() string {
	return fmt.Sprintf("Entry kanji=%v kana=%v gloss=%v",
		e.Kanji, e.Kana, e.Sense)
}

func PrintHighlight(print_string string, highlight_set map[rune]bool) {
	color.NoColor = false
	red := color.New(color.FgRed).PrintfFunc()
	for _, chr := range []rune(print_string) {
		_, exist := highlight_set[chr]
		if exist {
			red("%c", chr)
		} else {
			fmt.Printf("%c", chr)
		}
	}
	fmt.Printf("\n")
}

func PrintEntry(entry *Entry, highlight_set map[rune]bool, highlight_sub_entry int) {
	fmt.Printf("ID: %d\n", entry.Sequence)
	sub_entry_counter := 0
	fmt.Println("Kanji (漢字):")
	for i := 0; i < len(entry.Kanji); i++ {
		if sub_entry_counter == highlight_sub_entry {
			PrintHighlight(entry.Kanji[i].Kanji, highlight_set)
		} else {
			fmt.Println(entry.Kanji[i].Kanji)
		}
		sub_entry_counter = sub_entry_counter + 1
	}

	fmt.Println()

	fmt.Println("Kana (かな):")
	for i := 0; i < len(entry.Kana); i++ {
		if sub_entry_counter == highlight_sub_entry {
			PrintHighlight(entry.Kana[i].Kana, highlight_set)
		} else {
			fmt.Println(entry.Kana[i].Kana)
		}
		sub_entry_counter = sub_entry_counter + 1
	}

	fmt.Println()

	fmt.Println("Meaning (ENG / 英語):")
	for i := 0; i < len(entry.Sense.Gloss); i++ {
		if entry.Sense.Gloss[i].Lang == "" {
			fmt.Println(entry.Sense.Gloss[i].Gloss)
		}
	}

	fmt.Printf("\n---------------------\n")
}

func main() {
	lookup_word := ""

	argsWithoutProg := os.Args[1:]
	for _, arg := range argsWithoutProg {
		if arg == "--init" {
			fmt.Println("INIT db command invoked.")
			initDB()
			return
		} else {
			lookup_word = arg
		}
	}

	Lookup(lookup_word)
}
