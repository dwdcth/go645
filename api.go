package go645

type Client interface {
	ClientProvider
	//Read 发送读请求
	Read(ver ProtoVersion, address Address, itemCode int32) (*ReadData, bool, error)
	//ReadWithBlock  读请求使能块
	ReadWithBlock(ver ProtoVersion, address Address, data ReadRequestData) (*Protocol, error)
	//Broadcast 开始广播
	Broadcast(ver ProtoVersion, p InformationElement, control Control) error
}
