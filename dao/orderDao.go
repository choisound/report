package dao

import "time"

//Order 订单表
type Order struct {
	ID         int64     `column:"id"`
	TID        int64     `column:"t_id"`
	Money      float64   `column:"money"`
	OrderID    string    `column:"order_id"`
	QQ         string    `column:"qq"`
	KID        int64     `column:"k_id"`
	Report     int16     `column:"report"`
	ZfbID      string    `column:"zfb_id"`
	ReportTime time.Time `column:"report_time"`
}

//OrderData 订单表
type OrderData struct {
	ID         int64     `column:"id"`
	TID        int64     `column:"t_id"`
	Money      float64   `column:"money"`
	OrderID    string    `column:"order_id"`
	QQ         string    `column:"qq"`
	KID        int64     `column:"k_id"`
	Report     int16     `column:"report"`
	ZfbID      string    `column:"zfb_id"`
	ReportTime time.Time `column:"report_time"`
	TechName   string    `column:"techName"`
	TechZfb    string    `column:"techZfb"`
}

// ReportData 报账数据
type ReportData struct {
	//
	QQ      string  `column:"qq"`
	OrderID string  `column:"order_id"`
	Money   float64 `column:"money"`
	Report  int16   `column:"report"`
	Name    string  `column:"name"`
	Zfb     string  `column:"zfb"`
	Type    string  `column:"type"`
}

// GetOrderList 获取订单列表
func GetOrderList(orders interface{}, sql string, param []interface{}) error {
	return Query(orders, sql, param...)
}

// GetOrderCount 获取订单数量
func GetOrderCount(sql string, param []interface{}) (int64, error) {
	var num int64
	err := Query(&num, sql, param...)
	if err != nil {
		return -1, err
	}
	return num, nil
}

// UpdateOrderInfo 更新订单信息
func UpdateOrderInfo(sql string, param []interface{}) (int64, error) {
	num, err := Exec(sql, param...)
	if err != nil {
		return -1, err
	}
	return num, nil
}

// InsertZfbOrderList 插入支付宝订单列表
func InsertZfbOrderList(sql string, params []interface{}) (int64, error) {
	return Exec(sql, params...)
}

// GetReportInfo 获取报账数据
func GetReportInfo() ([]ReportData, error) {
	sql := "select 'kefu' as type,order_info.order_id, order_info.money, order_info.qq, order_info.report, kefu.name as name, kefu.phone as zfb from order_info,kefu where order_info.k_id = kefu.id and order_info.report & 2 = 0"
	reportData := []ReportData{}
	reportData1 := []ReportData{}
	err := Query(&reportData, sql)
	if err != nil {
		return nil, err
	}
	sql = "select 'tech' as type,order_info.order_id, order_info.qq, order_info.money, order_info.report,tech.zfbaccount as zfb, tech.name as name from order_info,tech where order_info.t_id = tech.id and order_info.report & 1 = 0"
	err = Query(&reportData1, sql)
	if err != nil {
		return nil, err
	}
	reportData = append(reportData, reportData1...)
	return reportData, err
}

// UpdateReportInfo 报账
func UpdateReportInfo(sqlArgs []SQLArgs) error {
	return TxExecute(sqlArgs)
}
