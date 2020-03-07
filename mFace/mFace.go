package mFace

import (
	"../mConst"
)

type MServeBase interface {
	Status() mConst.MServe_Status
	Load() error
	Start() error
	Reload() error
}

type MServeStop interface {
	Stop() error
}

type MServeEnding interface {
	StartEnding() error
	OfficialEnding() error
}
