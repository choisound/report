package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reportapp/base"
	"reportapp/dao"
	"strconv"
	"time"
)

//Order 订单
type Order struct {
	ID         int64     `json:"id"`
	TID        int64     `json:"tId"`
	Money      float64   `json:"money"`
	OrderID    string    `json:"orderId"`
	QQ         string    `json:"qq"`
	KID        int64     `json:"kId"`
	Report     int16     `json:"report"`
	ZfbID      string    `json:"zfbId"`
	ReportTime time.Time `json:"reportTime"`
}

// RequestCondition 筛选条件
type RequestCondition struct {
	TechID          int64     `json:"techId"`
	OrderID         string    `json:"orderId"`
	Report          int16     `json:"report"`
	QQ              string    `json:"qq"`
	Money           float64   `json:"money"`
	ReportStartTime time.Time `json:"reportStartTime"`
	ReportEndTime   time.Time `json:"reportEndTime"`
	ZfbID           string    `json:"zfbId"`
	KID             int64     `json:"kId"`
	Page            int64     `json:"page"`
	Size            int64     `json:"size"`
}

//ZfbOrder 支付宝订单
type ZfbOrder struct {
	OrderID string  `json:"orderId"`
	Money   float64 `json:"money"`
	ZfbID   string  `json:"zfbId"`
}

// OrderResult 订单结果
type OrderResult struct {
	Total     int64           `json:"total"`
	OrderList []dao.OrderData `json:"orderList"`
}

// ReportData 报账数据
type ReportData struct {
	Money    int    `json:"money"`
	OrderID  string `json:"orderId"`
	QQ       string `json:"qq"`
	TechName string `json:"techName"`
	TechZfb  string `json:"techZfb"`
	KID      int64  `json:"kId"`
}

// AddZfbOrder 添加订单
func AddZfbOrder(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("read body err, %v\n", err)
		sendJSONMessage(w, base.Result{
			Code: 20001,
			Msg:  "param error" + err.Error(),
			Data: nil,
		})
		return
	}
	var a []ZfbOrder
	if err = json.Unmarshal(body, &a); err != nil {
		// fmt.Printf("Unmarshal err, %+v\n", err)
		sendJSONMessage(w, base.Result{
			Code: 20001,
			Msg:  "param error" + err.Error(),
			Data: nil,
		})
		return
	}
	// fmt.Printf("%+v", a)
	var params []interface{}
	sql := "insert IGNORE  into order_info(order_id,zfb_id,money) values "
	for index, v := range a {
		params = append(params, v.OrderID, v.ZfbID, v.Money)
		if index != 0 {
			sql += ","
		}
		sql += "(?,?,?)"
	}
	_, err = dao.InsertZfbOrderList(sql, params)
	if err != nil {
		fmt.Printf("%+v\n", err)
		sendJSONMessage(w, base.Result{
			Code: 20001,
			Msg:  "数据库操作出错",
			Data: nil,
		})
		return
	}
	sendJSONMessage(w, base.Result{
		Code: 20000,
		Msg:  "插入成功",
		Data: nil,
	})
}

// ReportDataRes 报账数据
type ReportDataRes struct {
	Name        string   `json:"Name"`
	Zfb         string   `json:"Zfb"`
	Money       float64  `json:"Money"`
	Tag         string   `json:"Tag"`
	Type        string   `json:"type"`
	OrderIDList []string `json:"orderIdList"`
}

// GetReportData 获取报账信息
func GetReportData(w http.ResponseWriter, r *http.Request) {
	reportData, err := dao.GetReportInfo()
	mapReport := make(map[string][]dao.ReportData)
	if err != nil {
		sendJSONMessage(w, base.Result{
			Code: 20001,
			Msg:  "数据库操作出错",
			Data: nil,
		})
		return
	}
	for _, item := range reportData {
		list := mapReport[item.Zfb]
		mapReport[item.Zfb] = append(list, item)
	}
	res := []ReportDataRes{}
	for _, v := range mapReport {
		var allMoney float64 = 0
		var tag string = "("
		var zfb string = ""
		var name string = ""
		var typet string = "kefu"
		var orderIDList []string
		for _, val := range v {
			if val.Type == "tech" {
				typet = "tech"
				allMoney += val.Money * 0.8
			} else {
				allMoney += val.Money * 0.1
			}
			orderIDList = append(orderIDList, val.OrderID)
			tag += "群号：" + val.QQ + " 金额：" + strconv.FormatFloat(val.Money, 'f', 0, 64) + ","
			name = val.Name
			zfb = val.Zfb
		}
		tag += ")"
		res = append(res, ReportDataRes{
			Name:        name,
			Zfb:         zfb,
			Tag:         tag,
			Money:       allMoney,
			Type:        typet,
			OrderIDList: orderIDList,
		})
	}
	sendJSONMessage(w, base.Result{
		Code: 20000,
		Msg:  "获取成功",
		Data: res,
	})
}

// ReportOrder 报账
func ReportOrder(w http.ResponseWriter, r *http.Request) {
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
	var orderIDList []ReportDataRes
	if len(body) > 0 {
		if err = json.Unmarshal(body, &orderIDList); err != nil {
			// fmt.Printf("Unmarshal err, %+v\n", err)
			sendJSONMessage(w, base.Result{
				Code: 20001,
				Msg:  "param error" + err.Error(),
				Data: nil,
			})
			return
		}
	}
	fmt.Printf("orderIDList:%+v\n", orderIDList)
	sqlArgs := []dao.SQLArgs{}
	for _, val := range orderIDList {
		params := []interface{}{}
		sql := "update order_info set report = ?|report where order_id = ?"
		val1 := 2
		if val.Type == "tech" {
			val1 = 1
		}
		for i := 0; i < len(val.OrderIDList)-1; i++ {
			sql += " or order_id = ?"
		}
		params = append(params, val1)
		for _, v1 := range val.OrderIDList {
			params = append(params, v1)
		}
		sqlArgs = append(sqlArgs, dao.SQLArgs{
			SQL:    sql,
			Params: params,
		})
	}
	err = dao.UpdateReportInfo(sqlArgs)
	if err != nil {
		sendJSONMessage(w, base.Result{
			Code: 20001,
			Msg:  "param error" + err.Error(),
			Data: nil,
		})
		return
	}
	sendJSONMessage(w, base.Result{
		Code: 20000,
		Msg:  "报账成功",
		Data: nil,
	})
}

// GetOrderByCondition 获取订单
func GetOrderByCondition(w http.ResponseWriter, r *http.Request) {
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
	var requestCondition RequestCondition
	//表示获取全部
	requestCondition.Report = -1
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
	var orderResult OrderResult
	flag := 1
	var params []interface{}
	sql := ""
	dataSQL := " select order_info.*,tech.zfbaccount as techZfb,tech.name as techName from order_info left join tech on tech.id = order_info.t_id "
	countSQL := "select count(*) from order_info left join tech on tech.id = order_info.t_id "
	if requestCondition.TechID > 0 {
		flag = 0
		sql += " where order_info.t_id = ? "
		params = append(params, requestCondition.TechID)
	}
	if len(requestCondition.OrderID) > 0 {
		if flag == 1 {
			sql += " where "
		} else {
			sql += " and "
		}
		sql += " order_info.order_id = ?"
		params = append(params, requestCondition.OrderID)

		flag = 0
	}
	if requestCondition.Report >= 0 {
		if flag == 1 {
			sql += " where "
		} else {
			sql += " and "
		}
		flag = 0
		sql += " order_info.report = ?"
		params = append(params, requestCondition.Report)
	}
	if len(requestCondition.QQ) > 0 {
		if flag == 1 {
			sql += " where "
		} else {
			sql += " and "
		}
		flag = 0
		sql += " order_info.qq = ?"
		params = append(params, requestCondition.QQ)
	}
	if requestCondition.Money > 0 {
		if flag == 1 {
			sql += " where "
		} else {
			sql += " and "
		}
		flag = 0
		sql += " order_info.money = ?"
		params = append(params, requestCondition.Money)
	}
	if !requestCondition.ReportStartTime.IsZero() && !requestCondition.ReportEndTime.IsZero() && requestCondition.ReportStartTime.Before(requestCondition.ReportEndTime) {
		if flag == 1 {
			sql += " where "
		} else {
			sql += " and "
		}
		flag = 0
		sql += " order_info.report_time between (?,?)"
		params = append(params, requestCondition.ReportStartTime, requestCondition.ReportEndTime)
	}

	if len(requestCondition.ZfbID) > 0 {
		if flag == 1 {
			sql += " where "
		} else {
			sql += " and "
		}
		flag = 0
		sql += " order_info.zfb_id = ?"
		params = append(params, requestCondition.ZfbID)
	}

	if requestCondition.KID > 0 {
		if flag == 1 {
			sql += " where "
		} else {
			sql += " and "
		}
		flag = 0
		sql += " order_info.k_id = ?"
		params = append(params, requestCondition.KID)
	}
	fmt.Printf("sql:%+v param:%v\n", countSQL+sql, params)
	var countNum int64 = 0
	countNum, err = dao.GetOrderCount(countSQL+sql, params)
	if err != nil {
		sendJSONMessage(w, base.Result{
			Code: 20001,
			Msg:  "获取列表失败" + err.Error(),
			Data: nil,
		})
		return
	}
	orderResult.Total = countNum
	sql += " limit ?,? "
	if requestCondition.Page > 0 && requestCondition.Size > 0 {
		params = append(params, (requestCondition.Page-1)*requestCondition.Size, requestCondition.Size)
	} else {
		params = append(params, 0, 10)
	}
	orders := make([]dao.OrderData, 0)
	fmt.Printf("sql:%+v params:%+v\n", sql, params)
	err = dao.GetOrderList(&orders, dataSQL+sql, params)
	// fmt.Printf("%+v\n", orders)
	if err != nil {
		fmt.Printf("err:%+v\n", err)
		sendJSONMessage(w, base.Result{
			Code: 20001,
			Msg:  "获取列表失败" + err.Error(),
			Data: nil,
		})
		return
	}
	orderResult.OrderList = orders
	sendJSONMessage(w, base.Result{
		Code: 20000,
		Msg:  "获取成功",
		Data: orderResult,
	})
}

// AddReportData 添加报账信息
func AddReportData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("read body err, %v\n", err)
		sendJSONMessage(w, base.Result{
			Code: 20001,
			Msg:  "param error",
			Data: nil,
		})
		return
	}
	NoExistData := make([]ReportData, 0)
	var a []ReportData
	if err = json.Unmarshal(body, &a); err != nil {
		fmt.Printf("Unmarshal err, %+v\n", err)
		sendJSONMessage(w, base.Result{
			Code: 20001,
			Msg:  "param error" + err.Error(),
			Data: nil,
		})
		return
	}
	ZfbMap := make(map[string]int64)

	var TID int64
	for _, item := range a {
		TID = 0
		if _, ok := ZfbMap[item.TechZfb]; !ok {
			err = dao.Query(&TID, "select count(*) from tech where zfbaccount = ?", item.TechZfb)
			if err != nil {
				sendJSONMessage(w, base.Result{
					Code: 20001,
					Msg:  "select Error" + err.Error(),
					Data: nil,
				})
				return
			}
			if TID == 0 {
				TID, err = dao.Insert("insert into tech(name, zfbaccount) values(?,?) ", item.TechName, item.TechZfb)
				if err != nil {
					sendJSONMessage(w, base.Result{
						Code: 20001,
						Msg:  "select Error" + err.Error(),
						Data: nil,
					})
				}

				return
			}
		}
		// fmt.Printf("%+v\n", TID)
		params := []interface{}{
			item.QQ,
			TID,
			item.KID,
			item.OrderID,
		}
		var num int64
		fmt.Printf("%+v %+v %+v %+v\n", item.QQ, TID, item.KID, item.OrderID)
		num, err = dao.UpdateOrderInfo("update order_info set qq = ?,t_id=?, k_id=? where order_id = ? and t_id = 0 ", params)
		if err != nil {
			sendJSONMessage(w, base.Result{
				Code: 20001,
				Msg:  "select Error" + err.Error(),
				Data: nil,
			})
			return
		}
		if num == 0 {
			NoExistData = append(NoExistData, item)
			// fmt.Printf("%+v is no exist!", item.OrderID)
		}
	}
	sendJSONMessage(w, base.Result{
		Code: 20000,
		Msg:  "插入成功",
		Data: NoExistData,
	})
}
