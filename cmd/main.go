package main

import (
	"fmt"
	"inverted-index/index"
)

func main() {
	// Набор документов
	documents := []index.Document{
		{DID: 0, Text: "This is a sample document."},
		{DID: 1, Text: "Another document with sample words."},
		{DID: 2, Text: "This is a third document."},
		{DID: 3, Text: "This is a new document for test"},
	}
	inmem, err := index.NewInvertedIndex()
	if err != nil {
		panic(err)
	}
	for _, d := range documents {
		fmt.Printf("%d: %s\n", d.DID, d.Text)
		inmem.AddDocument(d)
	}

	// Строим индекс в памяти
	//index, err := index.NewInMemoryIndexFrom(documents)
	//if err != nil {
	//    panic(err)
	//}

	// Выводим индекс
	fmt.Println(inmem)

	if err = inmem.Save("/tmp/inverted.index"); err != nil {
		panic(err)
	}
	s, _ := index.LoadFromFile("/tmp/inverted.index", documents)
	fmt.Print("this: ")
	it := s.Search("this")
	for i := it; i.HasNext(); i.Next() {
		fmt.Printf("%d ", i.Get())
	}
	fmt.Println()

	fmt.Print("another: ")
	it = s.Search("another")
	for i := it; i.HasNext(); i.Next() {
		fmt.Printf("%d ", i.Get())
	}
	fmt.Println()
	fmt.Print("is: ")
	it = s.Search("is")
	for i := it; i.HasNext(); i.Next() {
		fmt.Printf("%d ", i.Get())
	}
	fmt.Println()
	resIt := index.And(index.And(s.Search("this"), s.Search("is")), s.Search("new"))
	fmt.Printf("this && is && new: ")
	for i := resIt; i.HasNext(); i.Next() {
		fmt.Printf("%d ", i.Get())
	}
	fmt.Println()
	newIt := index.And(index.Or(s.Search("this"), s.Search("another")), s.Search("a"))
	fmt.Printf("(this || another) && a: ")
	for i := newIt; i.HasNext(); i.Next() {
		fmt.Printf("%d ", i.Get())
	}
	fmt.Println()
	newnewIt2 := index.Or(index.And(s.Search("this"), s.Search("another")), s.Search("a"))
	fmt.Printf("(this && another) || a: ")
	for i := newnewIt2; i.HasNext(); i.Next() {
		fmt.Printf("%d ", i.Get())
	}
	fmt.Println()

}
