package main

import (
	"reportapp/controller"
)

func main() {
	// build("E:\\GoProject\\dsapp\\main.go")
	// id, err := dao.Insert("insert into conference_info(acronym, name, location, submission_deadline, start_date, detail)values(?,?,?,?,?,?)", "1", "1", "1", "1", "1", "1")
	// fmt.Printf("%+v %+v\n", id, err)
	// thread.TeTest()
	// thread.MainThread()
	// var ss []base.TagInfo
	// dao.Query(&ss, "select * from tag_info")
	// fmt.Printf("%+v\n", ss)
	// da, rdr := dao.QueryConferenceTagGroup()
	// fmt.Printf("%+v %+v\n", da, rdr)
	// tc, er := dao.QueryTag(553)
	// fmt.Printf("%+v %+v\n", tc, er)
	// tcd, err := dao.QueryTagConference()
	// fmt.Printf("%+v %+v\n", tcd, err)

	// var up []controller.TagCounter
	// err := dao.Query(&up, "select tag_info.id, tag_info.name, count(*) as cnum from tag_info left join conference_tag_info on tag_info.id = conference_tag_info.tag_id where tag_info.id < 27 GROUP BY tag_info.id ")
	// if err == nil {
	// 	fmt.Printf("%+v\n", up)
	// } else {
	// 	fmt.Printf("err := %+v\n", err)
	// }
	controller.StartController()
	// rows, err1 := dao.Exec("update user_profile set chinese_name = ? where id = ?", "方文崇 upd", 1)
	// fmt.Printf("%+v %+v\n", rows, err1)

}
