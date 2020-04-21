package ctx

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCtx_EvaluateExpr(t *testing.T) {
	Convey("test 1", t, func() {
		ctx := NewCtx()
		ctx.EvaluateExpr("int64(64)")
		ctx.EvaluateExpr("int64()")
		ctx.EvaluateExpr("int64")
	})
}
