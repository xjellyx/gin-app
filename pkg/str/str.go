package str

import (
	"bytes"
	"crypto/rand"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"math/big"
)

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	var result string
	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err)
		}
		result += string(charset[index.Int64()])
	}

	return result
}

// UTF82GBK :  transform UTF8 rune into GBK byte array
func UTF82GBK(str string) ([]byte, error) {
	gb18030 := simplifiedchinese.All[0]
	return io.ReadAll(transform.NewReader(bytes.NewReader([]byte(str)), gb18030.NewEncoder()))
}

// GBK2UTF8 : transform  GBK byte array into UTF8 string
func GBK2UTF8(src []byte) (string, error) {
	gb18030 := simplifiedchinese.All[0]
	b, err := io.ReadAll(transform.NewReader(bytes.NewReader(src), gb18030.NewDecoder()))
	return string(b), err
}
