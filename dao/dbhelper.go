package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
)

// SQLArgs sql封装类
type SQLArgs struct {
	SQL    string
	Params []interface{}
}

//DbConfig 数据库配置
type DbConfig struct {
	DBConfig struct {
		Name        string `toml:"name"`
		Host        string `toml:"host"`
		Port        int    `toml:"port"`
		User        string `toml:"user"`
		Database    string `toml:"database"`
		Password    string `toml:"password"`
		Network     string `toml:"network"`
		MaxTime     int    `toml:"maxTime"`
		MaxOpenConn int    `toml:"maxOpenConn"`
		MaxIdleConn int    `toml:"maxIdleConn"`
	}
}

// DB mysql链接
var mDB *sql.DB

// Exec 执行语句
func Exec(sql string, param ...interface{}) (int64, error) {
	result, err := mDB.Exec(sql, param...)
	if err != nil {
		return 0, err
	}
	rowsaffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsaffected, nil
}

// Insert 执行语句 返回id
func Insert(sql string, param ...interface{}) (int64, error) {
	stmt, err := mDB.Prepare(sql)
	if err != nil {
		fmt.Printf("sql:%v err1 %v \n\n", sql, err)
		return -1, err
	}

	res, err1 := stmt.Exec(param...)
	if err1 != nil {
		fmt.Printf("sql:%v err%v \n\n", sql, err1)
		return -1, err
	}
	return res.LastInsertId()
}

//Query 查询
func Query(result interface{}, sql string, param ...interface{}) error {
	// 接收数据指针为空
	if result == nil {
		fmt.Println("传入结构体为空，查询出错")
		return errors.New("传入结构体为空，查询出错")
	}
	// 执行语句
	rows, err := mDB.Query(sql, param...)
	if err != nil {
		return err
	}
	// 返回
	defer rows.Close()
	var index int = 0
	var indexValue reflect.Value
	var valueSlice reflect.Value
	// 遍历结果集进行初始化
	for rows.Next() {
		columnValueMap := make(map[string]interface{})
		sqlColumn, err := rows.Columns()
		// 获取列数失败，返回错误
		if err != nil {
			return err
		}
		// fmt.Printf("%+v %+v\n", sql, sqlColumn)
		// 生成interface{}指针数据作为scan方法参数
		sqlValues := make([]([]byte), len(sqlColumn))
		sqlValuesPointer := make([]interface{}, len(sqlColumn))
		for i := 0; i < len(sqlColumn); i++ {
			sqlValuesPointer[i] = &sqlValues[i]
		}
		// 扫描数据
		err = rows.Scan(sqlValuesPointer...)
		// 扫描出错返回错误
		if err != nil {
			return err
		}
		// 将查询出来的一行数据放在map中
		for i := 0; i < len(sqlColumn); i++ {
			if sqlValuesPointer != nil {
				columnValueMap[sqlColumn[i]] = sqlValues[i]
				// fmt.Printf("name:%+v value:%+v\n", sqlColumn[i], sqlValues[i])
			}
		}
		if len(sqlColumn) == 0 {
			return nil
		}
		resultType := reflect.TypeOf(result).Elem().Kind()
		// 当查询结果需要返回单个数字
		if resultType == reflect.Uint || resultType == reflect.Int || resultType == reflect.Int16 || resultType == reflect.Int32 || resultType == reflect.Int64 {
			singleValue := reflect.ValueOf(result).Elem()
			if len(columnValueMap) == 1 {
				for _, val := range columnValueMap {
					ret, err := strconv.ParseInt(string(val.([]byte)), 10, 64)
					if err != nil {
						return err
					}
					singleValue.SetInt(ret)
				}
				return nil
			}
			return errors.New("查询数量sql语句返回长度超过1一个，Sql语句为" + sql + ";请检查你的sql语句")
			// 当查询结果需要返回一个结构体
		} else if resultType != reflect.Array && resultType != reflect.Slice {
			structValue := reflect.ValueOf(result).Elem()
			structType := reflect.TypeOf(result).Elem()
			fillValueInStruct(columnValueMap, structValue, structType)
			// reflect.ValueOf(result).Set
			return nil
		}
		// fmt.Printf("sqlValue:%+v %+v\n", sqlValues, result)
		// 当查询结果需要返回结构体数组
		if index == 0 {
			//如果是首个元素需要构造一个Slice和获取结构体的类型
			tt := reflect.ValueOf(result)
			valueSlice = reflect.MakeSlice(tt.Elem().Type(), 1, 1)
			indexValue = valueSlice.Index(0)
			//给元素填充元素
			err = fillValueInStruct(columnValueMap, valueSlice.Index(0), indexValue.Type())
			if err != nil {
				return err
			}
		} else {
			// 反射新建一个对象
			newVal := reflect.New(indexValue.Type()).Elem()
			err = fillValueInStruct(columnValueMap, newVal, indexValue.Type())
			if err != nil {
				return err
			}
			valueSlice = reflect.Append(valueSlice, newVal)
		}
		index++
	}
	if index == 0 {
		// fmt.Printf("query nil")
		return nil
	}
	//填充数据
	val := reflect.ValueOf(result)
	val = val.Elem()
	val.Set(valueSlice)
	return nil
}

// fillValueInStruct 给结构体填充数据
func fillValueInStruct(columnValueMap map[string]interface{}, structValue reflect.Value, structType reflect.Type) (err error) {
	typeOfType := structValue.Type()
	//获取遍历结构体成员域
	for i := 0; i < structValue.NumField(); i++ {
		var columnName string
		//如果结构体有标签 以标签为主 否则以结构体名称为主
		if structType.Field(i).Tag.Get("column") != "" {
			columnName = structType.Field(i).Tag.Get("column")
		} else {
			columnName = typeOfType.Field(i).Name
		}
		// fmt.Printf("columnName:  %+v\n", columnName)
		//返回数据 不能放在类型转换之后，不然保存的现场将会是上一次执行成功的现场
		defer func() {
			if r := recover(); r != nil {
				switch x := r.(type) {
				case string:
					err = errors.New(x + " column ：" + columnName)
				case error:
					err = errors.New(fmt.Sprintf("%s", x) + " column:" + columnName)
				default:
					err = errors.New("Unknow panic column:" + columnName)
				}
				fmt.Printf("error %+v\n", err)
			}
		}()
		//获取具体某个成员
		field := structValue.Field(i)
		// fmt.Printf("%+v\n", columnValueMap)
		//判断查询结果中是否包含该成员
		if mapValue, ok := columnValueMap[columnName]; ok {
			// fmt.Printf("columnName:%+v  mapValue:%+v\n", columnName, mapValue)
			//包含的话判断该成员的类型，进行数据填充
			byVal := mapValue.([]byte)
			kind := field.Type().Kind()
			switch kind {
			case reflect.Bool:
				b, err := strconv.ParseBool(string(byVal))
				if err != nil {
					return err
				}
				field.SetBool(b)
			case reflect.Int:
				ret, err := strconv.ParseInt(string(byVal), 10, 0)
				if err != nil {
					return err
				}
				field.SetInt(ret)
			case reflect.Int8:
				ret, err := strconv.ParseInt(string(byVal), 10, 8)
				if err != nil {
					return err
				}
				field.SetInt(ret)
			case reflect.Int16:
				ret, err := strconv.ParseInt(string(byVal), 10, 16)
				if err != nil {
					return err
				}
				field.SetInt(ret)
			case reflect.Int32:
				ret, err := strconv.ParseInt(string(byVal), 10, 32)
				if err != nil {
					return err
				}
				field.SetInt(ret)
			case reflect.Int64:
				ret, err := strconv.ParseInt(string(byVal), 10, 64)
				if err != nil {
					return err
				}
				field.SetInt(ret)
			case reflect.Uint:
				ret, err := strconv.ParseUint(string(byVal), 10, 0)
				if err != nil {
					return err
				}
				field.SetUint(ret)
			case reflect.Uint8:
				ret, err := strconv.ParseUint(string(byVal), 10, 8)
				if err != nil {
					return err
				}
				field.SetUint(ret)
			case reflect.Uint16:
				ret, err := strconv.ParseUint(string(byVal), 10, 16)
				if err != nil {
					return err
				}
				field.SetUint(ret)
			case reflect.Uint32:
				ret, err := strconv.ParseUint(string(byVal), 10, 32)
				if err != nil {
					return err
				}
				field.SetUint(ret)
			case reflect.Uint64:
				ret, err := strconv.ParseUint(string(byVal), 10, 64)
				if err != nil {
					return err
				}
				field.SetUint(ret)
			case reflect.Uintptr:
				ret, err := strconv.ParseUint(string(byVal), 10, 64)
				if err != nil {
					return err
				}
				field.SetUint(ret)
			case reflect.Float32:
				ret, err := strconv.ParseFloat(string(byVal), 32)
				if err != nil {
					return err
				}
				field.SetFloat(ret)
			case reflect.Float64:
				ret, err := strconv.ParseFloat(string(byVal), 32)
				if err != nil {
					return err
				}
				field.SetFloat(ret)
			case reflect.String:
				field.SetString(string(byVal))
			default:
				break
			}
		}
		// fmt.Printf("%d. %s %s = %v %+v %+v\n", i, typeOfType.Field(i).Name, field.Type(), field.Interface(), structType.Field(i).Tag.Get("column"), field)
	}
	return err
}

// GetDBConfig 获取数据库配置
func GetDBConfig() (*DbConfig, error) {
	var dbConfig DbConfig
	if _, err := toml.DecodeFile("config/db.toml", &dbConfig); err != nil {
		return nil, err
	}
	return &dbConfig, nil
}

func init() {
	mainconfig, err := GetDBConfig()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	fmt.Printf("%+v\n", mainconfig)
	config := mainconfig.DBConfig
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8&parseTime=true", config.User, config.Password, config.Network, config.Host, config.Port, config.Database)
	// fmt.Println(dsn)
	DB, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("Open mysql failed,err:%v\n", err)
		return
	}
	DB.SetConnMaxLifetime(time.Duration(config.MaxTime) * time.Second)
	DB.SetMaxOpenConns(config.MaxOpenConn)
	DB.SetMaxIdleConns(config.MaxIdleConn)
	mDB = DB
}

// TxExecute 批量执行语句 事务
func TxExecute(sqlArgs []SQLArgs) error {
	sqlTx, err := mDB.Begin()
	if err != nil {
		return err
	}
	for _, val := range sqlArgs {
		_, err := sqlTx.Exec(val.SQL, val.Params...)
		if err != nil {
			sqlTx.Rollback()
			return err
		}
	}
	return sqlTx.Commit()
}
