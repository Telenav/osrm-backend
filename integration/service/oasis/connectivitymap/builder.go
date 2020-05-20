package connectivitymap

import (
	"sync"

	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/common"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"
	"github.com/golang/glog"
)

type connectivityMapBuilder struct {
	iterator      spatialindexer.PlacesIterator
	finder        spatialindexer.Finder
	ranker        spatialindexer.Ranker
	distanceLimit float64
	id2NearbyIDs  ID2NearByIDsMap

	numOfWorker         int
	workerWaitGroup     *sync.WaitGroup
	aggregatorWaitGroup *sync.WaitGroup
	aggregatorC         chan placeIDWithNearByPlaceIDs
}

func newConnectivityMapBuilder(iterator spatialindexer.PlacesIterator, finder spatialindexer.Finder,
	ranker spatialindexer.Ranker, distanceLimit float64, numOfWorker int) *connectivityMapBuilder {
	builder := &connectivityMapBuilder{
		iterator:      iterator,
		finder:        finder,
		ranker:        ranker,
		distanceLimit: distanceLimit,
		id2NearbyIDs:  make(ID2NearByIDsMap),

		numOfWorker:         numOfWorker,
		workerWaitGroup:     &sync.WaitGroup{},
		aggregatorWaitGroup: &sync.WaitGroup{},
		aggregatorC:         make(chan placeIDWithNearByPlaceIDs, 10000),
	}

	if numOfWorker < 1 {
		glog.Fatal("numOfWorker should never be smaller than 1, recommend using NumCPU()\n")
	}

	return builder
}

/*
                               ->  worker (fetch task  ->  find  ->  rank)
                             /                                               \
                            /                                                 \
  Input Iterator(channel)    --->  worker (fetch task  ->  find  ->  rank)      ---> aggregatorChannel -> feed to map
                            \                                                 /
                             \                                               /
                                ->  worker (fetch task  ->  find  ->  rank)

                                    . . .(more workers)
*/

func (builder *connectivityMapBuilder) build() ID2NearByIDsMap {
	builder.process()
	builder.aggregate()
	builder.wait()

	return builder.id2NearbyIDs
}

func (builder *connectivityMapBuilder) process() {
	inputC := builder.iterator.IteratePlaces()

	for i := 0; i < builder.numOfWorker; i++ {
		builder.workerWaitGroup.Add(1)
		go builder.work(i, inputC, builder.aggregatorC)
	}

	glog.Infof("builder's process is finished, start number of %d workers.\n", builder.numOfWorker)
}

func (builder *connectivityMapBuilder) work(workerID int, source <-chan common.PlaceInfo, sink chan<- placeIDWithNearByPlaceIDs) {
	defer builder.workerWaitGroup.Done()

	counter := 0
	for p := range source {
		counter += 1
		nearbyIDs := builder.finder.FindNearByPlaceIDs(*p.Location, builder.distanceLimit, spatialindexer.UnlimitedCount)
		rankedResults := builder.ranker.RankPlaceIDsByShortestDistance(*p.Location, nearbyIDs)

		sink <- placeIDWithNearByPlaceIDs{
			id:  p.ID,
			ids: rankedResults,
		}
	}

	glog.Infof("Worker_%d finished handling %d tasks.\n", workerID, counter)
}

func (builder *connectivityMapBuilder) aggregate() {
	builder.aggregatorWaitGroup.Add(1)

	go func() {
		counter := 0
		for item := range builder.aggregatorC {
			counter += 1
			builder.id2NearbyIDs[item.id] = item.ids
		}

		glog.Infof("Aggregation is finished with handling %d items.\n", counter)
		builder.aggregatorWaitGroup.Done()
	}()
}

func (builder *connectivityMapBuilder) wait() {
	builder.workerWaitGroup.Wait()
	close(builder.aggregatorC)
	builder.aggregatorWaitGroup.Wait()
}

type placeIDWithNearByPlaceIDs struct {
	id  common.PlaceID
	ids []*common.RankedPlaceInfo
}

func (builder *connectivityMapBuilder) buildInSerial() ID2NearByIDsMap {
	glog.Warning("This function is only used for compare result of worker's build().\n`")
	internalResult := make(chan placeIDWithNearByPlaceIDs, 10000)
	m := make(ID2NearByIDsMap)

	go func() {
		for p := range builder.iterator.IteratePlaces() {
			nearbyIDs := builder.finder.FindNearByPlaceIDs(*p.Location, builder.distanceLimit, spatialindexer.UnlimitedCount)
			rankedResults := builder.ranker.RankPlaceIDsByGreatCircleDistance(*p.Location, nearbyIDs)

			internalResult <- placeIDWithNearByPlaceIDs{
				id:  p.ID,
				ids: rankedResults,
			}
		}
		close(internalResult)
	}()

	for item := range internalResult {
		m[item.id] = item.ids
	}

	return m
}
