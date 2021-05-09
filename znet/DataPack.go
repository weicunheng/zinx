package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"zinx/utils"
	"zinx/zface"
)

type DataPack struct {
}

// 构造方法
func NewDataPack() DataPack  {
	return DataPack{}
}

// 获取Message Head的长度
func (d *DataPack) GetHeadLen() int {
	// HeadLen = DataLen + IdLen
	return 8
}

// 封包，接受消息，封装成bytes数组
func (d *DataPack) Pack(message zface.IMessage) ([]byte, error) {
	// 把message 打包为二进制数据
	// 1. 创建一个存放byte缓冲；本质就是二进制切片， 存放二进制组
	byteBuffer := bytes.NewBuffer([]byte{})

	// 2. 通过binary.Write(w io.Writer, order ByteOrder, data interface{}) error 将输入写入[]byte中
	if err := binary.Write(byteBuffer, binary.LittleEndian, message.GetMsgLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(byteBuffer, binary.LittleEndian, message.GetMsgId()); err != nil {
		return nil, err
	}
	if err := binary.Write(byteBuffer, binary.LittleEndian, message.GetMsgData()); err != nil {
		return nil, err
	}

	byteBufferBytes := byteBuffer.Bytes()
	fmt.Println("byteBufferBytes:", byteBufferBytes)
	return byteBufferBytes, nil

}

// 拆包, 接受bytes数组，返回拆包后的Message
func (d *DataPack) UnPack(message []byte) (zface.IMessage, error) {
	// 1. 创建一个BytesRead
	messageByteData := bytes.NewReader(message)
	msg := Message{}
	// 2. 读取message data len
	if err := binary.Read(messageByteData, binary.LittleEndian, &msg.MsgLen); err != nil {
		return nil, err
	}
	// 3. 读取message id
	if err := binary.Read(messageByteData, binary.LittleEndian, &msg.MsgId); err != nil {
		return nil, err
	}
	// 4. 读取message data
	if err := binary.Read(messageByteData, binary.LittleEndian, &msg.Msg); err != nil {
		return nil, err
	}
	if utils.GlobalObject.MaxPackageSize > 0 && msg.GetMsgLen() > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("数据包过大！！")
	}

	// 这就意味着指针类型的receiver 方法实现接口时，只有指针类型的对象实现了该接口。
	// 对应上面的例子来说，只有&msg实现了GetMsgData接口，而msg根本没有实现该接口。
	return &msg, nil
}
