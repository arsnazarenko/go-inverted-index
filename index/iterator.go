package index

import (
	"cmp"
	"slices"
)

type Iterator struct {
	index *InvertedIndex
	pl    PostingList
	pos   int
}

func (it *Iterator) HasNext() bool {
	return it.pos < len(it.pl)
}

func (it *Iterator) Next() {
	it.pos++
}

func (it *Iterator) Get() DocumentID {
	return it.pl[it.pos]
}

func Or(it1, it2 *Iterator) *Iterator {
	res := mergeSlices(it1.pl, it2.pl)

	return &Iterator{
		index: it1.index,
		pl:    res,
		pos:   0,
	}
}

func And(it1, it2 *Iterator) *Iterator {
	res := intersectSlices(it1.pl, it2.pl)
	return &Iterator{
		index: it1.index,
		pl:    res,
		pos:   0,
	}

}

func Not(it *Iterator) *Iterator {
	fullDocsPl := make(PostingList, 0, len(it.index.docs))
	for _, d := range it.index.docs {
		fullDocsPl = append(fullDocsPl, d.DID)
	}
	res := subtractSlices(it.pl, fullDocsPl)
	return &Iterator{
		index: it.index,
		pl:    res,
		pos:   0,
	}
}

func mergeSlices[T cmp.Ordered](s1, s2 []T) []T {
	// Создаем карту для отслеживания уже добавленных элементов
	seen := make(map[T]bool)
	// Создаем новый слайс для результата
	// Добавляем все элементы из всех слайсов в карту и результат
	for _, v := range s1 {
		seen[v] = true
	}

	for _, v := range s2 {
		if !seen[v] {
			s1 = append(s1, v)
			seen[v] = true
		}
	}
	// Сортируем результирующий слайс
	slices.Sort(s1)
	return s1
}

func intersectSlices[T cmp.Ordered](slice1, slice2 []T) []T {
	// Создаем карту для отслеживания элементов из первого слайса
	seen := make(map[T]bool)
	for _, v := range slice1 {
		seen[v] = true
	}
	// Создаем новый слайс для результата
	result := make([]T, 0)
	// Итерируем по второму слайсу, добавляя элементы в результат, если они присутствуют в карте
	for _, v := range slice2 {
		if seen[v] {
			result = append(result, v)
		}
	}
	slices.Sort(result)
	// Сортируем результирующий слайс
	return result
}

func subtractSlices[T cmp.Ordered](slice1, slice2 []T) []T {
	// Создаем карту для отслеживания элементов из второго слайса
	seen := make(map[T]bool)
	var small, big []T = slice1, slice2
	if len(slice1) > len(slice2) {
		small = slice2
		big = slice1

	}
	for _, v := range small {
		seen[v] = true
	}
	// Создаем новый слайс для результата
	result := make([]T, 0)
	// Итерируем по первому слайсу, добавляя элементы в результат, если они не присутствуют в карте
	for _, v := range big {
		if !seen[v] {
			result = append(result, v)
		}
	}
	slices.Sort(result)
	return result
}
