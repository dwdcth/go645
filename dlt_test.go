package go645

import (
	"fmt"

	"github.com/goburrow/serial"
	"log"
	"testing"
	"time"
)

/*

前导字节
在主站发送帧信息之前，先发送4个字节FEH，以唤醒接收方。

收到命令帧后的响应延时 Td：20ms≤Td≤500ms。
字节之间停顿时间 Tb：Tb≤500ms
*/
/*
FE FE FE FE 68 08 02 02 41 22 00 68 11 04 35 37 33 37 2A 16
68 08 02 02 41 22 00 68 91 0A 35 37 33 37 3B 35 35 74 55 33 51 16
68 08 02 02 41 22 00 68   93 06  3B 35 35 74 55 33  79 16

68 01 00 00 00 00 00 68  93 06  34 33 33 33 33 33 9D 16
解析数据：
04000402  002241020208
3B 35 35 74 55 33
08 02 02 41 22 00
00 22 41 02 02 08

FE FE FE FE(前导) 68(起始) 01 00 00 00 00 00(表号)68(起始) 91(控制码) 08(数据长度) 33 34 34 35(数据类型) 77 66 55 44(数据值)   B0(校验码) 16(停止位)
其中数据按照bcd编码+33h，小端序
33 34 34 35 77 66 55 44  => 11223344
0  1   1  2 44 33 22 11
*/

func TestDlt(t *testing.T) {
	//bytes.Reader{}
	/*
	   	发送命令： FE FE FE FE 68 AA AA AA AA AA AA 68 13 00 DF 16

	      回复:68 01 00 00 00 00 00 68 93 06 34 33 33 33 33 33 9D 16

	      其中电表地址:34 33 33 33 33 33减去0x33为01 00 00 00 00 00,然后倒序后通信地址为00 00 00 00 00 01 控制码为0x93 长度:0x06 CS校验:0x9d 帧尾：0x16
	*/
	//dlta()

	//dltb()

	//dltc()

	//read1997()
	read2007()
}

//func dlta() {
//	//调用ClientProvider的构造函数, 返回结构体指针
//	p := dlt.NewClientProvider()
//	//windows 下面就是 com开头的，比如说 com3
//	//mac OS/linux/unix  下面就是 /dev/下面的，比如说 dev/tty.usbserial-14320
//	p.Address = "com1"
//	p.BaudRate = 9600
//	p.DataBits = 8
//	p.Parity = "E" //偶数
//	p.StopBits = 1
//	p.Timeout = 100 * time.Millisecond
//
//	client := dltcon.NewClient(p)
//	client.LogMode(true)
//	err := client.Start()
//	if err != nil {
//		fmt.Println("start err,", err)
//		return
//	}
//	addr := "000000000001"
//	//MeterNumber是表号 005223440001
//	//DataMarker是数据标识别 02010300
//	test := &dlt.Dlt645ConfigClient{addr, "00010000"}
//	for {
//		value, err := test.SendMessageToSerial(client)
//		if err != nil {
//			fmt.Println("readHoldErr,", err)
//		} else {
//			fmt.Printf("%#v\n", value)
//		}
//		time.Sleep(time.Second * 3)
//	}
//}

//33 33 34 33 9a 98
//00 00 01 00 67 65
// 11227973

// 68 01 00 00 00 00 00 68 91 08 33 33 34 33 59 b3 55 44 dc 16
// 11228026

var item = []int32{
	0x0001_0000, 0x90_10, //正向有功
	0x0201_0100, 0xb6_11, //a 相电压
	0x0202_0100, 0xb6_21, //a 相电流
	0x0203_0000, 0xb6_30, //瞬时有功
	0x0204_0000, 0xb6_40, //瞬时无功
	0x0002_0000, 0x90_20, //反向有功
	0x0003_0000, 0x91_10, //正向无功
	0x0004_0000, 0x91_20, //反向无功
	0x0206_0000, 0xB6_50, //总功率因数
}

type MyLog struct {
}

func (m MyLog) Write(address string, station string, data []byte) {
	fmt.Println(address, station, data)
}

func read2007() {
	p2 := NewRTUClientProvider(
		WithSerialConfig(serial.Config{
			Address:  "COM1",
			BaudRate: 2400,
			DataBits: 8,
			StopBits: 1,
			Parity:   "E",
			Timeout:  time.Second * 30}),
		WithEnableLogger(),
		WithLogSaver(MyLog{}),
	)
	//c1 := NewClient(p1)
	c2 := NewClient(p2)

	for i := 0; i < len(item); i = i + 2 {
		// 68 01 00 00 00 00 00 68 11 04 33 33 34 33 B3 16
		// 68 01 00 00 00 00 00 68 11 04 33 33 34 33 B3 16
		//pr, _, err := c.Read(NewAddress("000000000001", LittleEndian), 0x00_01_00_00) //正向有功电能
		fmt.Println("***************************")
		//pr, _, err := c1.Read(NewAddress("000000000001", LittleEndian), item[i]) //a 相电压
		//if err != nil {
		//	log.Print(err.Error())
		//} else {
		//	println(pr.GetValue())
		//}
		//time.Sleep(500 * time.Millisecond)
		pr, _, err := c2.Read(Ver2007, NewAddress("000000000001", LittleEndian), item[i]) //a 相电压
		if err != nil {
			log.Print(err.Error())
		} else {
			println(pr.GetValue())
		}
		time.Sleep(500 * time.Millisecond)
	}
	c2.Close()
}

func read1997() {

	p2 := NewRTUClientProvider(
		WithSerialConfig(serial.Config{
			Address:  "COM1",
			BaudRate: 2400,
			DataBits: 8,
			StopBits: 1,
			Parity:   "E",
			Timeout:  time.Second * 30}),
		WithEnableLogger(),
	)
	//c1 := NewClient(p1)
	c2 := NewClient(p2)

	for i := 0; i < len(item); i = i + 2 {
		// 68 01 00 00 00 00 00 68 11 04 33 33 34 33 B3 16
		// 68 01 00 00 00 00 00 68 11 04 33 33 34 33 B3 16
		//pr, _, err := c.Read(NewAddress("000000000001", LittleEndian), 0x00_01_00_00) //正向有功电能
		fmt.Println("***************************")
		//pr, _, err := c1.Read(NewAddress("000000000001", LittleEndian), item[i]) //a 相电压
		//if err != nil {
		//	log.Print(err.Error())
		//} else {
		//	println(pr.GetValue())
		//}
		//time.Sleep(500 * time.Millisecond)
		pr, _, err := c2.Read(Ver2007, NewAddress("000000000001", LittleEndian), item[i+1]) //a 相电压
		if err != nil {
			log.Print(err.Error())
		} else {
			println(pr.GetValue())
		}
		time.Sleep(500 * time.Millisecond)
	}
	c2.Close()
}
