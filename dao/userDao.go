package dao

// User 用户结构体
type User struct {
	ID       int64  `column:"id"`
	Username string `column:"username"`
	Password string `column:"password"`
	Verify   int    `column:"verify"`
	Type     int    `column:"type"`
}

// Permission 权限
type Permission struct {
	ID          int64  `column:"id"`
	Name        string `column:"name"`
	Description string `column:"description"`
	URL         string `column:"url"`
}

// AddUser 添加用户
func AddUser(username string, password string, userType int) (int64, error) {
	return Insert("insert into user(username, password, type) values(?,?,?)", username, password, userType)
}

//GetUser 通过账号密码获取用户
func GetUser(username string, password string) (User, error) {
	var user User
	err := Query(&user, "select * from user where username = ? and password = ? and verify = 1", username, password)
	return user, err
}

//VerifyUser 验证用户
func VerifyUser(username string, password string) error {
	_, err := Exec("update user set verify = 1 where username = ? and password = ?", username, password)
	return err
}

// RenewPassword 更新密码
func RenewPassword(id int64, username string, password string) error {
	_, err := Exec("update user set password = ? where username = ? and id = ?", password, username, id)
	return err
}

// GetPermission 获取权限
func GetPermission(uid int64, p interface{}) error {
	err := Query(p, "select permission.* from (select p_id from user_role left join role_permission on user_role.r_id=role_permission.r_id where user_role.u_id= ?) as rl left join permission on rl.p_id=permission.id", uid)
	return err
}
