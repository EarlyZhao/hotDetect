package hotDetect

import "sync/atomic"

type ring struct {
	buckets []*freqCount
	size int
	currentBucketIndex int32
}

func newRing(windowSize int, limitItemCount int64, id string) *ring {
	r := &ring{size: windowSize}
	r.buckets = make([]*freqCount, r.size)
	for offset := range r.buckets{
		r.buckets[offset] = newFc(limitItemCount, id)
	}

	return r
}


func (r *ring) sliding(){
	// 理论上只有一个gorouting在修改
	var indexShouldBe int32
	currentIndex := r.getCurrentBucketIndex()
	if currentIndex == int32(len(r.buckets) - 1) {
		indexShouldBe = 0
	}else{
		indexShouldBe = currentIndex + 1
	}

	r.buckets[indexShouldBe].Reset()
	atomic.StoreInt32(&r.currentBucketIndex, indexShouldBe)
}

func (r *ring) lastWindow() *freqCount {
	var lastWindow int32
	current := r.getCurrentBucketIndex()
	if current == 0 {
		lastWindow = int32(len(r.buckets) - 1)
	}else {
		lastWindow = current - 1
	}

	return r.buckets[lastWindow]
}

func (r *ring) currentWindow() *freqCount {
	current := r.getCurrentBucketIndex()
	return r.buckets[current]
}

func (r *ring) lastWindowTop(topN int) []TopItem{
	return r.lastWindow().TopKey(topN)
}

func (r *ring) inrc(key interface{}, num int64) {
	r.currentWindow().Add(key, num)
}

func (r *ring) getCurrentBucketIndex() int32{
	return atomic.LoadInt32(&r.currentBucketIndex)
}