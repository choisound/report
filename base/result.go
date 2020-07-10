package base

const (
	//SUCCESS 成功
	SUCCESS = 20000
	//ERROR 错误
	ERROR = 20001
)

// Result 返回结果
type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
