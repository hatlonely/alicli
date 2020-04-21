package main

import (
	"fmt"
	"os"

	"github.com/hpifu/go-kit/hconf"
	"github.com/hpifu/go-kit/hflag"

	"github.com/hatlonely/alicli/internal/ctx"
	"github.com/hatlonely/alicli/internal/http"
	"github.com/hatlonely/alicli/internal/ots"
	"github.com/hatlonely/alicli/internal/workflow"
)

var AppVersion = "unknown"

type Options struct {
	CtxFile  string `hflag:"--ctx-file, -c; usage: context file path"`
	WorkFile string `hflag:"--work-file, -w; required; usage: work file path"`
}

func init() {
	NewHTTPJobWithAliyunPOP := func(ctx *ctx.Ctx, plugins map[string]interface{}) workflow.Job {
		job := http.NewJob(ctx, plugins)
		hjob := job.(*http.Job)
		hjob.AddPlugin("aliyunpop", nil)
		return hjob
	}

	workflow.Register("http", http.NewJob)
	workflow.Register("ots", ots.NewJob)
	workflow.Register("imm", NewHTTPJobWithAliyunPOP)
	workflow.Register("kms", NewHTTPJobWithAliyunPOP)
	workflow.Register("ecs", NewHTTPJobWithAliyunPOP)
	workflow.Register("nas", NewHTTPJobWithAliyunPOP)
}

func main() {
	version := hflag.Bool("v", false, "print current version")
	options := &Options{}

	if err := hflag.Bind(options); err != nil {
		panic(err)
	}
	if err := hflag.Parse(); err != nil {
		fmt.Println(hflag.Usage())
		panic(err)
	}
	if *version {
		fmt.Println(AppVersion)
		os.Exit(0)
	}

	ctxConfig, err := hconf.New("yaml", "local", options.CtxFile)
	if err != nil {
		panic(err)
	}
	workConfig, err := hconf.New("yaml", "local", options.WorkFile)
	if err != nil {
		panic(err)
	}
	define, err := workConfig.Get("define")
	if err != nil {
		panic(err)
	}
	flows, err := workConfig.Get("workflow")
	if err != nil {
		panic(err)
	}
	global, err := ctxConfig.Get("")
	if err != nil {
		panic(err)
	}
	wf, err := workflow.NewWorkflow(global, define, flows)
	if err != nil {
		panic(err)
	}
	if err := wf.Run(); err != nil {
		panic(err)
	}
}
