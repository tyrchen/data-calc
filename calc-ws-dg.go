package main

import (
	. "calcapp/calc"
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type ZgData struct {
	Inst Bpoint
	Zg   Point
	Gz   Point
	Gzmm Point
	Gf   Point
	Gfmm Point
	Gf1  Point
}

type DgData struct {
	Dg   Point
	Gz   Point
	Gzmm Point
	Gf   Point
	Gfmm Point
	Gf1  Point
}

type CalcData struct {
	Method string
	Pos    Value
	Dg     [2]DgData
	Zg     [2][ZG_NUM_SHOW]ZgData
}

type BpZgData struct {
	Method string
	Values [ZG_NUM_SHOW][COLS]Bpoint
}

var (
	values *BigData
)

func calcHandler(ws *websocket.Conn) {
	var err error

	// send initial data
	sendData(0, ws)
	sendBpZg(ws)

	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Can't receive")
			clear()
			break
		}

		fmt.Println("Received back from client: " + reply)

		runCommand(ws, reply)
	}
}

func wsSend(ws *websocket.Conn, msg string) {
	if err := websocket.Message.Send(ws, msg); err != nil {
		fmt.Println("Can't send")
		clear()
	}
}

func sendData(pos Value, ws *websocket.Conn) {
	b, _ := json.Marshal(getValues(pos))
	wsSend(ws, string(b))
}

func sendBpZg(ws *websocket.Conn) {
	var data BpZgData
	data.Method = "bpzg"
	for i := 0; i < ZG_NUM_SHOW; i++ {
		for j := 0; j < COLS; j++ {
			data.Values[i][j] = values.Data[i].Bp[j]
		}
	}
	b, _ := json.Marshal(data)
	wsSend(ws, string(b))
}

func runCommand(ws *websocket.Conn, reply string) {
	var v interface{}
	json.Unmarshal([]byte(reply), &v)
	m := v.(map[string]interface{})
	switch m["method"].(string) {
	case "calc":
		inst, pos := Bpoint(m["inst"].(float64)), Value(m["pos"].(float64))
		calc(inst, pos)
		sendData(pos, ws)

	case "close":
		clear()
		wsSend(ws, `{"method":"close","value":"true"}`)
	}
}

func calc(inst Bpoint, pos Value) {
	values.Run(inst, pos)
}

func getValues(pos Value) (ret CalcData) {
	var i Value
	ret.Method = "calc"
	ret.Pos = pos

	for i = 0; i < 2; i++ {
		ret.Dg[i].Dg = values.Dg[i+pos]
		ret.Dg[i].Gz = values.Gz[i+pos]
		ret.Dg[i].Gzmm = values.Gzmm[i+pos]
		ret.Dg[i].Gf = values.Gf[i+pos]
		ret.Dg[i].Gfmm = values.Gfmm[i+pos]
		ret.Dg[i].Gf1 = values.Gf1[i+pos]

		for j := 0; j < ZG_NUM_SHOW; j++ {
			ret.Zg[i][j].Inst = values.Data[j].Inst[i+pos]
			ret.Zg[i][j].Zg = values.Data[j].Zg[i+pos]
			ret.Zg[i][j].Gz = values.Data[j].Gz[i+pos]
			ret.Zg[i][j].Gzmm = values.Data[j].Gzmm[i+pos]
			ret.Zg[i][j].Gf = values.Data[j].Gf[i+pos]
			ret.Zg[i][j].Gfmm = values.Data[j].Gfmm[i+pos]
			ret.Zg[i][j].Gf1 = values.Data[j].Gf1[i+pos]
		}
	}

	return ret
}

func initValues() {
	values = new(BigData)
	values.Init()
}

func clear() {
	values.Clear()
	initValues()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <port>", os.Args[0])
		os.Exit(1)
	}

	sport, _ := strconv.Atoi(os.Args[1])
	port := uint(sport)

	initValues()

	http.Handle("/", websocket.Handler(calcHandler))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
