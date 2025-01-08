package xhash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"net"
	"os"
)

type xHash struct{}

var XHash = &xHash{}

// Base64Encode base64加密
func (x *xHash) Base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// Base64Decode base64 解密
func (x *xHash) Base64Decode(data string) (string, error) {
	decode, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	return string(decode), nil
}

func (x *xHash) ParseIPFromDomain(domain string) ([]string, error) {
	n, err1 := net.LookupHost(domain)
	if err1 != nil {
		return nil, err1
	}
	return n, nil
}

func (x *xHash) Md5(data []byte, salt ...string) string {
	m := md5.New()

	m.Write(data)
	if len(salt) > 0 {
		m.Write([]byte(salt[0]))
	}
	enc := m.Sum(nil)
	return fmt.Sprintf("%x", enc)

}

func (x *xHash) FileMd5Hash(fileName string) (string, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return x.Md5(data), nil

}

func (x *xHash) Sha1(data []byte, salt ...string) string {
	m := sha1.New()
	m.Write(data)
	if len(salt) > 0 {
		m.Write([]byte(salt[0]))
	}
	enc := m.Sum(nil)
	return fmt.Sprintf("%x", enc)
}

func (x *xHash) Sha256(data []byte, salt ...string) string {
	m := sha256.New()
	m.Write(data)
	if len(salt) > 0 {
		m.Write([]byte(salt[0]))
	}
	enc := m.Sum(nil)
	return fmt.Sprintf("%x", enc)
}

func (x *xHash) Sha512(data []byte, salt ...string) string {
	m := sha512.New()
	m.Write(data)
	if len(salt) > 0 {
		m.Write([]byte(salt[0]))
	}
	enc := m.Sum(nil)
	return fmt.Sprintf("%x", enc)
}
