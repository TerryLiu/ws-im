package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")
var G_house map[string]*Room

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
	G_house = make(map[string]*Room, 100)

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		log.Printf("room:%v",query.Get("room"))
		rid:=query.Get("room")
		if  rid=="" {
			rid="0"
		}
		log.Printf("uid:%v",query.Get("uid"))
		uid:=query.Get("uid")
		if  uid=="" {
			w.WriteHeader(http.StatusBadRequest)
		}else{
			proom,ok := G_house[rid]
			if !ok {
				proom = newRoom()
				go proom.run()
				G_house[rid]=proom
			}
			serveWs(proom, w, r, uid)
		}


	})
	log.Print("Listen Port",*addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
