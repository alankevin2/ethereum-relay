package utils

import (
	"io/ioutil"
	"math"
	"math/rand"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// @title GetRandomNumber
// @description 取隨機整數
// @param min int 整數最小值
// @param max int 整數最大值
// @return res int 隨機整數
func GetRandomNumber(min int, max int) (res int) {
	rand.Seed(time.Now().UnixNano())
	return int(math.Floor(rand.Float64()*float64(max-min+1))) + min
}

// @title GetRandomString
// @description 生成隨機字串
// @param n int 字串長度
// @return string 指定長度字串
func GetRandomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		rand.Seed(time.Now().UnixNano())
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

// @title SlceStringInsert
// @description 將字串存入切片的指定index
// @param s []string 目標字串切片
// @param index int 指定index
// @param val string 欲存入字串
// @return res []string 調整後的新切片
func SlceStringInsert(s []string, index int, val string) (res []string) {
	head := s[:index]
	res = append(head, val)
	if len(s) > index+1 {
		end := s[index+1:]
		res = append(res, end...)
	}
	return res
}

// @title SlceIntInsert
// @description 將int元素存入切片的指定index
// @param s []int 目標int切片
// @param index int 指定index
// @param val int 欲存入int
// @return res []int 調整後的新切片
func SlceIntInsert(s []int, index int, val int) (res []int) {
	head := s[:index]
	res = append(head, val)
	if len(s) > index+1 {
		end := s[index+1:]
		res = append(res, end...)
	}
	return res
}

// @title ConvertJsonToMap
// @description 轉換json檔為map
// @param jsonPath string json檔存放path
// @return map[string]interface{}, error
func ConvertJsonToMap(jsonPath string) (res map[string]interface{}, err error) {
	file, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return nil, err
	}
	res = make(map[string]interface{})
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	_ = json.Unmarshal([]byte(file), &res)
	return res, nil
}
