package hotDetect

import (
	"math/rand"
	"sort"
	"sync"
	"sync/atomic"
)

type freqCount struct {
	fu map[interface{}]int64
	limit int64
	size int64
	id string

	mutex sync.Mutex

	dropCount int64
	rejectCont int64
}


func newFc(limit int64, id string) *freqCount {
	return &freqCount{
		fu: make(map[interface{}]int64),
		limit: limit,
		id: id,
	}
}


func(l *freqCount) Add(id interface{}, num int64){
	needSkp := false
	// 达到了限制的60%
	// 先粗暴丢弃50%
	// todo： 丢弃比例和超限比例同步上涨
	// 先除,精度丢失比较大
	// 理论上有可能溢出，实际上?
	if l.getSize() * 10 / 6 >= l.limit{
		randVal := rand.Intn(2)
		if randVal > 0{
			needSkp = true
		}
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	if _, ok := l.fu[id];ok {
		l.fu[id] += num
	}else{
		// 超过限制直接丢弃
		if l.getSize() >= l.limit {
			l.rejectCont ++
			return
		}
		// 有概率地丢弃
		if needSkp{
			l.dropCount ++
			return
		}

		l.fu[id] = num
		atomic.AddInt64(&l.size, 1)
	}

}

func (l *freqCount) Reset(){
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.fu = make(map[interface{}]int64)
	atomic.StoreInt64(&l.size, 0)
	l.dropCount = 0
	l.rejectCont = 0
}
//
func (l *freqCount) TopKey(n int) (list []TopItem){
	if l.getSize() == 0 || n <= 0 {
		return
	}

	alllist := l.getSlice()
	//取出top-n
	list = priorityQueue(alllist, n)
	// 从大到小排序
	sort.Slice(list, func(i, j int) bool {
		return list[i].Freq > list[j].Freq
	})
	return
}

func (l *freqCount) getSlice() []TopItem{
	alllist := make(itemSlice, 0, l.getSize())
	l.mutex.Lock()
	for key, value := range l.fu {
		alllist = append(alllist, TopItem{ Key: key, Freq: value })
	}
	l.mutex.Unlock()

	return alllist
}

func (l *freqCount) getSize() int64 {
	return atomic.LoadInt64(&l.size)
}