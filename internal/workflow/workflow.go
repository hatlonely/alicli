package workflow

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hpifu/go-kit/href"
	"github.com/hpifu/go-kit/hstr"
)

type Workflow struct {
	ctx   *Ctx
	flows []interface{}
}

func NewWorkflow(global interface{}, define interface{}, workflow interface{}) (*Workflow, error) {
	ctx := NewCtx()
	if err := ctx.Set("global", global); err != nil {
		return nil, err
	}
	if err := ctx.Set("define", define); err != nil {
		return nil, err
	}
	w := &Workflow{ctx: ctx}
	if _, err := w.Evaluate(&define); err != nil {
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

func (w *Workflow) Evaluate(data *interface{}) (interface{}, error) {
	switch (*data).(type) {
	case string:
		v := (*data).(string)
		if strings.HasPrefix(v, "{{") && strings.HasSuffix(v, "}}") {
			text := v[2 : len(v)-2]
			if strings.HasPrefix(text, "type.int64(") {
				val, err := strconv.ParseInt(text[11:len(text)-1], 10, 64)
				if err != nil {
					return nil, err
				}
				*data = val
			} else {
				key := strings.TrimSpace(text)
				val, err := w.ctx.Get(key)
				if err != nil {
					return nil, err
				}
				*data = val
			}
		}
	case map[string]interface{}:
		for k, v := range (*data).(map[string]interface{}) {
			if val, err := w.Evaluate(&v); err != nil {
				return nil, err
			} else {
				(*data).(map[string]interface{})[k] = val
			}
		}
	case map[interface{}]interface{}:
		for k, v := range (*data).(map[interface{}]interface{}) {
			if val, err := w.Evaluate(&v); err != nil {
				return nil, err
			} else {
				(*data).(map[interface{}]interface{})[k] = val
			}
		}
	case []string:
		for i, v := range (*data).([]string) {
			var vi interface{}
			vi = v
			if val, err := w.Evaluate(&vi); err != nil {
				return nil, err
			} else {
				(*data).([]interface{})[i] = val
			}
		}
	case []interface{}:
		for i, v := range (*data).([]interface{}) {
			if val, err := w.Evaluate(&v); err != nil {
				return nil, err
			} else {
				(*data).([]interface{})[i] = val
			}
		}
	default:
	}

	return *data, nil
}

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
		if _, err := w.Evaluate(&info.Detail); err != nil {
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

type JobGenerator func(ctx *Ctx, plugins map[string]interface{}) Job

var typeJobMap = map[string]JobGenerator{}

func Register(typename string, generator JobGenerator) {
	typeJobMap[typename] = generator
}

func init() {
	Register("echo", NewEchoJob)
}

func (w *Workflow) CreateJob(info *JobInfo) Job {
	if constructor, ok := typeJobMap[info.Type]; ok {
		return constructor(w.ctx, info.Plugins)
	}

	return nil
}
