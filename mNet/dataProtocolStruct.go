package mNet

import (
	"markV5/mError"
	"markV5/mFace"
)

func newDataProtocolStruct(symbol string, stages ...uint64) (mFace.MDataProtocolStruct, mFace.MError) {
	if symbol == "" {
		return nil, mError.NilParam
	}

	dps := &dpStruct{
		symbol: symbol,
		stages: make([]uint64, 0),
		data:   make([][]byte, 0),
	}

	dps.stages = append(dps.stages, uint64(len([]byte(symbol))))

	if len(stages) > 0 {
		for _, stage := range stages {
			dps.stages = append(dps.stages, uint64(stage))
		}
	}

	dps.stages = append(dps.stages, uint64(0)) // 尾部追加的 0 表示实际数据

	return dps, mError.Nil
}

// 真实数据强制对应的长度为0，追加在 stages 最后面
// 那么记录真实数据长度的 stage ， 强制为倒数第二位 ， 即为 stages[len(stages) - 2]
type dpStruct struct {
	symbol string
	stages []uint64 // 该参数保存的是每个 stage 的实际数据占用字节数，例如 "header" 占用6个字节 , 数字50占用4个字节
	data   [][]byte
}

func (dps *dpStruct) Symbol() string {
	return dps.symbol
}

func (dps *dpStruct) LengthOfStages() int {
	return len(dps.stages)
}

func (dps *dpStruct) ValueOfStages(index uint8) uint64 {
	if index > uint8(len(dps.stages)) {
		return 0
	}
	return dps.stages[index]
}

func (dps *dpStruct) AppendData(data []byte) {
	dps.data = append(dps.data, data)
}

func (dps *dpStruct) CleanData() {
	dps.data = make([][]byte, 0)
}

func (dps *dpStruct) Data() [][]byte {
	return dps.data
}

func (dps *dpStruct) SymbolLength() uint64 {
	return dps.stages[0]
}

// 整个头部长度 ， 最后实际数据长度的预设为0，所以追加上也没关系
func (dps *dpStruct) HeaderLength() uint64 {
	headerLength := uint64(0)
	for _, stage := range dps.stages {
		headerLength += stage
	}
	return headerLength
}

func (dps *dpStruct) RangeOfDataLength() (uint64, uint64) {
	location, length := uint64(0), uint64(0)
	for i := 0; i < len(dps.stages)-2; i++ {
		location += dps.stages[i]
	}
	length = dps.stages[len(dps.stages)-2]
	return location, length
}
