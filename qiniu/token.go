package qiniu

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// Base64编码
func safeBase64Encode(buf []byte) string {
	s := base64.StdEncoding.EncodeToString(buf)

	s = strings.Replace(s, "+", "-", -1)
	s = strings.Replace(s, "/", "_", -1)

	return s
}

// 上传Token
func sumUploadToken(access_key, secret_key string, params map[string]interface{}) string {
	delete(params, "upload_token")

	buf, _ := json.Marshal(params)
	s := safeBase64Encode(buf)

	secret := hmac.New(sha1.New, []byte(secret_key))
	secret.Write([]byte(s))
	sign := safeBase64Encode(secret.Sum(nil))

	return fmt.Sprintf("%s:%s:%s", access_key, sign, s)
}

// 管理Token (文件删除/移动/复制、获取仓库/文件列表/空间大小...)
func sumMgToken(access_key, secret_key, uri, content_type, post string) string {
	def_ctp := "application/x-www-form-urlencoded"

	s := ""
	if uri != "" {
		u, err := url.Parse(uri)
		if err == nil {
			if u.Path != "" {
				s = fmt.Sprintf("%s%s", s, u.Path)
			}
			if u.RawQuery != "" {
				s = fmt.Sprintf("%s?%s", s, u.RawQuery)
			}
		}
	}
	s += "\n"
	if post != "" && content_type == def_ctp {
		s = fmt.Sprintf("%s%s", s, post)
	}
	secret := hmac.New(sha1.New, []byte(secret_key))
	secret.Write([]byte(s))
	sign := safeBase64Encode(secret.Sum(nil))

	return fmt.Sprintf("QBox %s:%s", access_key, sign)
}

// 下载Token
func sumDownToken(access_key, secret_key, uri string) string {
	secret := hmac.New(sha1.New, []byte(secret_key))
	secret.Write([]byte(uri))
	sign := safeBase64Encode(secret.Sum(nil))
	return fmt.Sprintf("%s:%s", access_key, sign)
}
