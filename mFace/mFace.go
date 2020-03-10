package mFace

import (
	"../mConst"
)

type MServe interface {
	Status() mConst.MServe_Status
	Start() error
}

type MServeLoad interface {
	Load() error
	Reload() error
}

type MServeStop interface {
	Stop() error
}

type MServeEnding interface {
	StartEnding() error
	OfficialEnding() error
}
