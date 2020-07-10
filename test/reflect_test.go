package test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"testing"
	"time"
)

// TestTet 测试
func TestTet(t *testing.T) {
	const DF = "Jan 2, 2006"

	dstr := "internet of things"
	Td, err := time.Parse(DF, dstr)
	fmt.Printf("%+v %+v\n", Td, err)
	dstr = "Aug 24, 2020"

	Td, err = time.Parse(DF, dstr)
	fmt.Printf("%+v %+v %s\n", Td, err, Td.Format("2006-01-01 15:04:05"))
	// conf := base.ConferenceInfo{
	// 	Acronym:            "SWEHCS-2021",
	// 	Name:               "Tutorial and workshop on uncertainty in machine learning",
	// 	Location:           "Ghent, Belgium",
	// 	SubmissionDeadline: "Jun 11, 2020",
	// 	StartData:          ":Sep 14, 2020",
	// 	Detail:             "",
	// 	URL:                "",
	// }
	// id, err1 := dao.Insert("insert into conference_info(acronym, name, location, submission_deadline, start_date, detail)values(?,?,?,?,?,?)",
	// 	conf.Acronym, conf.Name, conf.Location, conf.SubmissionDeadline, conf.StartData, conf.Detail)
	// fmt.Printf("%+v\n%+v\n", id, err1)
	// var up []dao.UserProfile
	// err := dao.Query(&up, "select id, acc_name, chinese_name, scholar_field, introduction from user_profile where ", 1, 5)
	// rows, err1 := dao.Exec("update user_profile set chinese_name = ? where id = ?", "方文崇 updata", 1)
	// fmt.Printf("%+v %+v\n", rows, err1)mmmn
	// if err != nil {
	// 	fmt.Printf("%+v\n", up)
	// }
	// for i := 0; i < 10; i++ {
	// 	go testFor(i)
	// }
	// tr := reflect.ValueOf(up)
	// tet(tr)
	// time.Sleep(time.Second * time.Duration(10))
	tet()
}

func testFor(tr int) {
	fmt.Println(tr)
}

func Byte2Int(data []byte) int64 {
	ret, _ := strconv.ParseInt(string(data), 10, 64)

	return ret
}
func tet() {
	fmt.Println("dddd")

	// fmt.Printf("%+v\n", tt.Type())
	// var ecd thread.EasychairDatableCrawl
	// data, err1 := ecd.CrawlData("https://www.easychair.org/cfp/")
	// fmt.Printf("%+v %+v\n", string(data), err1)
	// dao.SaddItem("http://localhost:8080", "mainContentSet")
	// res, err := dao.SpopItem("mainContentSet")
	// fmt.Printf("%+v %+v\n", res, err)
	// thread.MainThread()
	// Conf := base.ConferenceInfo{ID: 0, Acronym: "IEEE EDGE 2020", Name: "2020 IEEE International Conference on Edge Computing", Location: "Beijing, China", SubmissionDeadline: "Oct 19, 2020", StartData: "edge computing", Detail: "", URL: "https://www.easychair.org/cfp/ieeeedge2020"}
	// cf, _ := json.Marshal(Conf)
	byVal := []byte{53, 53, 51}
	// byVal := []byte{6: 1, 7: 255}
	fmt.Printf("%+v\n", Byte2Int(byVal))
	var ints int64
	buf := bytes.NewBuffer(byVal)
	binary.Read(buf, binary.LittleEndian, &ints)
	fmt.Printf("%+v\n", ints)
	// str, err := dao.SpopItem("easychair_redis_set")
	// // newSlice := reflect.MakeSlice(tt.Type(), 100, 100)
	// fmt.Printf("%+v %+v\n", str, err)
}
