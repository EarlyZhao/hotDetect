package hotDetect

import (
	"math/rand"
	"time"
)

type Detect struct {
	ring *ring
	channel chan interface{}
	conf *Config
}

type Config struct {
	// 有多少个窗口
	windowSize int
	// 窗口保持时间(多久滑动一次), 单位秒
	slidingTime int
	// 单个窗口探测的对象总量
	limitDetectItem int64
	// 标识
	id string
	// 采样率 最大100
	sampling int
	// 窗口滑动时回调，回传上一窗口topN
	// 也可以通过 LastWindowTop 自行拉取上一窗口topN
	callbackAfterSliding func(tops []TopItem)
	// 如果注册了 callbackAfterSliding 需要配置一个参数，拉取靠前的数量, 默认 top 10
	wantTopNum int
}

func NewConfig(windowSize int, slidingTime int, limitItem int64, id string,
	           sampling int, callbackEverySliding func(tops []TopItem), topNum int) *Config{

	return &Config{
		windowSize: windowSize,
		slidingTime: slidingTime,
		limitDetectItem: limitItem,
		id: id,
		sampling: sampling,
		callbackAfterSliding: callbackEverySliding,
		wantTopNum: topNum,
	}
}

func DefualtConf(id string, callback func(top []TopItem), ) *Config{
	return &Config{windowSize: 10, slidingTime: 3, limitDetectItem: 60000, id: id, sampling: 100, wantTopNum: 500,
		callbackAfterSliding: callback}
}

func NewDetect(conf *Config) (d *Detect) {
	r := newRing(conf.windowSize, conf.limitDetectItem, conf.id)

	d = &Detect{ring: r, channel: make(chan interface{}, conf.limitDetectItem), conf: conf}
	// 定时切窗口
	go func() {
		for{
			time.Sleep(time.Second * time.Duration(conf.slidingTime))
			d.sliding()
		}
	}()
	// 异步消费
	go d.comsumingChannel()

	return
}

func (d *Detect) LastWindowTop(topN int) (list []TopItem){
	list =  d.ring.lastWindowTop(topN)
	return
}
// 埋点上报
func (d *Detect) Record(key interface{}) {
	if d.conf.sampling < 100 {
		if rand.Intn(101) > d.conf.sampling{
			return
		}
	}
	select {
		case d.channel <- key:
		default:
			//metricTotalCount.Inc(itemAbandon, d.conf.id)
	}
}

func (d *Detect) sliding(){
	defer func() {
		if err := recover();err != nil {
			//log.Error("%v hotDetect panic ERROR when sliding : %v", d.conf.id, err)
			//metricTotalCount.Inc(errorCount, d.conf.id)
			return
		}
	}()
	// 切到下一个窗口
	d.ring.sliding()

	// 使用方回调
	if d.conf.callbackAfterSliding != nil {
		wantTopNum := d.conf.wantTopNum
		if wantTopNum == 0 { wantTopNum = 10 }

		list := d.LastWindowTop(wantTopNum)
		if len(list) > 0 {
		}
		go func() {
			defer func() {
				if err := recover();err != nil {
					//log.Error("%v hotDetect panic ERROR when callback : %v", d.conf.id, err)
					//metricTotalCount.Inc(errorCount, d.conf.id)
					return
				}
			}()
			d.conf.callbackAfterSliding(list)
		}()
	}
	// 其他埋点上报
}


func (d *Detect) comsumingChannel() {
	defer func() {
		if err := recover();err != nil {
			//log.Error("%v hotDetect panic ERROR when comsumingChannel : %v", d.conf.id, err)
			//metricTotalCount.Inc(errorCount, d.conf.id)
			return
		}
	}()

	for{
		select {
			case key := <- d.channel:
				d.ring.inrc(key, 1)
				//metricTotalCount.Inc(itemAdd, d.conf.id)
		}
	}
}