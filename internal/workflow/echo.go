package workflow

import "github.com/hpifu/go-kit/href"

type EchoJob struct {}

type EchoJobDetail struct {
	Message string
}

func NewEchoJob(ctx *Ctx, plugins map[string]interface{}) Job {
	return &EchoJob{}
}

func (j *EchoJob) Do(v interface{}) (interface{}, error) {
	detail := &EchoJobDetail{}
	if err := href.InterfaceToStruct(v, detail); err != nil {
		return nil, err
	}
	return detail.Message, nil
}
