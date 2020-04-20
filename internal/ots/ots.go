package ots

import (
	"reflect"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/hpifu/go-kit/href"

	"github.com/hatlonely/alicli/internal/ctx"
	"github.com/hatlonely/alicli/internal/workflow"
)

type Job struct {
	ctx     *ctx.Ctx
	plugins map[string]interface{}
}

type ClientInfo struct {
	Endpoint        string
	Instance        string
	AccessKeyID     string
	AccessKeySecret string
}

type JobDetail struct {
	Client ClientInfo
	Method string
	Params interface{}
}

func NewJob(ctx *ctx.Ctx, plugins map[string]interface{}) workflow.Job {
	return &Job{
		ctx:     ctx,
		plugins: plugins,
	}
}

func (j *Job) Do(v interface{}) (interface{}, error) {
	detail := &JobDetail{}
	if err := href.InterfaceToStruct(v, detail); err != nil {
		return nil, err
	}

	client := tablestore.NewClient(
		detail.Client.Endpoint, detail.Client.Instance,
		detail.Client.AccessKeyID, detail.Client.AccessKeySecret,
	)

	method := reflect.ValueOf(client).MethodByName(detail.Method)
	request := reflect.New(method.Type().In(0).Elem()).Interface()
	if err := href.InterfaceToStruct(detail.Params, request); err != nil {
		panic(err)
	}

	vals := method.Call([]reflect.Value{
		reflect.ValueOf(request),
	})
	var result []interface{}
	for _, val := range vals {
		result = append(result, val.Interface())
	}

	return result, nil
}
