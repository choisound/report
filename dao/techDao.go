package dao

// Tech 技术类
type Tech struct {
	ZfbAccount string `column:"zfbaccount"`
	Name       string `column:"name"`
	Phone      string `column:"phone"`
	Password   string `column:"password"`
	Area       string `column:"area"`
	QQ         string `column:"qq"`
}

// TechDetail 技术详情
type TechDetail struct {
	ZfbAccount string `column:"zfbaccount"`
	Name       string `column:"name"`
	Phone      string `column:"phone"`
	Password   string `column:"password"`
	Area       string `column:"area"`
	QQ         string `column:"qq"`
	Count      int    `column:"num"`
}

// TechRes 技术类请求结果
type TechRes struct {
	TechDetailList []TechDetail
	Num            int64
}

// GetTechInfo 获取技术数据
func GetTechInfo(page int64, limit int64) (*TechRes, error) {
	sql := "select count(1) from tech"
	var count int64
	err := Query(&count, sql)
	if err != nil {
		return nil, err
	}
	sql = "select tech.*,count(1) as num from tech left join order_info on tech.id = order_info.t_id GROUP BY tech.id order by num limit ?, ?"
	techDetail := []TechDetail{}
	err = Query(&techDetail, sql, (page-1)*limit, limit)
	techRes := TechRes{
		TechDetailList: techDetail,
		Num:            count,
	}
	return &techRes, err
}
