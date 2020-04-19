package http

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/url"
	"sort"
	"strings"
	"time"
)

var myrand = rand.New(rand.NewSource(time.Now().UnixNano()))

func PercentEncode(str string) string {
	str = url.QueryEscape(str)
	str = strings.ReplaceAll(str, "+", "%20")
	str = strings.ReplaceAll(str, "*", "%2A")
	str = strings.ReplaceAll(str, "%7E", "~")
	return str
}

func Signature(methods string, params map[string]interface{}, accessKeySecret string) string {
	var kvs []string
	for k, v := range params {
		kvs = append(kvs, fmt.Sprintf("%v=%v", PercentEncode(k), PercentEncode(v.(string))))
	}
	sort.Strings(kvs)
	str := strings.Join(kvs, "&")
	toSign := methods + "&" + PercentEncode("/") + "&" + PercentEncode(str)
	h := hmac.New(sha1.New, []byte(accessKeySecret+"&"))
	h.Write([]byte(toSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func MakePopParams(methods string, params map[string]interface{}, accessKeyID string, accessKeySecret string) map[string]interface{} {
	if _, ok := params["AccessKeyId"]; !ok {
		params["AccessKeyId"] = accessKeyID
	}
	if _, ok := params["Format"]; !ok {
		params["Format"] = "JSON"
	}
	if _, ok := params["Version"]; !ok {
		params["Version"] = "2017-09-06"
	}
	if _, ok := params["Timestamp"]; !ok {
		params["Timestamp"] = time.Now().UTC().Format(time.RFC3339)
	}
	if _, ok := params["SignatureMethod"]; !ok {
		params["SignatureMethod"] = "HMAC-SHA1"
	}
	if _, ok := params["SignatureVersion"]; !ok {
		params["SignatureVersion"] = "1.0"
	}
	if _, ok := params["SignatureNonce"]; !ok {
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, myrand.Uint64())
		params["SignatureNonce"] = hex.EncodeToString(buf)
	}
	params["Signature"] = Signature(methods, params, accessKeySecret)
	return params
}
