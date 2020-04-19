package workflow

import (
	"fmt"
	"testing"

	"github.com/hpifu/go-kit/href"
	. "github.com/smartystreets/goconvey/convey"
)

func TestInterfaceToStruct(t *testing.T) {
	bj := &JobInfo{}
	Convey("test case1", t, func() {
		So(href.InterfaceToStruct(map[interface{}]interface{}{
			"description": "测试 echo",
			"type":        "Echo",
			"plugins":     []interface{}{"abc", "def"},
			"detail": map[interface{}]interface{}{
				"Message": "hello world",
			},
		}, bj), ShouldBeNil)

		fmt.Printf("%#v", bj)
	})
}
