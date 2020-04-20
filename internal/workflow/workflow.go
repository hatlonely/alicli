package workflow

import (
	"errors"
	"fmt"

	"github.com/hpifu/go-kit/href"
	"github.com/hpifu/go-kit/hstr"

	"github.com/hatlonely/alicli/internal/ctx"
)

type Workflow struct {
	ctx   *ctx.Ctx
	flows []interface{}
}

func NewWorkflow(global interface{}, define interface{}, workflow interface{}) (*Workflow, error) {
	nc := ctx.NewCtx()
	if err := nc.Set("global", global); err != nil {
		return nil, err
	}
	if err := nc.Set("define", define); err != nil {
		return nil, err
	}
	w := &Workflow{ctx: nc}
	if _, err := nc.Evaluate(&define); err != nil {
		return nil, err
	}

	var ok bool
	if w.flows, ok = workflow.([]interface{}); !ok {
		return nil, errors.New("workflow is not []interface{}")
	}

	return w, nil
}

type JobInfo struct {
	Description string
	Type        string
	Plugins     map[string]interface{}
	Detail      interface{}
	Result      interface{}
}

type Job interface {
	Do(v interface{}) (interface{}, error)
}

var GreenTxt = hstr.NewFontStyle(hstr.ForegroundGreen)

func (w *Workflow) Run() error {
	for _, v := range w.flows {
		info := &JobInfo{}
		if err := href.InterfaceToStruct(v, info); err != nil {
			return err
		}
		job := w.CreateJob(info)
		if job == nil {
			return fmt.Errorf("unknown job type [%#v]", info)
		}
		fmt.Println(info.Description)
		fmt.Println(hstr.Indent("  ", info.Type))
		if _, err := w.ctx.Evaluate(&info.Detail); err != nil {
			panic(err)
		}
		fmt.Println(hstr.Indent("    ", "Detail"))
		fmt.Println(hstr.Indent("      ", hstr.ToYamlString(info.Detail)))
		res, err := job.Do(info.Detail)
		if err != nil {
			return err
		}
		info.Result = res
		fmt.Println(hstr.Indent("    ", "Result"))
		fmt.Println(GreenTxt.Render(hstr.Indent("      ", hstr.ToYamlString(info.Result))))
		fmt.Println()
	}

	return nil
}

type JobGenerator func(nc *ctx.Ctx, plugins map[string]interface{}) Job

var typeJobMap = map[string]JobGenerator{}

func Register(typename string, generator JobGenerator) {
	typeJobMap[typename] = generator
}

func (w *Workflow) CreateJob(info *JobInfo) Job {
	if constructor, ok := typeJobMap[info.Type]; ok {
		return constructor(w.ctx, info.Plugins)
	}

	return nil
}
