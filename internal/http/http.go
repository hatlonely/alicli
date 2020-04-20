package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hpifu/go-kit/href"
	"github.com/spf13/cast"

	"github.com/hatlonely/alicli/internal/ctx"
	"github.com/hatlonely/alicli/internal/workflow"
)

func init() {
	workflow.Register("http", NewJob)
}

type Job struct {
	ctx     *ctx.Ctx
	plugins map[string]interface{}
}

type JobDetail struct {
	URL    string
	Method string
	Header map[string]string
	Params map[string]interface{}
	Form   map[string]string
}

func NewJob(ctx *ctx.Ctx, plugins map[string]interface{}) workflow.Job {
	return &Job{
		ctx:     ctx,
		plugins: plugins,
	}
}

func (j *Job) AliyunPopPlugin(detail *JobDetail) (err error) {
	ak, ok := detail.Params["accessKeyID"]
	if !ok || ak == "" {
		ak, err = j.ctx.Get("global.accessKeyID")
		if err != nil {
			return err
		}
	}
	sk, ok := detail.Params["accessKeyID"]
	if !ok || sk == "" {
		sk, err = j.ctx.Get("global.accessKeySecret")
		if err != nil {
			return err
		}
	}

	for k, v := range detail.Params {
		detail.Params[k] = cast.ToString(v)
	}
	detail.Params = MakePopParams(detail.Method, detail.Params, ak.(string), sk.(string))

	return nil
}

func (j *Job) Do(v interface{}) (interface{}, error) {
	detail := &JobDetail{}
	if err := href.InterfaceToStruct(v, detail); err != nil {
		return nil, err
	}

	if _, ok := j.plugins["aliyunpop"]; ok {
		if err := j.AliyunPopPlugin(detail); err != nil {
			return nil, err
		}
	}
	var reader io.Reader
	if detail.Form != nil {
		buf, _ := json.Marshal(detail.Form)
		reader = bytes.NewReader(buf)
	}

	req, err := http.NewRequest(detail.Method, detail.URL, reader)
	if err != nil {
		return nil, err
	}

	if detail.Header != nil {
		for k, v := range detail.Header {
			req.Header.Add(k, v)
		}
	}

	if detail.Params != nil {
		values := &url.Values{}
		for k, v := range detail.Params {
			switch v.(type) {
			case string:
				values.Add(k, v.(string))
			case []string:
				for _, i := range v.([]string) {
					values.Add(k, i)
				}
			case int:
				values.Add(k, strconv.Itoa(v.(int)))
			case []int:
				for _, i := range v.([]int) {
					values.Add(k, strconv.Itoa(i))
				}
			case []interface{}:
				for _, i := range v.([]interface{}) {
					values.Add(k, cast.ToString(i))
				}
			default:
				values.Add(k, fmt.Sprintf("%v", v))
			}
		}
		req.URL.RawQuery = values.Encode()
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, 1000*time.Millisecond)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
		},
		Timeout: 5000 * time.Millisecond,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result interface{}
	if err := json.Unmarshal(buf, &result); err != nil {
		result = string(buf)
	}

	return map[string]interface{}{
		"header": NormalizeHeader(res.Header),
		"status": res.Status,
		"result": result,
	}, nil
}

func NormalizeHeader(header http.Header) []map[string]string {
	var kvs []map[string]string
	for k, vs := range header {
		for _, v := range vs {
			kvs = append(kvs, map[string]string{
				k: v,
			})
		}
	}

	return kvs
}
