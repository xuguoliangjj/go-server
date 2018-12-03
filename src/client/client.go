package main

import (
	"base"
	"fmt"
	"log"
	"network"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
)

var (
	CLIENT *network.ClientSocket
)
func main() {
	cfg := &base.Config{}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	//初始ini配置文件
	dir += "/SXZ_SERVER.CFG"
	cfg.Read(dir)
	UserNetIP, UserNetPort := cfg.Get2("NetGate_WANAddress", ":")
	//UserNetIP, UserNetPort := "101.132.178.159", "31700"
	port,_ := strconv.Atoi(UserNetPort)
	CLIENT = new(network.ClientSocket)
	CLIENT.Init(UserNetIP, port)
	PACKET := new(EventProcess)
	PACKET.Init(1)
	CLIENT.BindPacketFunc(PACKET.PacketFunc)
	if CLIENT.Start(){
		PACKET.LoginAccount()
	}

	InitCmd()
	//PACKET.LoginGame()
	//for{
	//	PACKET.LoginAccount()
	//}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	fmt.Printf("client exit ------- signal:[%v]", s)
}