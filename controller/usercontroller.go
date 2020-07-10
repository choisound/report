package controller

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"reportapp/base"
	"reportapp/dao"
	"strconv"
)

//UserData 用户凭证
type UserData struct {
	ID         int64    `json:"id"`
	Username   string   `json:"username"`
	Permission []string `json:"permission"`
	Name       string   `json:"name"`
}

// Token 返回给前端
type Token struct {
	Token string `json:"token"`
}

// Resgister 用户注册
func Resgister(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	userType := r.FormValue("type")
	typeValue, err1 := strconv.Atoi(userType)
	if err1 != nil {
		sendJSONMessage(w, base.Result{
			Code: base.ERROR,
			Msg:  "注册用户失败！！！",
			Data: nil,
		})
		return
	}
	password = getMd5String2([]byte(password))
	id, err := dao.AddUser(username, password, typeValue)

	if err != nil {
		sendJSONMessage(w, base.Result{
			Code: base.ERROR,
			Msg:  "注册用户失败！！！",
			Data: nil,
		})
		return
	}
	sendJSONMessage(w, base.Result{
		Code: base.SUCCESS,
		Msg:  "注册成功",
		Data: id,
	})
}

// Login 用户登录
func Login(w http.ResponseWriter, r *http.Request) { //返回数据格式是json
	username := r.FormValue("username")
	password := r.FormValue("password")
	password = getMd5String2([]byte(password))
	user, err := dao.GetUser(username, password)
	user.Password = ""
	if err != nil && len(user.Username) != 0 {
		sendJSONMessage(w, base.Result{
			Code: base.ERROR,
			Msg:  "登录失败！！！",
			Data: nil,
		})
		return
	}
	permission := make([]dao.Permission, 0)
	err = dao.GetPermission(user.ID, &permission)
	if err != nil {
		if err != nil && len(user.Username) != 0 {
			sendJSONMessage(w, base.Result{
				Code: base.ERROR,
				Msg:  "获取权限失败！！！",
				Data: nil,
			})
			return
		}
	}
	permissionStr := make([]string, len(permission))
	// fmt.Printf("%+v\n", permission)
	for _, url := range permission {
		permissionStr = append(permissionStr, url.URL)
		// fmt.Printf("%+v\n", url.URL)
	}
	// fmt.Printf("%+v\n", permissionStr)
	userData := UserData{
		ID:         user.ID,
		Username:   user.Username,
		Name:       "test",
		Permission: permissionStr,
	}
	token, err := GenerateToken(&userData)
	if err != nil {
		sendJSONMessage(w, base.Result{
			Code: base.ERROR,
			Msg:  "登录失败！！！",
			Data: nil,
		})
		return
	}
	sendJSONMessage(w, base.Result{
		Code: base.SUCCESS,
		Msg:  "登录成功",
		Data: Token{Token: token},
	})
}

//获取md5
func getMd5String2(b []byte) string {
	return fmt.Sprintf("%x", md5.Sum(b))
}
