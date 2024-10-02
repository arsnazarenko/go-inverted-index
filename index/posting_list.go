package index

import (
	"encoding/binary"
)

// Постинг лист - список документов, содержащих данный терм
type DocumentID int
type PostingList []DocumentID
// Кодирование PostingList с помощью pForDelta
func (pl PostingList) Encode() []byte {
	var (
		b   [binary.MaxVarintLen64]byte
		res []byte
		n   int = 0
	)

	n = binary.PutUvarint(b[:], uint64(pl[0]))
	res = append(res, b[:n]...)

	// Сжатие оставшихся ID с помощью pForDelta
	for i := 1; i < len(pl); i++ {
		delta := pl[i] - pl[i-1]
		n = binary.PutUvarint(b[:], uint64(delta))
		res = append(res, b[:n]...)
	}
	return res
}

// Декодирование PostingList с помощью pForDelta
func (pl *PostingList) Decode(data []byte) int{
    if len(data) == 0 {
        return -1
    }
	var tmpId uint64 = 0
	var tmpN int = 0

	tmpId, tmpN = binary.Uvarint(data[0:])
	*pl = append(*pl, DocumentID(tmpId))
	// Декодирование оставшихся ID с помощью pForDelta
	var n int = tmpN
	var id uint64 = tmpId
	for n < len(data) {
		tmpId, tmpN = binary.Uvarint(data[n:])
		n += tmpN
		id += tmpId
		*pl = append(*pl, DocumentID(id))
	}
    return n
}
