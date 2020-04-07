package connectivitymap

import (
	"sync"

	"github.com/Telenav/osrm-backend/integration/service/spatialindexer"
	"github.com/golang/glog"
)

type connectivityMapBuilder struct {
	iterator      spatialindexer.PointsIterator
	finder        spatialindexer.Finder
	ranker        spatialindexer.Ranker
	distanceLimit float64
	id2NearbyIDs  ID2NearByIDsMap

	numOfWorker         int
	dispatchWaitGroup   *sync.WaitGroup
	workerWaitGroup     *sync.WaitGroup
	aggregatorWaitGroup *sync.WaitGroup
	tasks4WorkerC       []chan spatialindexer.PointInfo
	aggregatorC         chan placeIDWithNearByPlaceIDs
}

func newConnectivityMapBuilder(iterator spatialindexer.PointsIterator, finder spatialindexer.Finder,
	ranker spatialindexer.Ranker, distanceLimit float64, numOfWorker int) *connectivityMapBuilder {
	builder := &connectivityMapBuilder{
		iterator:      iterator,
		finder:        finder,
		ranker:        ranker,
		distanceLimit: distanceLimit,
		id2NearbyIDs:  make(ID2NearByIDsMap),

		numOfWorker:         numOfWorker,
		dispatchWaitGroup:   &sync.WaitGroup{},
		workerWaitGroup:     &sync.WaitGroup{},
		aggregatorWaitGroup: &sync.WaitGroup{},
		tasks4WorkerC:       make([]chan spatialindexer.PointInfo, numOfWorker),
		aggregatorC:         make(chan placeIDWithNearByPlaceIDs, 10000),
	}

	if numOfWorker < 1 {
		glog.Fatal("numOfWorker should never be smaller than 1, recommend using NumCPU()\n")
	}

	for i := range builder.tasks4WorkerC {
		builder.tasks4WorkerC[i] = make(chan spatialindexer.PointInfo, 50)
	}

	return builder
}

/*
                                    ->   tasksChannel   ---> worker
                                  /                                   \
                                 /                                     \
Input Iterator ==> dispatcher    --->    tasksChannel   ---> worker      ---> aggregatorChannel -> iterate and put result to map
                                 \                                     /
                                  \                                   /
                                    ->   tasksChannel   ---> worker
*/

func (builder *connectivityMapBuilder) build() ID2NearByIDsMap {
	builder.dispatch()
	builder.process()
	builder.aggregate()
	builder.wait()

	return builder.id2NearbyIDs
}

func (builder *connectivityMapBuilder) dispatch() {
	builder.dispatchWaitGroup.Add(1)

	go func() {
		counter := 0
		for p := range builder.iterator.IteratePoints() {
			builder.tasks4WorkerC[counter%builder.numOfWorker] <- p
			counter += 1
		}

		for i := 0; i < builder.numOfWorker; i++ {
			close(builder.tasks4WorkerC[i])
		}

		glog.Infof("builder's dispatch is finished.  Total input is %d.\n", counter)
		builder.dispatchWaitGroup.Done()
	}()
}

func (builder *connectivityMapBuilder) process() {
	for i := 0; i < builder.numOfWorker; i++ {
		builder.workerWaitGroup.Add(1)
		go builder.work(builder.workerWaitGroup, i, builder.tasks4WorkerC[i], builder.aggregatorC)
	}

	glog.Infof("builder's process is finished, start number of %d workers.\n", builder.numOfWorker)
}

func (builder *connectivityMapBuilder) work(wg *sync.WaitGroup, workerID int, source <-chan spatialindexer.PointInfo, sink chan<- placeIDWithNearByPlaceIDs) {
	defer wg.Done()

	counter := 0
	for p := range source {
		counter += 1

		nearbyIDs := builder.finder.FindNearByPointIDs(p.Location, builder.distanceLimit, spatialindexer.UnlimitedCount)
		rankedResults := builder.ranker.RankPointIDsByShortestDistance(p.Location, nearbyIDs)

		ids := make([]IDAndDistance, 0, len(rankedResults))
		for _, r := range rankedResults {
			ids = append(ids, IDAndDistance{
				ID:       r.ID,
				Distance: r.Distance,
			})
		}

		sink <- placeIDWithNearByPlaceIDs{
			id:  p.ID,
			ids: ids,
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
	builder.dispatchWaitGroup.Wait()
	builder.workerWaitGroup.Wait()
	close(builder.aggregatorC)
	builder.aggregatorWaitGroup.Wait()
}

type placeIDWithNearByPlaceIDs struct {
	id  spatialindexer.PointID
	ids []IDAndDistance
}

func (builder *connectivityMapBuilder) buildInSerial() ID2NearByIDsMap {
	glog.Warning("This function is only used for compare result of worker's build().\n`")
	internalResult := make(chan placeIDWithNearByPlaceIDs, 10000)
	m := make(ID2NearByIDsMap)

	go func() {
		for p := range builder.iterator.IteratePoints() {
			nearbyIDs := builder.finder.FindNearByPointIDs(p.Location, builder.distanceLimit, spatialindexer.UnlimitedCount)
			rankedResults := builder.ranker.RankPointIDsByGreatCircleDistance(p.Location, nearbyIDs)

			ids := make([]IDAndDistance, 0, len(rankedResults))
			for _, r := range rankedResults {
				ids = append(ids, IDAndDistance{
					ID:       r.ID,
					Distance: r.Distance,
				})
			}
			internalResult <- placeIDWithNearByPlaceIDs{
				id:  p.ID,
				ids: ids,
			}
		}
		close(internalResult)
	}()

	for item := range internalResult {
		m[item.id] = item.ids
	}

	return m
}
