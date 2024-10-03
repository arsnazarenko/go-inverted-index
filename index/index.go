package index

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

// Постинг лист - список документов, содержащих данный терм
// Индекс в памяти
type InvertedIndex struct {
	docs  []Document
	index map[string]PostingList
}

type Document struct {
	DID  DocumentID
	Text string
}

func NewInvertedIndex() (*InvertedIndex, error) {
	return &InvertedIndex{
		docs:  []Document{},
		index: make(map[string]PostingList),
	}, nil
}

// Строим индекс в памяти
func NewInvertedIndexFrom(docs []Document) (*InvertedIndex, error) {
	if len(docs) == 0 {
		return nil, fmt.Errorf("docs len in 0")
	}

	inmem, err := NewInvertedIndex()
	if err != nil {
		return nil, err
	}

	for _, d := range docs {
		inmem.AddDocument(d)
	}
	return inmem, nil
}

// Нормализация текста

func (in *InvertedIndex) Search(term string) *Iterator {
	if v, ok := in.index[term]; ok {
		return &Iterator{
			index: in,
			pl:    v,
			pos:   0,
		}
	}
	return &Iterator{
		index: in,
		pl:    PostingList{},
		pos:   0,
	}
}

func (in *InvertedIndex) AddDocument(doc Document) error {
	normalizedText := normalizeText(doc.Text)
	in.docs = append(in.docs, doc)
	// Разбиение текста на токены
	tokens := strings.Fields(normalizedText)
	for _, token := range tokens {
		// Добавление документа в постинг лист
		if _, ok := in.index[token]; ok {
			in.index[token] = append(in.index[token], doc.DID)
		} else {
			in.index[token] = PostingList{doc.DID}
		}
	}
	return nil
}

// |PAYLOAD|META|META_OFFSET(2B)|
//
// |len|token|posting_list|len|token|posting_list|len|token|posting_list|

func (in *InvertedIndex) Save(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	b := bytes.NewBuffer(make([]byte, 0, 4096))
	offsets := []uint16{}
	var tmp [binary.MaxVarintLen64]byte

	for token, pl := range in.index {
		offsets = append(offsets, uint16(b.Len()))
		n := binary.PutUvarint(tmp[:], uint64(len(token)))
		raw := pl.Encode()
		b.Write(tmp[:n])
		b.WriteString(token)
		b.Write(raw)
	}
	metaOffset := b.Len()
	for i := 1; i < len(offsets); i++ {
		binary.LittleEndian.PutUint16(tmp[:], offsets[i])
		b.Write(tmp[:2])
	}
	binary.LittleEndian.PutUint16(tmp[:], uint16(metaOffset))
	b.Write(tmp[:2])

	n, err := file.Write(b.Bytes())
	if err != nil || n != b.Len() {
		return fmt.Errorf("Write error: %w", err)
	}
	file.Sync()
	return nil
}

// |len|str|pl len|str|pl len|str|pl|...|2, 4|5|
// |          |          |          |
func LoadFromFile(path string, docs []Document) (*InvertedIndex, error) {
	var (
		res map[string]PostingList = make(map[string]PostingList)
	)

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	metaOffset := binary.LittleEndian.Uint16(raw[len(raw)-2:])

	payload := raw[:metaOffset]
	meta := raw[metaOffset : len(raw)-2]

	offsets := []int{0}
	for i := 0; i < len(meta); i += 2 {

		off := binary.LittleEndian.Uint16(meta[i : i+2])
		offsets = append(offsets, int(off))
	}
	offsets = append(offsets, int(metaOffset))
	for i := 0; i < len(offsets)-1; i++ {
		start, end := offsets[i], offsets[i+1]
		view := payload[start:end]
		idx := 0
		tokenLen, n := binary.Uvarint(view[idx:])
		idx += n

		token := string(view[idx : idx+int(tokenLen)])
		idx += int(tokenLen)
		var pl PostingList
		pl.Decode(view[idx:])
		res[token] = pl
	}
	return &InvertedIndex{
		docs:  docs,
		index: res,
	}, nil
}

func normalizeText(text string) string {
	// Удаляем знаки препинания и приводим к нижнему регистру
	r := strings.NewReplacer(".", "", ",", "", "!", "", "?", "", ":", "")
	return strings.ToLower(r.Replace(text))
}

func (in *InvertedIndex) String() string {
	b := strings.Builder{}
	for term, postingList := range in.index {
		b.WriteString(fmt.Sprintf("%s: %v\n", term, postingList))
	}
	return b.String()
}
