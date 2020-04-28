package httptool

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func Post(req interface{}, url string, rsp interface{}) ([]byte, error) {
	reqByt, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	rq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqByt))
	if err != nil {
		return nil, err
	}
	rq.Header.Set("Content-Type", "application/json")
	rq.Header.Set("Connection", "keep-alive")

	body, err := client.Do(rq)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(body.Body)
	if err != nil {
		return nil, err
	}
	if rsp != nil {
		if err := json.Unmarshal(b, rsp); err != nil {
			return nil, err
		}
	}
	return b, nil
}

func Get(url string) (body *http.Response, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Host", "jandan.net")
	req.Header.Set("Connection", "keep-alive")
	return client.Do(req)
}

func UploadFile(fileName string, fileData []byte, url string) (*http.Response, error) {

	rqbody := new(bytes.Buffer)

	mWriter := multipart.NewWriter(rqbody)

	iow, err := mWriter.CreateFormFile("file", fileName)
	if err != nil {
		return nil, err
	}

	iow.Write(fileData)

	mWriter.Close()

	req, err := http.NewRequest("POST", url, rqbody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", mWriter.FormDataContentType())

	client := &http.Client{}

	return client.Do(req)
}

func isFileExist(filename string, filesize int64) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		fmt.Println(info)
		return false
	}
	if filesize == info.Size() {
		fmt.Println("安装包已存在！", info.Name(), info.Size(), info.ModTime())
		return true
	}
	del := os.Remove(filename)
	if del != nil {
		fmt.Println(del)
	}
	return false
}

func AliDownloadFile(url string, localPath string, fb func(length, downLen int64)) (string, error) {
	var (
		fsize   int64
		buf     = make([]byte, 32*1024)
		written int64
	)

	client := new(http.Client)
	//client.Timeout = time.Second * 60 //设置超时时间
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	contents := resp.Header.Get("Content-Disposition")
	index := strings.Index(contents, "=")
	if index == 0 {
		return "", errors.New("content-Disposition error")
	}
	localPath = localPath + contents[index+1:]
	tmpFilePath := localPath + ".download"
	file, err := os.Create(tmpFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	if resp.Body == nil {
		return "", errors.New("body is null")
	}
	defer resp.Body.Close()
	//下面是 io.copyBuffer() 的简化版本
	for {
		//读取bytes
		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			//写入bytes
			nw, ew := file.Write(buf[0:nr])
			//数据长度大于0
			if nw > 0 {
				written += int64(nw)
			}
			//写入出错
			if ew != nil {
				err = ew
				break
			}
			//读取是数据长度不等于写入的数据长度
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
		if fb != nil {
			fb(fsize, written)
		}
	}
	if err == nil {
		file.Close()
		err = os.Rename(tmpFilePath, localPath)
	}
	return localPath, err
}
