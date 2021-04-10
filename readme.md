
# 自动探测热点
	- 一个自动探测热点的SDK【本地探测】
	- 接口埋点统计，埋点对象可以是数字、string等
	- 通过滑动窗口+频率计数+优先队列计算Top-K，K由业务方参数决定
	- 通过业务回调回传上一次窗口热点ID TOP-K，窗口时间由业务方参数决定
## 项目简介 

   [ ![SDK原理图](https://i0.hdslb.com/bfs/live/32f6e3d5adddd865a7d1a72a315c5006c2345ed9.png)]


# 接入案例
    go-live/app/service/xroom 房间信息服务。可查看NewDetect调用的地方
    
# 代码示例

创建一个热点探测实例：
```
import "go-live/pkg/library/hotDetect"

func callback((list []hotDetect.TopItem)) {
    // hotDetect会将top-k的热点通过该回调传回来
    // 在该回调中处理热点ID
    // 可参考 xroom refreshHotDetectRoom 的实现
}

detectSomeThing := hotDetect.NewDetect(hotDetect.NewConfig(windowSize int, // 窗口数量
                                                        slidingTime int, // 窗口保持时间
                                                        limitItem int64, // 窗口探测数量上限，预防极端情况
                                                        id string,   // 自定义标识, 日志/埋点上报用
	                                                    sampling int, // 采样率
	                                                    callback, // 回调
                                                        topNum int) // top-k  k的设置

//也可以用默认配置 hotDetect.NewDetect(hotDetect.DefualtConf("someID", callback))
```

埋点：
```
detectSomeThing.Record(123)
```
# 相关资料
    https://info.bilibili.co/pages/viewpage.action?pageId=154602332
