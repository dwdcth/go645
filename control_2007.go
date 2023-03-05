package go645

func init() {
	controlVersionMap[Ver2007] = map[ControlKind]ControlType{
		IsSlave:  1 << 7,
		SlaveErr: 1 << 6,
		HasNext:  1 << 5,
		//Retain 保留
		Retain: 0b00000,
		//Broadcast 广播校时
		Broadcast: 0b01000,
		// ReadNext 读后续10001
		ReadNext: 0b10010,
		//ReadAddress 读通讯地址
		ReadAddress: 0b10011,
		//Write 写数据
		Write: 0b10100,
		//WriteAddress 读通讯地址
		WriteAddress: 0b10101,
		//ToChangeCommunicationRate 更改通讯速率
		ToChangeCommunicationRate: 0b10111,
		Freeze:                    0b10110,
		//PassWord 修改密码
		PassWord:       0b11000,
		ResetMaxDemand: 0b11001,
		//ResetEM 电表清零
		ResetEM:    0b11010,
		ResetEvent: 0b11011,
		Read:       0x11,
	}
}
