package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"reportapp/base"
	"runtime"

	"github.com/codegangsta/negroni"
)

func resigerController() (map[string]func(http.ResponseWriter, *http.Request), map[string]func(http.ResponseWriter, *http.Request)) {
	handlerMap := make(map[string]func(http.ResponseWriter, *http.Request))
	handlerMap["/api/resgister"] = Resgister
	handlerMap["/api/login"] = Login
	protectHandleMap := make(map[string]func(http.ResponseWriter, *http.Request))
	protectHandleMap["/api/uploadzfbOrder"] = AddZfbOrder
	protectHandleMap["/api/getOrderList"] = GetOrderByCondition
	protectHandleMap["/api/addReportData"] = AddReportData
	protectHandleMap["/api/getReportData"] = GetReportData
	protectHandleMap["/api/reportOrder"] = ReportOrder
	protectHandleMap["/api/getTechData"] = GetTechData
	return handlerMap, protectHandleMap
}

// StartController 开始监听请求
func StartController() {
	handlerMap, protectHandleMap := resigerController()
	for k, v := range handlerMap {
		fmt.Printf("%+v : %+v\n", k, runtime.FuncForPC(reflect.ValueOf(v).Pointer()).Name())
		http.HandleFunc(k, v)
	}
	for k, v := range protectHandleMap {
		fmt.Printf("%+v : %+v\n", k, runtime.FuncForPC(reflect.ValueOf(v).Pointer()).Name())
		http.Handle(k, negroni.New(
			negroni.HandlerFunc(ValidateTokenMiddleware),
			negroni.Wrap(http.HandlerFunc(v)),
		))
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) // 启动静态文件服务
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("ListenAndServe:	", err)
	}
}

// sendJSONMessage 发送数据
func sendJSONMessage(w http.ResponseWriter, res base.Result) {
	w.Header().Set("Access-Control-Allow-Origin", "*")                                //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Access-Token,token") //header的类型
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Expose-Headers", "*")
	w.WriteHeader(http.StatusOK) //返回数据格式是json
	rj, _ := json.Marshal(res)
	io.WriteString(w, string(rj))
}
