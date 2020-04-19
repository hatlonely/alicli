package workflow

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWorkFlow_Evaluate(t *testing.T) {
	Convey("test evaluate", t, func() {
		ctx := NewCtx()
		ctx.Set("global.accessKeyID", "access key")
		ctx.Set("global.Infos", map[string]interface{}{
			"key1": "val1",
			"key2": 2,
		})
		w := NewWorkFlow(ctx)

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
		ctx := &Ctx{}
		w := NewWorkFlow(ctx)

		err := w.Run([]interface{}{
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
	})
}
