package main

import (
	"flag"
	"fmt"
	"runtime"
	"sync"
	"time"
)

var (
	url = flag.String("url", "ws://127.0.0.1:8080", "ws服务器的URL地址")
	rooms = flag.Uint("rooms",100,"要创建的房间数量")
	ids = flag.Uint("ids",10,"每个房间要连接的客户端数量")
	s = flag.Uint64("s",60,"每个客户端连接时长(秒)")

	G_online = make(map[RoomId]*ClintInRoom, 100)

	)

type RoomId string
type ClintInRoom struct {
	clients []*client
	rwLock sync.RWMutex
}


func NewRoomId(id int) RoomId {
	return RoomId(fmt.Sprintf("u%d",id))
}
func NewClintInRoom(num int) *ClintInRoom {
	return &ClintInRoom{
		clients: make([]*client, 0, num),
	}
}
func (this *ClintInRoom) addClient(c *client)  {
	this.rwLock.Lock()
	defer this.rwLock.Unlock()

	this.clients=append(this.clients, c)

}
func main() {
	flag.Parse()

	var ss=time.Second*time.Duration(int64(*s))
	for r := 1; r <= int(*rooms); r++ {
		for i := 0; i < int(*ids); i++ {
			path := fmt.Sprintf("/ws?room=r%v&uid=u0%v",  r, i)
			room := NewClintInRoom(100)
			room.addClient( Do(*url, path,ss))
			G_online[NewRoomId(r)] =room
			fmt.Printf("=============>>>>> inited %v\n",path)
			runtime.Gosched()
		}


	}
}
