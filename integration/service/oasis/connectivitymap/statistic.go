package connectivitymap

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"strings"

	"github.com/Telenav/osrm-backend/integration/api"
	"github.com/golang/glog"
)

const statisticFileName = "connectivity_map_statistic.json"

type statistic struct {
	Count                 int     `json:"count"`
	ValidCount            int     `json:"valid_count"`
	AverageNearByIDsCount int     `json:"average_nearby_ids_count"`
	MaxNearByIDsCount     int     `json:"max_nearby_ids_count"`
	MinNearByIDsCount     int     `json:"min_nearby_ids_count"`
	AverageMaxDistance    float64 `json:"average_value_max_distance"`
	MaxOfMaxDistance      float64 `json:"max_value_max_distance"`
	MinOfMaxDistance      float64 `json:"min_value_max_distance"`
	MaxRange              float64 `json:"maxrange_set_by_preprocessing"`
}

func newStatistic() *statistic {
	return &statistic{
		Count:                 0,
		ValidCount:            0,
		AverageNearByIDsCount: 0,
		MaxNearByIDsCount:     math.MinInt32,
		MinNearByIDsCount:     math.MaxInt32,
		AverageMaxDistance:    0.0,
		MaxOfMaxDistance:      0.0,
		MinOfMaxDistance:      math.MaxFloat64,
		MaxRange:              0.0,
	}
}

func (s *statistic) init() {
	s.Count = 0
	s.ValidCount = 0
	s.AverageNearByIDsCount = 0
	s.MaxNearByIDsCount = math.MinInt32
	s.MinNearByIDsCount = math.MaxInt32
	s.AverageMaxDistance = 0.0
	s.MaxOfMaxDistance = 0.0
	s.MinOfMaxDistance = math.MaxFloat64
	s.MaxRange = 0.0
}

func (s *statistic) build(m ID2NearByIDsMap, MaxRange float64) *statistic {
	s.init()
	s.MaxRange = MaxRange

	totalNearByIDsCount := 0
	totalMaxDistance := 0.0
	for _, idAndDistanceArray := range m {
		s.Count += 1
		if len(idAndDistanceArray) == 0 {
			continue
		}
		s.ValidCount += 1

		prevTotalNearByIDsCount := totalNearByIDsCount
		totalNearByIDsCount += len(idAndDistanceArray)
		if prevTotalNearByIDsCount > totalNearByIDsCount {
			glog.Fatalf("Overflow during accumulate totalNearByIDsCount, before accumulate value = %v, after add %v new value is %v\n",
				prevTotalNearByIDsCount, len(idAndDistanceArray), totalNearByIDsCount)
		}
		s.MaxNearByIDsCount = max(s.MaxNearByIDsCount, len(idAndDistanceArray))
		s.MinNearByIDsCount = min(s.MinNearByIDsCount, len(idAndDistanceArray))

		prevTotalMaxDistance := totalMaxDistance

		maxDistance := 0.0
		for _, item := range idAndDistanceArray {
			maxDistance = math.Max(maxDistance, item.Distance)
		}

		totalMaxDistance += maxDistance
		if prevTotalMaxDistance > totalMaxDistance {
			glog.Fatalf("Overflow during accumulate totalMaxDistance, before accumulate value = %#v, after add %v new value is %#v\n",
				prevTotalMaxDistance, maxDistance, totalMaxDistance)
		}
		s.MaxOfMaxDistance = math.Max(s.MaxOfMaxDistance, maxDistance)
		s.MinOfMaxDistance = math.Min(s.MinOfMaxDistance, maxDistance)
	}

	if s.ValidCount == 0 {
		glog.Warningf("connectivity's statistic detect 0 valid result.\n")
		return s
	}

	s.AverageNearByIDsCount = totalNearByIDsCount / s.ValidCount
	s.AverageMaxDistance = totalMaxDistance / (float64)(s.ValidCount)

	glog.Infof("Build statistic for ID2NearByIDsMap finished. %+v\n", s)
	return s
}

func (s *statistic) dump(folderPath string) error {
	if !strings.HasSuffix(folderPath, api.Slash) {
		folderPath += api.Slash
	}

	file, err := json.Marshal(*s)
	if err != nil {
		glog.Errorf("Marshal object %+v failed with error %+v\n", s, err)
		return err
	}

	if err := ioutil.WriteFile(folderPath+statisticFileName, file, 0644); err != nil {
		glog.Errorf("Dump %s to %s failed with error %+v\n", statisticFileName, folderPath, err)
		return err
	}

	glog.Infof("Finished dumpping statistic file of %s to %s\n", statisticFileName, folderPath)
	return nil
}

func (s *statistic) load(folderPath string) error {
	if !strings.HasSuffix(folderPath, api.Slash) {
		folderPath += api.Slash
	}

	file, err := ioutil.ReadFile(folderPath + statisticFileName)
	if err != nil {
		glog.Errorf("Load %s from %s failed with error %+v\n", statisticFileName, folderPath, err)
		return err
	}

	err = json.Unmarshal([]byte(file), s)
	if err != nil {
		glog.Errorf("Unmarshal statistic file %s from %s failed with error %+v\n", statisticFileName, folderPath, err)
		return err
	}

	return nil
}

// max returns the larger of x or y.
func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// min returns the smaller of x or y.
func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
