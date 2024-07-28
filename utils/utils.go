package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CamelName 下划线转驼峰
func CamelName(s string) string {
	a := strings.Split(s, "_")
	str := ""
	for _, v := range a {
		if len(v) < 1 {
			return ""
		}
		strArry := []rune(v)
		if strArry[0] >= 97 && strArry[0] <= 122 {
			strArry[0] -= 32
		}
		str += string(strArry)
	}
	return str
}

// InArray 判断是否在切片里
func InArray(need interface{}, needArr interface{}) bool {
	arr := needArr.([]string)
	for _, v := range arr {
		if need == v {
			return true
		}
	}

	return false
}

// NowFormat 返回当前格式化2006-01-02 15:04:05时间
func NowFormat() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// TimeParse 返回格式化时间的时间戳
func TimeParse(timeStr string) int64 {
	t, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	return t.Unix()
}

// FileExists 判断文件是否存在
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// MkDir 创建文件夹
func MkDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// GetCurrDir 获取当前文件夹路径
func GetCurrDir() (string, error) {
	dir, err := os.Executable()
	if err != nil {
		return "", err
	}

	return strings.Replace(filepath.Dir(dir), "\\", "/", -1), nil
}

// FileIsExists 文件是否存在
func FileIsExists(path string) bool {
	if _, err := os.Stat(path); os.IsExist(err) {
		return true
	} else if err != nil {
		return false
	}
	return true
}

// SaveFile 保存文件(覆盖)
func SaveFile(file, data string) error {
	fileObj, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer fileObj.Close()

	if _, err = io.WriteString(fileObj, data); err != nil {
		return err
	}

	return nil
}

func HasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}

func HasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

// GetRandomString 随机生成指定位数的大写字母和数字的组合
func GetRandomString(l int) string {
	byt := []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	result := make([]byte, 0)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < l; i++ {
		result = append(result, byt[r.Intn(len(byt))])
	}

	return string(result)
}

// Md5 Md5加密
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
