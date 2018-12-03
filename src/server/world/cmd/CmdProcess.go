package cmd

import (
	"actor"
	"fmt"
	"runtime"
	"strconv"
	"server/world/toprank"
	"server/common"
)

type (
	CmdProcess struct {
		actor.Actor
	}

	ICmdProcess interface {
		actor.IActor
	}
)

func (this *CmdProcess) Init(num int) {
	this.Actor.Init(num)
	this.RegisterCall("cpus", func() {
		fmt.Println(runtime.NumCPU(), " cpus and ", runtime.GOMAXPROCS(0), " in use")
	})

	this.RegisterCall("routines", func() {
		fmt.Println("Current number of goroutines: ", runtime.NumGoroutine())
	})

	this.RegisterCall("setcpus", func(args string) {
		n, _ := strconv.Atoi(args)
		runtime.GOMAXPROCS(n)
		fmt.Println(runtime.NumCPU(), " cpus and ", runtime.GOMAXPROCS(0), " in use")
	})

	this.RegisterCall("startgc", func() {
		runtime.GC()
		fmt.Println("gc finished")
	})

	this.RegisterCall("InTopRank", func(argv0,argv1,argv2,argv3,argv4,argv5 string) {
		nType, _ := strconv.Atoi(argv0)
		id, _ := strconv.Atoi(argv1)
		name := argv2
		score, _ := strconv.Atoi(argv3)
		val0, _ := strconv.Atoi(argv4)
		val1, _ := strconv.Atoi(argv5)
		toprank.TOPMGR.SendMsg( "InTopRank", nType, uint64(id), name, score, val0, val1)
	})

	this.Actor.Start()
}

var(
	g_Cmd *CmdProcess
)

func Init(){
	g_Cmd = &CmdProcess{}
	g_Cmd.Init(1000)
	common.StartConsole(g_Cmd)
	InitWeb()
}
