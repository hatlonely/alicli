package http

import (
	"fmt"
	"testing"

	"github.com/hatlonely/alicli/internal/workflow"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHttpJob_Do(t *testing.T) {
	Convey("test http", t, func() {
		job := &Job{}
		res, err := job.Do(map[interface{}]interface{}{
			"url": "http://www.baidu.com/s",
			"header": map[string]interface{}{
				"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36",
			},
			"method": "POST",
			"params": map[interface{}]interface{}{
				"wd": "golang",
			},
			"form": map[interface{}]interface{}{},
		})
		So(err, ShouldBeNil)
		fmt.Println(res)
		So(res, ShouldNotBeNil)
	})
}

func TestWorkFlow(t *testing.T) {
	Convey("case 1", t, func() {
		ctx := workflow.NewCtx()
		ctx.Set("accessKeyID", "xxx")
		ctx.Set("accessKeySecret", "xxx")

		w := workflow.NewWorkFlow(ctx)

		err := w.Run([]interface{}{
			map[interface{}]interface{}{
				"description": "测试阿里云",
				"type":        "http",
				"plugins": map[interface{}]interface{}{
					"aliyunpop": "",
				},
				"detail": map[interface{}]interface{}{
					"method": "POST",
					"url":    "https://imm.cn-shanghai.aliyuncs.com",
					"params": map[interface{}]interface{}{
						"Action": "ListProjects",
					},
				},
			},
			//map[interface{}]interface{}{
			//	"description": "测试 http",
			//	"type":        "http",
			//	"detail": map[interface{}]interface{}{
			//		"method": "POST",
			//		"url":    "http://www.baidu.com/s",
			//		"header": map[string]interface{}{
			//			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36",
			//		},
			//		"params": map[interface{}]interface{}{
			//			"wd": "golang",
			//		},
			//		"form": map[interface{}]interface{}{},
			//	},
			//},
		})
		So(err, ShouldBeNil)
	})
}
