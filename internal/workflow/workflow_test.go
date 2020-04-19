package workflow

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

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
