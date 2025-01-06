package xfile

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type xFile struct{}

var XFile = &xFile{}

// PathIsDir 判断所给路径是否为文件夹
func (x *xFile) PathIsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// CreateDir 创建目录
func (x *xFile) CreateDir(path string, perm os.FileMode) error {
	err := os.MkdirAll(path, perm)
	if err != nil {
		logrus.Error("create dir failure:", err)
	}
	return err
}

// PathIsFile 判断所给路径是否为文件
func (x *xFile) PathIsFile(path string) bool {
	return !x.PathIsDir(path)
}

func (x *xFile) PathIsExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// WriteFile 写文件
func (x *xFile) WriteFile(path string, data []byte) error {
	err := os.WriteFile(path, data, 0644)
	return err
}

func (x *xFile) ReadFile(path string) (data []byte, err error) {
	return os.ReadFile(path)
}

func (x *xFile) IteratorPath(path string) ([]fs.DirEntry, error) {
	return os.ReadDir(path)
}

func (x *xFile) ZipFile(srcDir, zipFileName string) {
	err := os.RemoveAll(zipFileName)
	if err != nil {
		logrus.Error("remove file ", zipFileName, "failure")
	}
	zipFileInfo, _ := os.Create(zipFileName)
	defer func(zipFileInfo *os.File) {
		err := zipFileInfo.Close()
		if err != nil {
			logrus.Error(err)
		}
	}(zipFileInfo)

	archive := zip.NewWriter(zipFileInfo)
	defer func(archive *zip.Writer) {
		err := archive.Close()
		if err != nil {
			logrus.Error(err)
		}
	}(archive)

	_ = filepath.Walk(srcDir, func(path string, info fs.FileInfo, err error) error {
		logrus.Info("zip file src: ", srcDir, ",path:", path)
		if path == srcDir {
			return nil
		}
		header, _ := zip.FileInfoHeader(info)
		// header.Name = strings.TrimPrefix(path, src_dir+`\`)
		// log.Println(header.Name)
		if info.IsDir() {
			header.Name += `/`
		} else {
			header.Method = zip.Deflate
		}
		writer, _ := archive.CreateHeader(header)
		if !info.IsDir() {
			file, _ := os.Open(path)
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					logrus.Error(err)
				}
			}(file)
			_, err = io.Copy(writer, file)
		}

		return nil
	})
	logrus.Info("zip ", zipFileName, "finish...")
}

func (x *xFile) GetDirSize(dir string) (size string, err error) {
	cmd := exec.Command("du", "-sh", dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	rst, err := simplifiedchinese.GBK.NewDecoder().Bytes(out)
	info := string(rst)

	if err != nil {
		return fmt.Sprintf("[%s]", info), err
	}
	sizeReg := regexp.MustCompile(`(\S+)?\s+`)
	sizeInfo := sizeReg.FindAllStringSubmatch(info, -1)
	if len(sizeInfo) > 0 {
		return sizeInfo[0][1], nil
	}
	return fmt.Sprintf("[%s],out=%s", cmd.String(), info), errors.New("unable to get the size of the file or directory")
}

func (x *xFile) GetDirSizeToBytesCount(dir string) (count int64, err error) {
	size, err := x.GetDirSize(dir)
	if err != nil {
		return 0, fmt.Errorf("%s-%s", err.Error(), size)
	}
	size = strings.ToUpper(size)
	tempSize := 0
	tempSize, err = strconv.Atoi(size[:len(size)-1])
	if err != nil {
		return 0, err
	}
	if strings.HasSuffix(size, "K") {
		count = int64(tempSize) * 1024
	} else if strings.HasSuffix(size, "M") {
		count = int64(tempSize) * 1024 * 1024
	} else if strings.HasSuffix(size, "G") {
		count = int64(tempSize) * 1024 * 1024 * 1024
	} else if strings.HasSuffix(size, "T") {
		count = int64(tempSize) * 1024 * 1024 * 1024 * 1024
	} else if strings.HasSuffix(size, "P") {
		count = int64(tempSize) * 1024 * 1024 * 1024 * 1024 * 1024
	}
	return
}

func (x *xFile) ConvertSizeFromStrToBytes(sizeStr string) (count int64, err error) {

	size := strings.ToUpper(sizeStr)
	tempSize := 0
	tempSize, err = strconv.Atoi(size[:len(size)-1])
	if err != nil {
		return 0, err
	}
	if strings.HasSuffix(size, "K") {
		count = int64(tempSize) * 1024
	} else if strings.HasSuffix(size, "M") {
		count = int64(tempSize) * 1024 * 1024
	} else if strings.HasSuffix(size, "G") {
		count = int64(tempSize) * 1024 * 1024 * 1024
	} else if strings.HasSuffix(size, "T") {
		count = int64(tempSize) * 1024 * 1024 * 1024 * 1024
	} else if strings.HasSuffix(size, "P") {
		count = int64(tempSize) * 1024 * 1024 * 1024 * 1024 * 1024
	}
	return
}
