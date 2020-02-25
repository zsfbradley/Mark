package mFace

import "markV5/mConst"

/*
  配置:能够支持用户自定义，但是必须指定输出几个服务器参数，能够加载与重载
  服务器:能够支持自定义配置或使用默认配置，默认加载，手动启动，手动停止，手动重载
*/

type MServeLoad interface {
	Load() MError
	Reload() MError
}

type MServeStatus interface {
	Status() mConst.MServeStatus
}

type MServeStart interface {
	Start() MError
}

type MServeStop interface {
	Stop() MError
}

type MServeEnding interface {
	StartEnding() MError
	OfficialEnding() MError
}

type MServeDataChannel interface {
	DataInChannel() MChannel
	DataOutChannel() MChannel
}
