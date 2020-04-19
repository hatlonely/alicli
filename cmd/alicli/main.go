package main

import (
	"fmt"
	"os"

	"github.com/hpifu/go-kit/hconf"
	"github.com/hpifu/go-kit/hflag"

	"github.com/hatlonely/alicli/internal/http"
	"github.com/hatlonely/alicli/internal/ots"
	"github.com/hatlonely/alicli/internal/workflow"
)

var AppVersion = "unknown"

type Options struct {
	CtxFile  string `hflag:"--ctx-file, -c; usage: context file path"`
	WorkFile string `hflag:"--work-file, -w; required; usage: work file path"`
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
	kvs := map[string]string{}
	if err := ctxConfig.Unmarshal(&kvs); err != nil {
		panic(err)
	}

	ctx := workflow.NewCtx()
	for k, v := range kvs {
		ctx.Set(fmt.Sprintf("global.%v", k), v)
	}

	workConfig, err := hconf.New("yaml", "local", options.WorkFile)
	if err != nil {
		panic(err)
	}
	workflow.Register("http", http.NewJob)
	workflow.Register("ots", ots.NewJob)
	wf := workflow.NewWorkFlow(ctx)
	flows, err := workConfig.Get("workflow")
	if err != nil {
		panic(err)
	}
	if err := wf.Run(flows.([]interface{})); err != nil {
		panic(err)
	}
}
