package trafficproxyclient

import proxy "github.com/Telenav/osrm-backend/integration/pkg/trafficproxy"

// Eater is the interface that wraps the basic Eat method.
type Eater interface {

	// Eat consumes traffic responses.
	Eat(proxy.TrafficResponse)
}

// Feeder will continuesly feed traffic flows and incidents.
type Feeder struct {
	e []Eater
}

// NewFeeder creates a new traffic flows and incidents Feeder.
func NewFeeder() *Feeder {
	tf := Feeder{[]Eater{}}
	return &tf
}

// RegisterEaters add eaters for this feeder.
func (f *Feeder) RegisterEaters(e ...Eater) {
	f.e = append(f.e, e...)
}

// Run starts to feed traffic flows and incidents if any.
// It'll block until `Shutdown` called.
func (f *Feeder) Run() {
	//TODO:
}

// Shutdown stops the feed process.
func (f *Feeder) Shutdown() {
	//TODO:
}
