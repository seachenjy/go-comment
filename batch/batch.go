package batch

import (
	"sync"

	"github.com/seachenjy/go-comment/dao"
	"github.com/seachenjy/go-comment/log"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

//Work batch work
type Work struct {
	SourceID dao.SourceID
}

type scheduler struct{}

var schedulers []*scheduler
var works []*Work
var m = &sync.Mutex{}
var ready chan *scheduler

//PushWork every comment post triggers the statistics work
func PushWork(w *Work) {
	m.Lock()
	defer m.Unlock()
	works = append(works, w)
}

//Run run works
//Count scores and records by sourceID
func Run() {
	ready = make(chan *scheduler)
	l := rate.NewLimiter(20, 5)
	for i := 0; i < 100; i++ {
		schedulers = append(schedulers, &scheduler{})
	}
	for {
		if l.Allow() && len(works) > 0 && len(schedulers) > 0 {
			m.Lock()
			s := schedulers[0]
			w := works[0]
			schedulers = schedulers[1:]
			works = works[1:]
			s.do(w, ready)
			m.Unlock()
			log.GetLogger().WithFields(logrus.Fields{
				"works":      len(works),
				"schedulers": len(schedulers),
			}).Info("batch")
		}
	}
}

func (s *scheduler) do(w *Work, ready chan *scheduler) {

}
