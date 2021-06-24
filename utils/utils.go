package utils

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/mozillazg/go-pinyin"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func GetMillisecond() int64 {
	return time.Now().UnixNano() / 1e6
}

func GetGenerateRandom(length int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	randTab := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	size := len(randTab)
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = randTab[r.Intn(size)]
	}
	return string(bytes)
}

func VerifyEmail(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func VerifyMobile(mobile string) bool {
	if len(mobile) > 11 {
		return false
	}
	pattern := `^[0-9]*$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(mobile)
}

func VerifyIp (ip string) bool {
	ip = strings.Trim(ip, " ")
	pattern := `^(([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.)(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){2}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(ip)
}

//结构体转map[string]string
func ToMapString(in interface{}, tagName string) (map[string]string, error) {
	out := make(map[string]string)

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // 非结构体返回错误提示
		return nil, fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
	}

	t := v.Type()
	// 遍历结构体字段
	// 指定tagName值为map中key;字段值为map中value
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get(tagName); tagValue != "" {
			switch fi.Type.Kind() {
			case reflect.Int64, reflect.Int8, reflect.Int32, reflect.Int:
				out[tagValue] = strconv.FormatInt(v.Field(i).Int(), 10)
			case reflect.String:
				out[tagValue] = v.Field(i).String()
			case reflect.Slice:
				field, _ := json.Marshal(v.Field(i).Interface())
				out[tagValue] = string(field)
			}
		}
	}
	return out, nil
}

// 汉字转拼音
func HanToPinyin(keyword string) string {
	var (
		args   pinyin.Args
		result string
		single []string
	)
	args = pinyin.NewArgs()
	for _, v := range keyword {
		single = pinyin.LazyPinyin(string(v), args)
		if len(single) > 0 && single[0] != "" {
			result += single[0]
		} else {
			result += string(v)
		}
	}

	return result
}