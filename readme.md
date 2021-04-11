
# 自动探测热点
- 一个自动探测热点的SDK【本地探测】
- 接口埋点统计，埋点对象可以是数字、string等，一般是被访问对象id
- 通过滑动窗口+频率计数+优先队列计算Top-K，K由业务方参数决定
- 通过业务回调定时回传上一次窗口热点ID TOP-K，窗口时间由业务方参数决定

可以用于解决高并发下的热点问题。
## 项目简介 
   ![SDK原理图](https://github.com/EarlyZhao/hotDetect/blob/master/doc/hotDetect.png?raw=true)

 - 基于SDK自动探测出热点ID  
 - 将热点ID上报，或者将对应数据缓存到本地
 - 接口优先从本地缓存读取数据，miss才从下游存储中获取
 - 这样热点数据大部分会在缓存中被命中，也就解决了热点问题

基础理论模型、实践样例、衍生问题等可参考下面【资料】。

# 代码示例

创建一个热点探测实例：
```
import "hotDetect"

func callback((list []hotDetect.TopItem)) {
    // hotDetect会将top-k的热点通过该回调定时传回来
    // 在该回调中处理热点ID，业务可以将热点数据缓存到本地内存，或自行上报
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

埋点统计：
```
detectSomeThing.Record(某id)
```
# 资料
   
[高并发下的热点处理实践]( https://ruby-china.org/topics/40596)
