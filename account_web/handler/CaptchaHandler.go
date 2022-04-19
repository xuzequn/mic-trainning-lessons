package handler

import (
	"encoding/base64"
	"fmt"
	"github.com/dchest/captcha"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"os"
)

func GenCaptcha() error {
	fileName := "data.png"
	f, err := os.Create(fileName)
	if err != nil {
		zap.S().Error("GenCaptcha() 失败")
		return err
	}
	defer f.Close()
	var w io.WriterTo
	d := captcha.RandomDigits(captcha.DefaultLen)
	w = captcha.NewImage("", d, captcha.StdWidth, captcha.StdHeight)
	_, err = w.WriteTo(f)
	if err != nil {
		zap.S().Error("GenCaptcha() 失败")
		return err
	}
	fmt.Println(d)
	captcha := ""
	for _, item := range d {
		captcha += fmt.Sprintf("%d", item)
	}
	fmt.Println(captcha)
	b64, err := GetBase64(fileName)
	if err != nil {
		return err
	}
	fmt.Println(b64)
	return nil
}

func GetBase64(fileName string) (string, error) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	b := make([]byte, len(file)*4)
	base64.StdEncoding.Encode(b, file)
	s := string(b)
	return s, nil
}
