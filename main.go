package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"

	. "ws-im/cmd"
)

var (
	addr    = flag.String("addr", ":8080", "http service address")
	debug   = flag.String("pprof", "", "type -pprof to address for pprof http")
	G_house map[string]*Room
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()
	if x := *debug; x != "" {
		log.Printf("starting pprof server on %s", x)
		go func() {
			log.Printf("pprof server error: %v", http.ListenAndServe(x, nil))
		}()
	}
	G_house = make(map[string]*Room, 100)

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		log.Printf("room:%v",query.Get("room"))
		rid:=query.Get("room")
		if  rid=="" {
			rid="0"
		}
		// log.Printf("uid:%v",query.Get("uid"))
		uid:=query.Get("uid")
		if  uid=="" {
			w.WriteHeader(http.StatusBadRequest)
		}else{
			proom,ok := G_house[rid]
			if !ok {
				proom = NewRoom()
				go proom.Run()
				G_house[rid]=proom
			}
			ServeWs(proom, w, r, uid)
		}


	})

	// 服务器状态打印
	go G_Stats.Print(G_Ticker3s)

	log.Print("Listen Port",*addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
