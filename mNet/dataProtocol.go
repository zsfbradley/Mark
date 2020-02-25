package mNet

import (
	"markV5/mError"
	"markV5/mFace"
	"markV5/mTool"
)

func newDataProtocolWithDPS(dps mFace.MDataProtocolStruct) (mFace.MDataProtocol, mFace.MError) {
	if dps == nil {
		return nil, mError.NilParam
	}

	dp := &dataProtocol{
		dps:             dps,
		buffer:          make([]byte, 0),
		compDataChannel: make(chan [][]byte),
	}
	return dp, mError.Nil
}

func newDataProtocol(symbol string, stages ...uint64) (mFace.MDataProtocol, mFace.MError) {
	dps, err := newDataProtocolStruct(symbol, stages...)
	if err.NotNil() {
		return nil, err
	}

	dp := &dataProtocol{
		dps:             dps,
		buffer:          make([]byte, 0),
		compDataChannel: make(chan [][]byte),
	}
	return dp, mError.Nil
}

type dataProtocol struct {
	dps             mFace.MDataProtocolStruct
	buffer          []byte
	compDataChannel chan [][]byte // 不需要缓存空间，每个 Conn 都会独立拥有一个数据处理器
}

func (dp *dataProtocol) Unmarshal(data []byte) {
	dp.buffer = append(dp.buffer, data...)
	length := uint64(len(dp.buffer))

	i := uint64(0)
	for ; i < length; i++ { // 遍历整个缓存区
		if i+dp.dps.HeaderLength() > length { // 取出完整头部信息
			break // 如果不行则跳出，等待下一段数据
		}
		if string(dp.buffer[i:i+dp.dps.SymbolLength()]) != dp.dps.Symbol() { // 判断标志
			continue // 如果不行则结束本 i 位置，继续往后
		}
		location, ldLength := dp.dps.RangeOfDataLength()
		dataLength := mTool.ByteToInt(dp.buffer[i+location : i+location+ldLength]) // 获取尾部数据的长度
		if i+location+ldLength+dataLength > length {                               // 如果剩余长度不足以取出全部数据则跳出，等待下一段数据
			break
		}

		// 从 i 位置开始，按照 dp.dps.stages的下标顺序取出各个stage的实际数据，并保存，然后将游标移到取出数据的后一位
		for x := 0; x < dp.dps.LengthOfStages(); x++ {
			addUp := dp.dps.ValueOfStages(uint8(x))
			// 如果 dp.dps.stage[x] == 0 , 说明到了获取真实数据的时候，此时的尾部数据长度应为上面取得的 dataLength
			if addUp == 0 {
				addUp = dataLength
			}
			dp.dps.AppendData(dp.buffer[i : i+addUp])
			//dp.dps.data = append(dp.dps.data, dp.buffer[i:i+addUp])
			i += addUp
		}

		dp.compDataChannel <- dp.dps.Data()
		dp.dps.CleanData()

		i -= 1 // 应为在取数据的过程中都是长度，且结束此过程后有 i++的过程，所以要减去一个下标
	}

	if i != length {
		leftData := dp.buffer[i:]
		dp.buffer = make([]byte, 0)
		dp.buffer = append(dp.buffer, leftData...)
	} else {
		dp.buffer = make([]byte, 0)
	}
}

// 省略symbol、尾部数据长度
func (dp *dataProtocol) Marshal(data ...[]byte) []byte {
	datas := make([]byte, 0)
	if len(data)+2 != dp.dps.LengthOfStages() {
		return datas
	}

	datas = append(datas, []byte(dp.dps.Symbol())...)
	for i := 0; i < len(data)-1; i++ {
		datas = append(datas, data[i]...)
	}
	datas = append(datas, mTool.IntToByte(uint64(len(data[len(data)-1])))...)
	datas = append(datas, data[len(data)-1]...)

	return datas
}

func (dp *dataProtocol) CompletedDataChannel() chan [][]byte {
	return dp.compDataChannel
}
