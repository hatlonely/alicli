package workflow

import (
	"fmt"

	"github.com/hpifu/go-kit/href"
	"github.com/hpifu/go-kit/hstr"
)

type WorkFlow struct {
	ctx *Ctx
}

func NewWorkFlow(ctx *Ctx) *WorkFlow {
	return &WorkFlow{
		ctx: ctx,
	}
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

func (w *WorkFlow) Run(vs []interface{}) error {
	for _, v := range vs {
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

type JobGenerator func(ctx *Ctx, plugins map[string]interface{}) Job

var typeJobMap = map[string]JobGenerator{}

func Register(typename string, generator JobGenerator) {
	typeJobMap[typename] = generator
}

func init() {
	Register("echo", NewEchoJob)
}

func (w *WorkFlow) CreateJob(info *JobInfo) Job {
	if constructor, ok := typeJobMap[info.Type]; ok {
		return constructor(w.ctx, info.Plugins)
	}

	return nil
}
