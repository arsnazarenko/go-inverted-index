package index

import (
	"testing"

	"github.com/stretchr/testify/require"
)


func TestIndexSearch(t *testing.T) {
	documents := []Document{
        {DID: 0, Text: "This is a sample document."},
        {DID: 1, Text: "Another document with sample words."},
        {DID: 2, Text: "This is a third document."},
        {DID: 3, Text: "This is a new document for test"},
	}
    inmem, err := NewInMemoryIndex()
    if err != nil { panic(err) }
    for _, d := range documents {
        inmem.AddDocument(d)
    }
    cases := [] struct {
        index *InvertedIndex
        name string
        req func(*InvertedIndex) *Iterator
        expected []DocumentID
    } {
        {
            index: inmem,
        	name: "is",
        	req: func(in *InvertedIndex) *Iterator {
                return in.Search("is")
        	},
        	expected: []DocumentID{0, 2, 3},
        },
        {
            index: inmem,
        	name: "!is",
        	req: func(in *InvertedIndex) *Iterator {
                return Not(in.Search("is"))
        	},
        	expected: []DocumentID{1},
        },
        {
            index: inmem,
        	name: "another",
        	req: func(in *InvertedIndex) *Iterator {
                return in.Search("another")
        	},
        	expected: []DocumentID{1},
        },
        {
            index: inmem,
        	name: "this",
        	req: func(in *InvertedIndex) *Iterator {
                return in.Search("this")
        	},
        	expected: []DocumentID{0, 2, 3},
        },
        {
            index: inmem,
        	name: "this && is && new",
        	req: func(in *InvertedIndex) *Iterator {
                return And(And(in.Search("this"), in.Search("is")), in.Search("new"))
        	},
        	expected: []DocumentID{3},
        },
        {
            index: inmem,
        	name: "(this || another) && a",
        	req: func(in *InvertedIndex) *Iterator {
                return And(Or(in.Search("this"), in.Search("another")), in.Search("a"))
        	},
        	expected: []DocumentID{0, 2, 3},
        },
        {
            index: inmem,
        	name: "(this && another) || a",
        	req: func(in *InvertedIndex) *Iterator {
                return Or(And(in.Search("this"), in.Search("another")), in.Search("a"))
        	},
        	expected: []DocumentID{0, 1, 2, 3},
        },

    }
    for _, tt := range cases {
        t.Run("Case " + tt.name, func(t *testing.T) {
            it := tt.req(tt.index)
            cnt := 0
            for i := it; i.HasNext(); i.Next() {
                require.Equal(t, tt.expected[cnt], i.Get())
                cnt++
            }
        })
    }



}


func TestIndexSaveLoad(t *testing.T) {

}

func TestIndexCreate(t *testing.T) {

}
