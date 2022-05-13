package http

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"snake/dao"
	"snake/objects"
)

type Map struct {
	Secret string    `json:"secret"`
	Init   [256]byte `json:"init"`
	Flag   string    `json:"flag"`
}

type Ids struct {
	Ids []int `json:"ids"`
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
		if err != nil {
			log.Fatalln(err)
		}
		dao.SaveMap(objects.Level{
			Id:      0, //todo генерация id
			Secret:  _map.Secret,
			Counter: 0,
			Init:    _map.Init,
			Flag:    _map.Flag,
		})
		log.Println("saved to db")
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

func gameList(w http.ResponseWriter, r *http.Request) {
	//todo отдать все id игр либо как-то батчевать их
}

func play(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		errorResp(w, 500, err)
		return
	}
	msg := make(map[string]interface{})
	_ = conn.ReadJSON(msg)
	if msg["id"] != nil {
		switch msg["id"].(type) {
		case int:
			gameConn := NewGameConn(conn, msg["id"].(int))
			go gameConn.Play()
		default:
			errorResp(w, 401, errors.New("can't parse id"))
			conn.Close()
		}
		return
	}
	errorResp(w, 401, errors.New("can't find id"))
	conn.Close()
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
	http.HandleFunc("/gameList", gameList)
	http.HandleFunc("/play", play)
	http.ListenAndServe(":8080", nil)
}
