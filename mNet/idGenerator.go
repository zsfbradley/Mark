package mNet

import (
	"../mTool"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	idg *idGenerator
)

func DefaultIDG() *idGenerator {
	return idg
}

func InitializesIDG(center, worker string) error {
	if len(center) > 5 || len(worker) > 5 {
		return errors.New("both of center and worker length need less or equal 5")
	}
	idg = &idGenerator{
		time:   time.Now().UnixNano(),
		center: center,
		worker: worker,
		count:  int64(0),
	}
	idg.parse()
	return nil
}

type idGenerator struct {
	time   int64  // len = 41
	center string // len = 5
	worker string // len = 5
	count  int64  // len = 12
	sync.Mutex
}

func (idg *idGenerator) parse() {
	centerS := 5 - len(idg.center)
	format := fmt.Sprintf("%%0%dd%%s",centerS)
	idg.center = fmt.Sprintf(format,0 , idg.center)

	workerS := 5 - len(idg.worker)
	format = fmt.Sprintf("%%0%dd%%s",workerS)
	idg.worker = fmt.Sprintf(format,0 , idg.worker)
}

func (idg *idGenerator) NewID() string {
	idg.Lock()
	defer idg.Unlock()

	now := time.Now().UnixNano()
	if now != idg.time {
		idg.time = now
		idg.count = int64(0)
	}

	idg.count++

	id := fmt.Sprintf("0%041d%s%s%012d",
		idg.time,
		idg.center,
		idg.worker,
		idg.count)
	return id
}

func (idg *idGenerator) EncNewID() string {
	return mTool.MD5(idg.NewID())
}
