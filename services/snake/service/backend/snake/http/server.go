package http

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"snake/db"
	"snake/objects"
)

type Map struct {
	Secret string    `json:"secret"`
	Init   [256]byte `json:"init"`
	Flag   string    `json:"flag"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func home(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello"))
}

func create(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method)
	if r.Method == http.MethodPost {
		log.Println("Handling create")
		var _map Map
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&_map)
		if err != nil {
			errorResp(w, 500, err)
			return
		}

		if err != nil {
			log.Fatalln(err)
		}
		log.Println(_map)
		level := objects.Level{Id: 1, Secret: _map.Secret, Counter: 0, Init: _map.Init, Flag: _map.Flag}
		_, err = db.InsertDoc("test", "test", level)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("saved to db")
		// сохранение карты в базу
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		resp := make(map[string]string)
		resp["msg"] = "Created"
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}
	errorResp(w, 405, errors.New("method not allowed"))
}

func play(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		errorResp(w, 500, err)
	}
	gameConn := NewGameConn(conn, 1)
	go gameConn.Play()
}

func errorResp(w http.ResponseWriter, code int, err error) {
	resp := make(map[string]string)
	resp["msg"] = err.Error()
	jsonResp, _ := json.Marshal(resp)
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonResp)
}

func StartServ() {
	http.HandleFunc("/", home)
	http.HandleFunc("/create", create)
	http.HandleFunc("/play", play)
	http.ListenAndServe(":8080", nil)
}
