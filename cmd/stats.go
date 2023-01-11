package cmd

import (
	"fmt"
	"sync/atomic"
	"time"

)

var (
	// 服务器状态打印的定时器
	G_Ticker3s *time.Ticker = time.NewTicker(time.Second * 3)
	G_Stats *stats = &stats{}
)
// 自增ID的原子化写法
func AutoAdd(x *uint64) uint64 {
	// int32的最大取值为2147483647,之后就需要重置归0,否则会越界
	// CompareAndSwapInt32为原子化的交换操作
	atomic.CompareAndSwapUint64(x,2147483647,0)
	// AddInt32为原子化的自增操作
	atomic.CompareAndSwapUint64(x,*x,atomic.AddUint64(x, 1))
	return *x
}


// 状态统计结构体
type stats struct {
	clients    int64
	clientsMax int64
	lastmsgs   int64
	idsMax	uint64	//客户端最大的下标
}
// AddInt64为原子化的自增操作
func (s *stats) messageAdd()      { atomic.AddInt64(&s.lastmsgs, 1) }
func (s *stats) clientConnect()    { atomic.AddInt64(&s.clients, 1) }
func (s *stats) clientDisconnect() {
	if s.clients <= 0 {
		return
	}
	atomic.AddInt64(&s.clients, -1)
}
func (s *stats) IdsMaxAdd() uint64 {
	return AutoAdd(&s.idsMax)
}


func (s *stats) Print(t *time.Ticker) {

	for {
		<-t.C
		// 原子化的载入数据
		clients := atomic.LoadInt64(&s.clients)
		clientsMax := atomic.LoadInt64(&s.clientsMax)
		lastmsgs := atomic.LoadInt64(&s.lastmsgs)
		if clients > clientsMax {
			clientsMax = clients
			// 原子化的存储操作
			atomic.StoreInt64(&s.clientsMax, clientsMax)
		}
		fmt.Printf("\r%v 消息读取数:%v,客户端连接数:%v,客户端连接峰值:%v             ... ", time.Now().Format("15:04:05.000000"),lastmsgs,clients,clientsMax)
		// fmt.Printf("\r消息处理数:%v,客户端连接数:%v,客户端连接峰值:%v,客户端最大下标:%v ... ", lastmsgs,clients,clientsMax,s.idsMax)
		// G_Server.clients.Print()
	}


}