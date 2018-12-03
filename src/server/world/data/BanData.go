package data

import (
	"server/common"
	"base"
	"log"
	"strings"
)

type(
	BanData struct{
		BanName string
	}

	BanDataRes struct{
		common.BaseDataRes
	}
)

var(
	BANDATA	BanDataRes
)

func (this *BanDataRes) Read() bool {
	this.Init()
	var file base.CDataFile
	lineData := &base.RData{}

	if (!file.ReadDataFile("data/BanWord.dat")) {
		log.Fatalf("read BanWord.dat error")
		return false
	}

	for i := 0; i < file.RecordNum; i++{
		pData := BanData{}

		file.GetData(lineData)
		base.IFAssert(lineData.Type == base.DType_String, "read BanWord.dat BanName error")
		pData.BanName = lineData.String
		if pData.BanName == ""{
			continue
		}

		//this.AddData(pData.BanName, pData)
	}

	return true
}

func (this *BanDataRes) GetData(id int) *BanData {
	pData := this.BaseDataRes.GetBaseData(id)
	if pData != nil{
		return pData.(*BanData)
	}

	return nil
}

func ReplaceBanWord(msg string, replace string) string{
	for i,_ := range BANDATA.DataMap{
		msg = strings.Replace(msg, i.(string), replace, -1 )
	}
	return msg
}