package main

import (
	"actor"
	"message"
	"github.com/golang/protobuf/proto"
	"fmt"
	"base"
)

type (
	EventProcess struct {
		actor.Actor

		AccountId int
		PlayerId int
		AccountName string
		SimId int32
	}

	IEventProcess interface {
		actor.IActor
		LoginGame()
		LoginAccount()
	}
)

func SendPacket(packet proto.Message){
	buff := message.Encode(packet)
	buff = base.SetTcpEnd(buff)
	CLIENT.Send(buff)
}

func (this *EventProcess) PacketFunc(socketid int, buff []byte) bool {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("EventProcess PacketFunc", err)
		}
	}()

	packetId, data := message.Decode(buff)
	packet := message.GetPakcet(packetId)
	if packet == nil{
		return true
	}
	err := proto.Unmarshal(data, packet)
	if err == nil{
		bitstream := base.NewBitStream(make([]byte, 1024), 1024)
		if !message.GetProtoBufPacket(packet, bitstream) {
			return true
		}
		var io actor.CallIO
		io.Buff = bitstream.GetBuffer()
		io.SocketId = socketid
		this.Send(io)
		return true
	}

	return true
}

func (this *EventProcess) Init(num int) {
	this.Actor.Init(num)
	this.RegisterCall("W_C_SelectPlayerResponse", func(packet *message.W_C_SelectPlayerResponse) {
		this.AccountId = int(*packet.AccountId)
		nLen := len(packet.PlayerData)
		//fmt.Println(len(packet.PlayerData), this.AccountId, packet.PlayerData)
		if nLen == 0{
			packet1 := &message.C_W_CreatePlayerRequest{PacketHead:message.BuildPacketHead( this.AccountId, int(message.SERVICE_WORLDSERVER)),
				PlayerName:proto.String("我是大坏蛋"),
				Sex:proto.Int32(int32(0)),}
			SendPacket(packet1)
		}else{
			this.PlayerId = int(*packet.PlayerData[0].PlayerID)
			this.LoginGame()
		}
	})

	this.RegisterCall("W_C_CreatePlayerResponse", func(packet *message.W_C_CreatePlayerResponse) {
		if *packet.Error == 0 {
			this.PlayerId = int(*packet.PlayerId)
			this.LoginGame()
		}else{//创建失败

		}
	})

	this.RegisterCall("A_C_LoginRequest", func(packet *message.A_C_LoginRequest) {
		if *packet.Error == base.ACCOUNT_NOEXIST {
			packet1 := &message.C_A_RegisterRequest{PacketHead:message.BuildPacketHead( 0, int(message.SERVICE_ACCOUNTSERVER)),
				AccountName:packet.AccountName, SocketId: proto.Int32(0)}
			SendPacket(packet1)
		}
	})

	this.RegisterCall("A_C_RegisterResponse", func(packet *message.A_C_RegisterResponse) {
		//注册失败
		if *packet.Error != 0 {
		}
	})

	this.RegisterCall("W_C_LoginMap", func(packet *message.W_C_LoginMap) {
		this.SimId = *packet.Id
		fmt.Println("login map")
	})

	this.RegisterCall("W_C_Move", func(packet *message.W_C_Move) {
		if this.SimId == *packet.Id{
		}else{
			fmt.Printf("simobj:[%d], Pos:[x:%d, y:%d, z:%d], Rot[%d]", *packet.Id, *packet.Pos.X, *packet.Pos.Y, *packet.Pos.Z, *packet.Rotation)
		}
		//this.Move(0, 100.0)
	})
	this.Actor.Start()
}

func (this *EventProcess)  LoginGame(){
	packet1 := &message.C_W_Game_LoginRequset{PacketHead:message.BuildPacketHead( this.AccountId, int(message.SERVICE_WORLDSERVER)),
		PlayerId:proto.Int32(int32(this.PlayerId)),}
	SendPacket(packet1)
}

var(

)

func (this *EventProcess)  LoginAccount() {
	this.AccountName = fmt.Sprintf("test%d", base.RAND().RandI(3005, 3005))
	packet1 := &message.C_A_LoginRequest{PacketHead: message.BuildPacketHead(0, int(message.SERVICE_ACCOUNTSERVER)),
		AccountName: proto.String(this.AccountName), BuildNo: proto.String(base.BUILD_NO), SocketId: proto.Int32(0)}
	SendPacket(packet1)
}

var(
	PACKET *EventProcess
)
