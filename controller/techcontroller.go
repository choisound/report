package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reportapp/base"
	"reportapp/dao"
)

// TechRequestCondition 筛选条件
type TechRequestCondition struct {
	Page int64 `json:"page"`
	Size int64 `json:"size"`
}

// GetTechData 获取技术信息
func GetTechData(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("method:", r.Method)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// fmt.Printf("read body err, %v\n", err)
		sendJSONMessage(w, base.Result{
			Code: 20001,
			Msg:  "param error" + err.Error(),
			Data: nil,
		})
		return
	}
	var requestCondition TechRequestCondition
	if len(body) > 0 {
		if err = json.Unmarshal(body, &requestCondition); err != nil {
			// fmt.Printf("Unmarshal err, %+v\n", err)
			sendJSONMessage(w, base.Result{
				Code: 20001,
				Msg:  "param error" + err.Error(),
				Data: nil,
			})
			return
		}
	}
	fmt.Printf("json:%+v\n", string(body))
	fmt.Printf("requestCondition:%+v\n", requestCondition)
	var res *dao.TechRes
	res, err = dao.GetTechInfo(requestCondition.Page, requestCondition.Size)
	if err != nil {
		sendJSONMessage(w, base.Result{
			Code: 20001,
			Msg:  "获取列表失败" + err.Error(),
			Data: nil,
		})
		return
	}
	sendJSONMessage(w, base.Result{
		Code: 20000,
		Msg:  "获取成功",
		Data: *res,
	})
}
