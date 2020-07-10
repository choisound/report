package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"reportapp/base"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	//SecretKey 秘钥
	SecretKey = "testSecreKeyError"
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//UserClaims 用户信息
type UserClaims struct {
	User           *UserData          `json:"userInfo"`
	StandardClaims jwt.StandardClaims `json:"standardClaims"`
}

//Valid 实现 `type Claims interface` 的 `Valid() error` 方法,自定义校验内容
func (c UserClaims) Valid() (err error) {
	if c.StandardClaims.VerifyExpiresAt(time.Now().Unix(), true) == false {
		return errors.New("token is expired")
	}
	if c.User.ID < 1 {
		return errors.New("invalid user in jwt")
	}
	return
}

//GenerateToken 产生token
func GenerateToken(userData *UserData) (string, error) {
	// expire := time.Now().Add(expireDuration)
	// 将 uid，用户角色， 过期时间作为数据写入 token 中
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)
	stdClaims := jwt.StandardClaims{
		ExpiresAt: expireTime.Unix(),
		IssuedAt:  time.Now().Unix(),
		Id:        fmt.Sprintf("%d", userData.ID),
		// Issuer:    AppIss,
	}
	userClaims := UserClaims{
		User:           userData,
		StandardClaims: stdClaims,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	// SecretKey 用于对用户数据进行签名，不能暴露
	return token.SignedString([]byte(SecretKey))
}

//getToken 通过request 解析出token等信息
func getToken(r *http.Request, userClaims *UserClaims) error { //由request获取token
	// t := T{} // t是已经实现extract接口的对象，对request进行处理得到tokenString并生成为解密的token
	// request.ParseFromRequest的第三个参数是一个keyFunc，具体的直接看源代码
	// 该keyFunc参数需要接受一个“未解密的token”，并返回Secretkey的字节和错误信息
	// keyFunc被调用并传入未解密的token参数，返回解密好的token和可能出现的错误
	// 若解密是正确的，那么返回的token.valid = true
	keys := []string{
		"token",
		"Token",
		"Authorization",
	}
	var tokenString string = ""
	for _, key := range keys {
		tokenString = r.Header.Get(key)
		if len(tokenString) > 0 {
			break
		}
	}
	if tokenString == "" {
		return errors.New("no token is found in Authorization Bearer")
	}
	_, err := jwt.ParseWithClaims(tokenString, userClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SecretKey), nil
	})
	return err
}

// validatePermission 验证用户权限
func validatePermission(userClaims *UserClaims, path string) bool {
	for _, key := range userClaims.User.Permission {
		fmt.Printf("key:%+v path %+v\n", key, path)
		matched, _ := regexp.MatchString(key, path)
		if matched {
			return true
		}
	}
	return false
}

//ValidateTokenMiddleware 校验token
func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// fmt.Printf("header:%+v param:%+v\n", r.Header, r.Form)
	userClaims := UserClaims{}
	err := getToken(r, &userClaims)
	// fmt.Printf("userClaims:%v UserData:%+v err:%+v\n", userClaims.StandardClaims, userClaims.User, err)
	if err != nil && userClaims.Valid() != nil {
		sendJSONMessage(w, base.Result{
			Code: base.ERROR,
			Msg:  "token 解析失效",
			Data: err,
		})
		return
	}
	if !validatePermission(&userClaims, r.URL.Path) {
		sendJSONMessage(w, base.Result{
			Code: base.ERROR,
			Msg:  "权限不足",
			Data: nil,
		})
		return
	}
	next(w, r)
}
