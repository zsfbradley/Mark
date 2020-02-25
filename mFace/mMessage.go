package mFace

type MMessage interface {
	RouteID() string
	ConnID() string

	Marshal() []byte
	Unmarshal([]byte)
}
