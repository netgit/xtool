package xhttp

import (
	"crypto/tls"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

type xHttp struct{}

var XHttp = &xHttp{}

func (x *xHttp) DoPost(url string, header map[string]string, post map[string]interface{}) string {
	client := resty.New()
	if strings.HasPrefix(url, "https") {
		client = client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	if header == nil {
		header = map[string]string{
			"User-Agent": x.Chrome(false),
		}
	}
	client = client.SetHeaders(header)
	res, err := client.R().SetBody(post).Post(url)
	if err != nil {
		logrus.Error("http request error:", err)
		return ""
	}

	return string(res.Body())
}

func (x *xHttp) DoPostObj(url string, header map[string]string, post interface{}) (string, *http.Response) {
	client := resty.New()
	if strings.HasPrefix(url, "https") {
		client = client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	if header == nil {
		header = map[string]string{
			"User-Agent": x.Chrome(false),
		}
	}
	client = client.SetHeaders(header)
	res, err := client.R().SetBody(post).Post(url)
	if err != nil {
		logrus.Error("http request error:", err)
		return "", res.RawResponse
	}

	return string(res.Body()), res.RawResponse
}

func (x *xHttp) DoDel(url string, header map[string]string, post map[string]interface{}) string {
	client := resty.New()
	if strings.HasPrefix(url, "https") {
		client = client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	if header == nil {
		header = map[string]string{
			"User-Agent": x.Chrome(false),
		}
	}
	client = client.SetHeaders(header)
	res, err := client.R().SetBody(post).Delete(url)
	if err != nil {
		logrus.Error("http request error:", err)
		return ""
	}

	return string(res.Body())
}

func (x *xHttp) DoGet(url string, header map[string]string, data map[string]interface{}) string {
	client := resty.New()
	if strings.HasPrefix(url, "https") {
		client = client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	if header == nil {
		header = map[string]string{
			"User-Agent": x.Chrome(false),
		}
	}
	client = client.SetHeaders(header)
	res, err := client.R().SetBody(data).Get(url)
	if err != nil {
		logrus.Error("http request error:", err)
		return ""
	}

	return string(res.Body())
}

func (x *xHttp) DoPut(url string, header map[string]string, data map[string]interface{}) string {
	client := resty.New()
	if strings.HasPrefix(url, "https") {
		client = client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	if header == nil {
		header = map[string]string{
			"User-Agent": x.Chrome(false),
		}
	}
	client = client.SetHeaders(header)
	res, err := client.R().SetBody(data).Put(url)
	if err != nil {
		logrus.Error("http request error:", err)
		return ""
	}

	return string(res.Body())
}

// DoMove 方法,待测试
func (x *xHttp) DoMove(url string, header map[string]string, data map[string]interface{}) string {
	client := resty.New()
	if strings.HasPrefix(url, "https") {
		client = client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	if header == nil {
		header = map[string]string{
			"User-Agent": x.Chrome(false),
		}
	}
	client = client.SetHeaders(header)
	res, err := client.R().SetBody(data).Execute("MOVE", url)
	if err != nil {
		logrus.Error("http request error:", err)
		return ""
	}

	return string(res.Body())
}

func (x *xHttp) DoGetDownload(url, dirPath, saveFileName string, header map[string]string, data map[string]interface{}) string {
	client := resty.New()
	//这里是否需要判断文件夹是否存在
	if !pathIsExists(dirPath) {
		logrus.Error("directory is not exists")
		return ""
	}
	client.SetOutputDirectory(dirPath)
	if strings.HasPrefix(url, "https") {
		client = client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	if header == nil {
		header = map[string]string{
			"User-Agent": x.Chrome(false),
		}
	}
	client = client.SetHeaders(header)
	res, err := client.R().SetOutput(saveFileName).SetBody(data).Get(url)
	if err != nil {
		logrus.Error("http request error:", err)
		return ""
	}

	return string(res.Body())
}

func (x *xHttp) DoPostDownload(url, dirPath, saveFileName string, header map[string]string, data map[string]interface{}) string {
	client := resty.New()
	//这里是否需要判断文件夹是否存在
	if !pathIsExists(dirPath) {
		logrus.Error("directory is not exists")
		return ""
	}
	client.SetOutputDirectory(dirPath)
	if strings.HasPrefix(url, "https") {
		client = client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	if header == nil {
		header = map[string]string{
			"User-Agent": x.Chrome(false),
		}
	}
	client = client.SetHeaders(header)
	res, err := client.R().SetOutput(saveFileName).SetBody(data).Post(url)
	if err != nil {
		log.Println("http request error:", err)
		return ""
	}

	return string(res.Body())
}

// Chrome 获取Chrome浏览器User-Agent,isRand=true随机获取一个
func (x *xHttp) Chrome(isRand bool) string {
	chrome := uaList["chrome"]
	if chrome != nil && len(chrome) > 0 {
		if !isRand {
			return chrome[0]
		}
		return chrome[rand.Intn(len(chrome))]
	}
	return defaultUA
}

// RandUserAgent 随机获取浏览器User-Agent
func (x *xHttp) RandUserAgent() string {
	index := rand.Intn(len(uaList))
	i := 0
	subUserAgent := make([]string, 0)
	for _, ua := range uaList {
		if i == index {
			subUserAgent = ua
		}
	}
	if len(subUserAgent) > 0 {
		return subUserAgent[rand.Intn(len(subUserAgent))]
	}
	return defaultUA
}

func pathIsExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
