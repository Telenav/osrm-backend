package trafficproxyclient

import (
	"sync"
	"time"

	"github.com/golang/glog"

	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic"
	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficproxy"
)

// Feeder will continuesly feed traffic flows and incidents.
type Feeder struct {
	e []livetraffic.Eater
}

// NewFeeder creates a new traffic flows and incidents Feeder.
func NewFeeder() *Feeder {
	tf := Feeder{[]livetraffic.Eater{}}
	return &tf
}

// RegisterEaters add eaters for this feeder.
func (f *Feeder) RegisterEaters(e ...livetraffic.Eater) {
	f.e = append(f.e, e...)
}

// Run starts to feed traffic flows and incidents if any.
// It'll block until `Shutdown` called or some error occurred.
func (f *Feeder) Run() error {

	feeds := make(chan trafficproxy.TrafficResponse)

	// feed eater
	waitFeedingDone := make(chan struct{})
	go func() {
		f.feed(feeds)
		waitFeedingDone <- struct{}{}
	}()

	//streaming delta
	deltaErr := make(chan error)
	go func() {
		deltaErr <- StreamingDeltaFlowsIncidents(feeds)
	}()
	time.Sleep(100 * time.Millisecond) //wait a while to make sure streaming delta running first

	//get all
	getAllError := make(chan error)
	go func() {

		allResp, err := GetFlowsIncidents(nil)
		if err != nil {
			getAllError <- err
			return
		}
		feeds <- *allResp
		getAllError <- nil
	}()
	if err := <-getAllError; err != nil {
		// only warning, won't exit directly
		glog.Warning(err)
	}

	// wait for delta
	err := <-deltaErr
	if err != nil {
		glog.Warning(err)
	}
	close(feeds)

	// wait for feeding exit
	<-waitFeedingDone

	return err // return err from streaming delta
}

// Shutdown stops the feed process.
func (f *Feeder) Shutdown() {
	//TODO:
}

func (f *Feeder) feed(in <-chan trafficproxy.TrafficResponse) {
	for {
		resp, ok := <-in
		if !ok {
			break
		}

		var wg sync.WaitGroup
		for _, e := range f.e {
			wg.Add(1)
			go func(e livetraffic.Eater) {
				e.Eat(resp)
				wg.Done()
			}(e)
		}
		wg.Wait()
	}
}
