package hotDetect

import (
	"container/heap"
)

type TopItem struct {
	Freq int64
	Key interface{}
}


type itemSlice []TopItem
func (h itemSlice) Len() int           { return len(h) }
func (h itemSlice) Less(i, j int) bool { return h[i].Freq < h[j].Freq }
func (h itemSlice) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h itemSlice) Pop() interface{}   { l := h[0]; h = h[1:h.Len()]; return l }
func (h itemSlice) Push(i interface{}) { h = append(h, i.(TopItem)) }

// 优先队列取top
func priorityQueue(s itemSlice, n int) (list itemSlice) {
	if n > len(s) {
		n = len(s)
	}
	list = make(itemSlice, 0, n)
	list = append(list, s[0:n]...)
	heap.Init(list)

	for index, val := range s {
		if index < n { continue }
		// 比顶部的都还小
		if list[0].Freq >= val.Freq {
			continue // 丢弃
		}

		// 换掉头部最小的
		// 下沉成小顶堆
		list[0] = val
		heap.Fix(list, 0)
 	}
	return
}
