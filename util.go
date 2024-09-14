package goo

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func Md5(input string) string {
	hasher := md5.New()         // 创建一个md5 Hash对象
	hasher.Write([]byte(input)) // 将字符串写入hasher
	digest := hasher.Sum(nil)   // 计算哈希并获取16字节的结果

	// 将16字节的哈希转换为16进制字符串
	MD5String := fmt.Sprintf("%x", digest)

	return MD5String
}

func DurationToChinese(d time.Duration) string {
	days := d / (24 * time.Hour)
	hours := (d % (24 * time.Hour)) / time.Hour
	minutes := (d % time.Hour) / time.Minute
	seconds := (d % time.Minute) / time.Second

	if days == 0 && hours == 0 && minutes == 0 {
		return fmt.Sprintf("%d秒", seconds)
	}

	if days == 0 && hours == 0 {
		return fmt.Sprintf("%d分钟%d秒", minutes, seconds)
	}

	if days == 0 {
		return fmt.Sprintf("%d小时%d分钟%d秒", hours, minutes, seconds)
	}

	return fmt.Sprintf("%d天%d小时%d分钟%d秒", days, hours, minutes, seconds)
}

// 判断变量是否为0，只有数字类型才可能返回true
func IsNumZero(v any) bool {
	switch v := v.(type) {
	case int, int8, int16, int32, int64:
		return v == 0
	case uint, uint8, uint16, uint32, uint64, uintptr:
		return v == 0
	case float32, float64:
		return v == 0.0
	default:
		return false
	}
}

// isNumeric 使用 reflect 判断给定的 interface{} 类型是否为数字类型
func IsNumeric(data any) bool {
	value := reflect.ValueOf(data)

	if !value.IsValid() {
		return false
	}

	kind := value.Kind()

	// 判断是否为数字类型
	return kind >= reflect.Int && kind <= reflect.Float64
}

func IsInteger(data any) bool {
	value := reflect.ValueOf(data)

	if !value.IsValid() {
		return false
	}

	kind := value.Kind()

	// 判断是否为数字类型
	return kind >= reflect.Int && kind <= reflect.Int64
}

func IsFloat(data any) bool {
	value := reflect.ValueOf(data)

	if !value.IsValid() {
		return false
	}

	kind := value.Kind()

	// 判断是否为数字类型
	return kind == reflect.Float32 || kind == reflect.Float64
}

// 判断变量是否为空
func Empty(v any) bool {
	if v == nil {
		return true
	}

	val := reflect.ValueOf(v)

	if !val.IsValid() {
		return true
	}

	return val.IsZero()
}

// 判断map类型的key是否存在
func IsSet[C comparable, V any](m map[C]V, key C) bool {
	_, ok := m[key]

	return ok
}

// 判断map类型的key是否存在, 不存在时返回
func GetMapWsDef[C comparable, V any, DV any](m map[C]V, key C, def DV) (DV, bool) {
	v, ok := m[key]

	if !ok {
		return def, ok
	}

	return AnyConvert2T(v, def), ok
}

/*
断言 any 类型是否能转换为指定类型，如果是，返回断言后的结果，否则返回指定的值
*/
func AnyConvert2T[T any](v any, t T) T {
	vVal := reflect.ValueOf(v)
	tVal := reflect.ValueOf(t)

	if !tVal.IsValid() || !vVal.IsValid() {
		return t
	}

	//如果原类型string, 并且目标类型是int,尝试转换
	if vVal.Kind() == reflect.String && tVal.Kind() >= reflect.Int && tVal.Kind() <= reflect.Int64 {
		if num, err := strconv.ParseInt(v.(string), 10, 64); err == nil {
			return AnyConvert2T(num, t)
		}
	}

	//如果原类型是int, 并且目标类型是string,尝试转换
	if vVal.Kind() >= reflect.Int && vVal.Kind() <= reflect.Int64 && tVal.Kind() == reflect.String {
		vint64 := AnyConvert2T(v, int64(0))

		return AnyConvert2T(fmt.Sprintf("%d", vint64), t)
	}

	//如果原类型是float, 并且目标类型是string,尝试转换
	if (vVal.Kind() == reflect.Float32 || vVal.Kind() == reflect.Float64) && tVal.Kind() == reflect.String {
		vfloat64 := AnyConvert2T(v, float64(0))

		return AnyConvert2T(fmt.Sprintf("%f", vfloat64), t)
	}

	//如果原类型是string, 并且目标类型是float,尝试转换
	if vVal.Kind() == reflect.String && (tVal.Kind() == reflect.Float32 || tVal.Kind() == reflect.Float64) {
		if num, err := strconv.ParseFloat(v.(string), 64); err == nil {
			return AnyConvert2T(num, t)
		}
	}

	//其他情况
	if vVal.Type().ConvertibleTo(tVal.Type()) {
		return vVal.Convert(tVal.Type()).Interface().(T)
	}

	return t
}

func MarshalJson(v any) string {
	jon, err := json.Marshal(v)

	if err != nil {
		fmt.Printf("utilities.MarshalJson param %v error %v", v, err)
	}

	return string(jon)
}

func IsMap(data any) bool {
	return reflect.TypeOf(data).Kind() == reflect.Map
}

func IsStruct(data any) bool {
	val := reflect.ValueOf(data)
	switch val.Kind() {
	case reflect.Ptr:
		// 如果是指针类型，则检查它指向的对象是否为 struct
		return val.Elem().Kind() == reflect.Struct
	case reflect.Struct:
		// 如果本身就是 struct 类型
		return true
	default:
		// 其他情况返回 false
		return false
	}
}

// YYYY-MM-DD -> unix
// YYYY-MM-DD hh:mm:ss -> unix
func TimeString2Unix(t string) int64 {
	loc, _ := time.LoadLocation("Local") // 获取时区

	if len(t) <= 10 {
		t = fmt.Sprintf("%s 00:00:00", t)
	}
	timer, _ := time.ParseInLocation(time.DateTime, t, loc)

	return timer.Unix()
}

func TimeString2Time(t string) time.Time {
	loc, _ := time.LoadLocation("Local") // 获取时区

	if len(t) <= 10 {
		t = fmt.Sprintf("%s 00:00:00", t)
	}
	timer, _ := time.ParseInLocation(time.DateTime, t, loc)

	return timer
}

type Number interface {
	int | int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

// 安全的除法， 除数为0返回0和错误
func SafeDivide[T Number](numerator, denominator T) (T, error) {
	// 检查除数是否为 0
	if denominator == 0 {
		return 0, errors.New("division by zero")
	}

	return numerator / denominator, nil
}
