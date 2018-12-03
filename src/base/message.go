package base

import (
	"reflect"
	"strings"
	"fmt"
	"unsafe"
)

type(
	Message interface {
		ReadData(*BitStream)
		WriteData(*BitStream)
	}
)

var(
	gMessageCreatorFactoryMap	map[string] func() Message
	gMessageCreatorFactoryInit	bool
)

func getMessageName(message Message) string{
	sType := strings.ToLower(reflect.ValueOf(message).Type().String())
	index := strings.Index(sType, ".")
	if index!= -1{
		sType = sType[index+1:]
	}
	return sType
}

func RegisterMessage(message Message)  {
	if !gMessageCreatorFactoryInit{
		gMessageCreatorFactoryMap = make(map[string] func() Message)
		gMessageCreatorFactoryInit = true
	}

	gMessageCreatorFactoryMap[getMessageName(message)] = func() Message {
		message1 := reflect.New(reflect.ValueOf(message).Elem().Type()).Interface().(Message)
		return message1
	}
}

func GetMessage(messageName string) Message{
	CreateFunc, exist := gMessageCreatorFactoryMap[messageName]
	if exist{
		return CreateFunc()
	}
	return nil
}

//--- 结构体写入bitsream
func WriteData(message interface{},bitstream *BitStream) bool{
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("WriteData", err)
		}
	}()
	return  parseMessage(bitstream, message)
}

//--- 结构体读取bitsream
func ReadData(message interface{}, bitstream *BitStream, ) bool{
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("ReadData", err)
		}
	}()
	return parseBMessage(bitstream, message)
}

func getTypeString(classField reflect.StructField, classVal reflect.Value) string{
	paramType := classField.Type
	sType := ""
	if paramType.Kind() == reflect.Ptr{
		sType = "*" + paramType.Elem().Kind().String()
	}else if paramType.Kind() == reflect.Slice{
		sType = GetSliceTypeString(paramType.String())
	}else if paramType.Kind() == reflect.Array{
		sType = GetArrayTypeString(paramType.String())
	} else{
		sType = classField.Type.Kind().String()
	}
	return sType
}

func parseType(bitstream *BitStream, classField reflect.StructField, val reflect.Value) bool {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("parseType", err)
		}
	}()

	sType := getTypeString(classField, val)
	switch sType {
	case "*bool":
		value := false
		if !val.IsNil() {
			value = *val.Interface().(*bool)
		}
		bitstream.WriteFlag(value)
	case "*float64":
		value := float64(0)
		if !val.IsNil() {
			value = *val.Interface().(*float64)
		}
		bitstream.WriteFloat64(value)
	case "*float32":
		value := float32(0)
		if !val.IsNil() {
			value = *val.Interface().(*float32)
		}
		bitstream.WriteFloat(value)
	case "*int8":
		value := int8(0)
		if !val.IsNil() {
			value = *val.Interface().(*int8)
		}
		bitstream.WriteInt(int(value), 8)
	case "*uint8":
		value := uint8(0)
		if !val.IsNil() {
			value = *val.Interface().(*uint8)
		}
		bitstream.WriteInt(int(value), 8)
	case "*int16":
		value := int16(0)
		if !val.IsNil() {
			value = *val.Interface().(*int16)
		}
		bitstream.WriteInt(int(value), 16)
	case "*uint16":
		value := uint16(0)
		if !val.IsNil() {
			value = *val.Interface().(*uint16)
		}
		bitstream.WriteInt(int(value), 16)
	case "*int32":
		value := int32(0)
		if !val.IsNil() {
			value = *val.Interface().(*int32)
		}
		bitstream.WriteInt(int(value), 32)
	case "*uint32":
		value := uint32(0)
		if !val.IsNil() {
			value = *val.Interface().(*uint32)
		}
		bitstream.WriteInt(int(value), 32)
	case "*int64":
		value := int64(0)
		if !val.IsNil() {
			value = *val.Interface().(*int64)
		}
		bitstream.WriteInt64(int64(value), 64)
	case "*uint64":
		value := uint64(0)
		if !val.IsNil() {
			value = *val.Interface().(*uint64)
		}
		bitstream.WriteInt64(int64(value), 64)
	case "*string":
		value := string("")
		if !val.IsNil() {
			value = *val.Interface().(*string)
		}
		bitstream.WriteString(value)
	case "*int":
		value := int(0)
		if !val.IsNil() {
			value = *val.Interface().(*int)
		}
		bitstream.WriteInt(value,32)
	case "*struct":
		value := val.Elem().Interface()
		parseMessage(bitstream, value)
	case "bool":
		bitstream.WriteFlag(val.Interface().(bool))
	case "float64":
		bitstream.WriteFloat64(val.Interface().(float64))
	case "float32":
		bitstream.WriteFloat(val.Interface().(float32))
	case "int8":
		bitstream.WriteInt(int(val.Interface().(int8)), 8)
	case "uint8":
		bitstream.WriteInt(int(val.Interface().(uint8)), 8)
	case "int16":
		bitstream.WriteInt(int(val.Interface().(int16)), 16)
	case "uint16":
		bitstream.WriteInt(int(val.Interface().(uint16)), 16)
	case "int32":
		bitstream.WriteInt(int(val.Interface().(int32)), 32)
	case "uint32":
		bitstream.WriteInt(int(val.Interface().(uint32)), 32)
	case "int64":
		bitstream.WriteInt64(int64(val.Interface().(int64)), 64)
	case "uint64":
		bitstream.WriteInt64(int64(val.Interface().(uint64)), 64)
	case "string":
		bitstream.WriteString(val.Interface().(string))
	case "int":
		bitstream.WriteInt(val.Interface().(int),32)
	case "struct":
		parseMessage(bitstream, val.Interface())
	case "[]bool":
		value := val.Interface().([]bool)
		nLen := len(value)
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			bitstream.WriteFlag(value[i])
		}
	case "[]float64":
		value := val.Interface().([]float64)
		nLen := len(value)
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			bitstream.WriteFloat64(value[i])
		}
	case "[]float32":
		value := val.Interface().([]float32)
		nLen := len(value)
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			bitstream.WriteFloat(value[i])
		}
	case "[]int8":
		value := val.Interface().([]int8)
		nLen := len(value)
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			bitstream.WriteInt(int(value[i]), 8)
		}
	case "[]uint8":
		value := val.Interface().([]uint8)
		nLen := len(value)
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			bitstream.WriteInt(int(value[i]), 8)
		}
	case "[]int16":
		value := val.Interface().([]int16)
		nLen := len(value)
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			bitstream.WriteInt(int(value[i]), 16)
		}
	case "[]uint16":
		value := val.Interface().([]uint16)
		nLen := len(value)
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			bitstream.WriteInt(int(value[i]), 16)
		}
	case "[]int32":
		value := val.Interface().([]int32)
		nLen := len(value)
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			bitstream.WriteInt(int(value[i]), 32)
		}
	case "[]uint32":
		value := val.Interface().([]uint32)
		nLen := len(value)
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			bitstream.WriteInt(int(value[i]), 32)
		}
	case "[]int64":
		value := val.Interface().([]int64)
		nLen := len(value)
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			bitstream.WriteInt64(value[i], 64)
		}
	case "[]uint64":
		value := val.Interface().([]uint64)
		nLen := len(value)
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			bitstream.WriteInt64(int64(value[i]), 64)
		}
	case "[]string":
		value := val.Interface().([]string)
		nLen := len(value)
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			bitstream.WriteString(value[i])
		}
	case "[]int":
		value := val.Interface().([]int)
		nLen := len(value)
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			bitstream.WriteInt(value[i], 32)
		}
	case "[]struct":
		nLen := val.Len()
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			parseMessage(bitstream,  val.Index(i).Interface())
		}
	case "[*]bool":
		nLen := val.Len()
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			bitstream.WriteFlag(val.Index(i).Bool())
		}
	case "[*]float64":
		nLen := val.Len()
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			bitstream.WriteFloat64(val.Index(i).Float())
		}
	case "[*]float32":
		nLen := val.Len()
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			bitstream.WriteFloat(float32(val.Index(i).Float()))
		}
	case "[*]int8":
		nLen := val.Len()
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			bitstream.WriteInt(int(val.Index(i).Int()), 8)
		}
	case "[*]uint8":
		nLen := val.Len()
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			bitstream.WriteInt(int(val.Index(i).Int()), 8)
		}
	case "[*]int16":
		nLen := val.Len()
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			bitstream.WriteInt(int(val.Index(i).Int()), 16)
		}
	case "[*]uint16":
		nLen := val.Len()
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			bitstream.WriteInt(int(val.Index(i).Int()), 16)
		}
	case "[*]int32":
		nLen := val.Len()
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			bitstream.WriteInt(int(val.Index(i).Int()), 32)
		}
	case "[*]uint32":
		nLen := val.Len()
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			bitstream.WriteInt(int(val.Index(i).Int()), 32)
		}
	case "[*]int64":
		nLen := val.Len()
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			bitstream.WriteInt64(val.Index(i).Int(), 64)
		}
	case "[*]uint64":
		nLen := val.Len()
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			bitstream.WriteInt64(val.Index(i).Int(), 64)
		}
	case "[*]string":
		nLen := val.Len()
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			bitstream.WriteString(val.Index(i).String())
		}
	case "[*]int":
		nLen := val.Len()
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			bitstream.WriteInt(int(val.Index(i).Int()), 32)
		}
	case "[*]struct":
		nLen := val.Len()
		bitstream.WriteInt(nLen, 16)
		for i := 0; i < nLen; i++ {
			parseMessage(bitstream,  val.Index(i).Interface())
		}
	default:
		fmt.Println("parseType type not supported", sType,  classField.Type)
		panic("parseType type not supported")
		return false
	}
	return true
}

func parseMessage(bitstream *BitStream, message interface{}) bool{
	var protoVal reflect.Value
	protoType := reflect.TypeOf(message)
	if protoType.Kind() == reflect.Ptr {
		protoType = reflect.TypeOf(message).Elem()
		protoVal = reflect.ValueOf(message).Elem()
	}else{
		protoVal = reflect.ValueOf(message)
	}
	for i := 0; i < protoType.NumField(); i++{
		if !protoVal.Field(i).CanSet(){//小写成员只有只读
			continue
		}

		if !parseType(bitstream, protoType.Field(i), protoVal.Field(i)){
			return false//丢弃这个包
		}
	}
	return true
}

func parseBMessage(bitstream *BitStream, message interface{}) bool {
	var protoVal reflect.Value
	protoType := reflect.TypeOf(message)
	if protoType.Kind() == reflect.Ptr {
		protoType = reflect.TypeOf(message).Elem()
		protoVal = reflect.ValueOf(message).Elem()
	}else{
		protoVal = reflect.ValueOf(message)
	}
	for i := 0; i < protoType.NumField(); i++{
		if !protoVal.Field(i).CanSet(){//小写成员只有只读
			continue
		}

		if !parseBType(bitstream, protoType.Field(i), protoVal.Field(i)){
			return false//丢弃这个包
		}
	}
	return true
}

func parseBType(bitstream *BitStream, classField reflect.StructField, val reflect.Value) bool{
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("parseBType", err)
		}
	}()
	sType := getTypeString(classField, val)
	switch sType {
	case "*bool":
		*val.Interface().(*bool) = bitstream.ReadFlag()
	case "*float64":
		*val.Interface().(*float64) = bitstream.ReadFloat64()
	case "*float32":
		*val.Interface().(*float32) = bitstream.ReadFloat()
	case "*int8":
		*val.Interface().(*int8) = int8(bitstream.ReadInt(8))
	case "*uint8":
		*val.Interface().(*uint8) = uint8(bitstream.ReadInt(8))
	case "*int16":
		*val.Interface().(*int16) = int16(bitstream.ReadInt(16))
	case "*uint16":
		*val.Interface().(*uint16) = uint16(bitstream.ReadInt(16))
	case "*int32":
		*val.Interface().(*int32) = int32(bitstream.ReadInt(32))
	case "*uint32":
		*val.Interface().(*uint32) = uint32(bitstream.ReadInt(32))
	case "*int64":
		*val.Interface().(*int64) = int64(bitstream.ReadInt64(64))
	case "*uint64":
		*val.Interface().(*uint64) = uint64(bitstream.ReadInt64(64))
	case "*string":
		*val.Interface().(*string) = bitstream.ReadString()
	case "*int":
		*val.Interface().(*int) = bitstream.ReadInt(32)
	case "*struct":
		value1 := val.Elem().Interface()
		parseBMessage(bitstream, value1)
	case "bool":
		val.SetBool(bitstream.ReadFlag())
	case "float64":
		val.SetFloat(bitstream.ReadFloat64())
	case "float32":
		val.SetFloat(float64(bitstream.ReadFloat()))
	case "int8":
		val.SetInt(int64(bitstream.ReadInt(8)))
	case "uint8":
		val.SetInt(int64(bitstream.ReadInt(8)))
	case "int16":
		val.SetInt(int64(bitstream.ReadInt(16)))
	case "uint16":
		val.SetInt(int64(bitstream.ReadInt(16)))
	case "int32":
		val.SetInt(int64(bitstream.ReadInt(32)))
	case "uint32":
		val.SetInt(int64(bitstream.ReadInt(32)))
	case "int64":
		val.SetInt(int64(bitstream.ReadInt64(64)))
	case "uint64":
		val.SetInt(int64(bitstream.ReadInt64(64)))
	case "string":
		val.SetString(bitstream.ReadString())
	case "int":
		val.SetInt(int64(bitstream.ReadInt(32)))
	case "struct":
		value1 := val.Interface()
		parseBMessage(bitstream, value1)
	case "[]bool":
		nLen := bitstream.ReadInt(16)
		for i := 0; i < nLen; i++ {
			val.Set( reflect.Append(val, reflect.ValueOf(bitstream.ReadFlag())))
		}
	case "[]float64":
		nLen := bitstream.ReadInt(16)
		for i := 0; i < nLen; i++ {
			val.Set( reflect.Append(val, reflect.ValueOf(bitstream.ReadFloat64())))
		}
	case "[]float32":
		nLen := bitstream.ReadInt(16)
		for i := 0; i < nLen; i++ {
			val.Set( reflect.Append(val, reflect.ValueOf(bitstream.ReadFloat())))
		}
	case "[]int8":
		nLen := bitstream.ReadInt(16)
		for i := 0; i < nLen; i++ {
			val.Set( reflect.Append(val, reflect.ValueOf(int8(bitstream.ReadInt(8)))))
		}
	case "[]uint8":
		nLen := bitstream.ReadInt(16)
		for i := 0; i < nLen; i++ {
			val.Set( reflect.Append(val, reflect.ValueOf(uint8(bitstream.ReadInt(8)))))
		}
	case "[]int16":
		nLen := bitstream.ReadInt(16)
		for i := 0; i < nLen; i++ {
			val.Set( reflect.Append(val, reflect.ValueOf(int16(bitstream.ReadInt(16)))))
		}
	case "[]uint16":
		nLen := bitstream.ReadInt(16)
		for i := 0; i < nLen; i++ {
			val.Set( reflect.Append(val, reflect.ValueOf(uint16(bitstream.ReadInt(16)))))
		}
	case "[]int32":
		nLen := bitstream.ReadInt(16)
		for i := 0; i < nLen; i++ {
			val.Set( reflect.Append(val, reflect.ValueOf(int32(bitstream.ReadInt(32)))))
		}
	case "[]uint32":
		nLen := bitstream.ReadInt(16)
		for i := 0; i < nLen; i++ {
			val.Set( reflect.Append(val, reflect.ValueOf(uint32(bitstream.ReadInt(32)))))
		}
	case "[]int64":
		nLen := bitstream.ReadInt(16)
		for i := 0; i < nLen; i++ {
			val.Set( reflect.Append(val, reflect.ValueOf(int64(bitstream.ReadInt(64)))))
		}
	case "[]uint64":
		nLen := bitstream.ReadInt(16)
		for i := 0; i < nLen; i++ {
			val.Set( reflect.Append(val, reflect.ValueOf(uint64(bitstream.ReadInt(64)))))
		}
	case "[]string":
		nLen := bitstream.ReadInt(16)
		for i := 0; i < nLen; i++ {
			val.Set( reflect.Append(val, reflect.ValueOf(bitstream.ReadString())))
		}
	case "[]int":
		nLen := bitstream.ReadInt(16)
		for i := 0; i < nLen; i++ {
			val.Set( reflect.Append(val, reflect.ValueOf(bitstream.ReadInt(32))))
		}
	case "[]struct"://no support
		value1 := val.Elem().Interface().([]interface{})
		nLen := bitstream.ReadInt(16)
		for i := 0; i < nLen; i++ {
			parseBMessage(bitstream, value1[i])
		}
	case "[*]bool":
		nLen := bitstream.ReadInt(16)
		arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
		for i := 0; i < nLen; i++ {
			value :=  (*bool)(unsafe.Pointer(arrayPtr))
			arrayPtr = arrayPtr + 1
			*value =  bool(bitstream.ReadFlag())
		}
	case "[*]float64":
		nLen := bitstream.ReadInt(16)
		arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
		for i := 0; i < nLen; i++ {
			value :=  (*float64)(unsafe.Pointer(arrayPtr))
			arrayPtr = arrayPtr + 8
			*value =  float64(bitstream.ReadFloat64())
		}
	case "[*]float32":
		nLen := bitstream.ReadInt(16)
		arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
		for i := 0; i < nLen; i++ {
			value :=  (*float32)(unsafe.Pointer(arrayPtr))
			arrayPtr = arrayPtr + 4
			*value =  float32(bitstream.ReadFloat64())
		}
	case "[*]int8":
		nLen := bitstream.ReadInt(16)
		arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
		for i := 0; i < nLen; i++ {
			value :=  (*int8)(unsafe.Pointer(arrayPtr))
			arrayPtr = arrayPtr + 1
			*value =  int8(bitstream.ReadInt(8))
		}
	case "[*]uint8":
		nLen := bitstream.ReadInt(16)
		arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
		for i := 0; i < nLen; i++ {
			value :=  (*uint8)(unsafe.Pointer(arrayPtr))
			arrayPtr = arrayPtr + 1
			*value =  uint8(bitstream.ReadInt(8))
		}
	case "[*]int16":
		nLen := bitstream.ReadInt(16)
		arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
		for i := 0; i < nLen; i++ {
			value :=  (*int16)(unsafe.Pointer(arrayPtr))
			arrayPtr = arrayPtr + 2
			*value =  int16(bitstream.ReadInt(16))
		}
	case "[*]uint16":
		nLen := bitstream.ReadInt(16)
		arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
		for i := 0; i < nLen; i++ {
			value :=  (*uint16)(unsafe.Pointer(arrayPtr))
			arrayPtr = arrayPtr + 2
			*value =  uint16(bitstream.ReadInt(16))
		}
	case "[*]int32":
		nLen := bitstream.ReadInt(16)
		arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
		for i := 0; i < nLen; i++ {
			value :=  (*int32)(unsafe.Pointer(arrayPtr))
			arrayPtr = arrayPtr + 4
			*value =  int32(bitstream.ReadInt(32))
		}
	case "[*]uint32":
		nLen := bitstream.ReadInt(16)
		arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
		for i := 0; i < nLen; i++ {
			value :=  (*uint32)(unsafe.Pointer(arrayPtr))
			arrayPtr = arrayPtr + 4
			*value =  uint32(bitstream.ReadInt(32))
		}
	case "[*]int64":
		nLen := bitstream.ReadInt(16)
		arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
		for i := 0; i < nLen; i++ {
			value :=  (*int64)(unsafe.Pointer(arrayPtr))
			arrayPtr = arrayPtr + 8
			*value =  int64(bitstream.ReadInt64(64))
		}
	case "[*]uint64":
		nLen := bitstream.ReadInt(16)
		arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
		for i := 0; i < nLen; i++ {
			value :=  (*uint64)(unsafe.Pointer(arrayPtr))
			arrayPtr = arrayPtr + 8
			*value =  uint64(bitstream.ReadInt64(64))
		}
	case "[*]string":
		nLen := bitstream.ReadInt(16)
		arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
		for i := 0; i < nLen; i++ {
			value :=  (*string)(unsafe.Pointer(arrayPtr))
			arrayPtr = arrayPtr + 16
			*value =  string(bitstream.ReadString())
		}
	case "[*]int":
		nLen := bitstream.ReadInt(16)
		arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
		for i := 0; i < nLen; i++ {
			value :=  (*int)(unsafe.Pointer(arrayPtr))
			arrayPtr = arrayPtr + 8
			*value =  int(bitstream.ReadInt(32))
		}
	default:
		fmt.Println("parseBType type not supported", sType,  classField.Type)
		panic("parseBType type not supported")
		return false
	}
	return true
}