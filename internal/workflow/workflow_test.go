package workflow

import (
	"fmt"
	"testing"

	"github.com/hpifu/go-kit/href"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/alicli/internal/ctx"
)

func init() {
	Register("echo", NewEchoJob)
}

type EchoJob struct{}

type EchoJobDetail struct {
	Message string
}

func NewEchoJob(nc *ctx.Ctx, plugins map[string]interface{}) Job {
	return &EchoJob{}
}

func (j *EchoJob) Do(v interface{}) (interface{}, error) {
	detail := &EchoJobDetail{}
	if err := href.InterfaceToStruct(v, detail); err != nil {
		return nil, err
	}
	return detail.Message, nil
}

func TestWorkFlow_Evaluate(t *testing.T) {
	Convey("test evaluate", t, func() {
		w, err := NewWorkflow(map[string]interface{}{
			"accessKeyID": "access key",
			"Infos": map[string]interface{}{
				"key1": "val1",
				"key2": 2,
			},
		}, nil, nil)
		So(err, ShouldBeNil)

		mystruct := map[string]interface{}{
			"obj1": map[interface{}]interface{}{
				"ak": "{{global.accessKeyID}}",
				"sk": "xxx",
			},
			"obj2": "{{global.Infos}}",
		}

		var data interface{}
		data = mystruct
		w.Evaluate(&data)

		fmt.Println(mystruct)
	})
}

func TestWorkFlow(t *testing.T) {
	Convey("case 1", t, func() {
		w, err := NewWorkflow(nil, nil, []interface{}{
			map[interface{}]interface{}{
				"description": "测试 echo",
				"type":        "echo",
				"plugins":     map[interface{}]interface{}{"abc": "def", "def": 2},
				"detail": map[interface{}]interface{}{
					"message": "hello world",
				},
			},
			map[interface{}]interface{}{
				"description": "测试 echo",
				"type":        "echo",
				"plugins":     map[interface{}]interface{}{"abc": "def", "def": 2},
				"detail": map[interface{}]interface{}{
					"message": "hello world",
				},
			},
		})

		So(err, ShouldBeNil)
		So(w.Run(), ShouldBeNil)
	})
}
