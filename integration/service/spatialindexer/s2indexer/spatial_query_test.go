package s2indexer

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/service/spatialindexer"
	"github.com/golang/geo/s2"
	"github.com/golang/glog"
)

func TestSpatialIndexQuery1(t *testing.T) {
	fakeIndexer1 := S2Indexer{
		cellID2PointIDs: map[s2.CellID][]spatialindexer.PointID{
			9263834064756932608: {1, 2}, // 4/0010133
			9263851656942977024: {1, 2}, // 4/00101332
			9263847258896465920: {1, 2}, // 4/001013321
			9263843960361582592: {1, 2}, // 4/0010133210
			9263844784995303424: {1, 2}, // 4/00101332103
			9263844716275826688: {2},    // 4/001013321031
			9263844733455695872: {2},    // 4/0010133210312
			9263844746340597760: {2},    // 4/00101332103123
			9263844749561823232: {2},    // 4/001013321031233
			9263844749830258688: {2},    // 4/0010133210312332
			9263844750031585280: {2},    // 4/00101332103123323
			9263844749981253632: {2},    // 4/001013321031233230
			9263844749993836544: {2},    // 4/0010133210312332303
			9263844749996982272: {2},    // 4/00101332103123323033
			9263844749996720128: {2},    // 4/001013321031233230331
			9263844749996916736: {2},    // 4/0010133210312332303313
			9263844749996900352: {2},    // 4/00101332103123323033131
			9263844749996896256: {2},    // 4/001013321031233230331311
			9263844749996893184: {2},    // 4/0010133210312332303313110
			9263844749996892416: {2},    // 4/00101332103123323033131100
			9263844749996892608: {2},    // 4/001013321031233230331311003
			9263844749996892592: {2},    // 4/0010133210312332303313110031
			9263844749996892588: {2},    // 4/00101332103123323033131100311
			9263844749996892591: {2},    // 4/001013321031233230331311003113
			9263844853714780160: {1},    // 4/001013321032
			9263844836534910976: {1},    // 4/0010133210321
			9263844823650009088: {1},    // 4/00101332103210
			9263844822576267264: {1},    // 4/001013321032101
			9263844823381573632: {1},    // 4/0010133210321013
			9263844823180247040: {1},    // 4/00101332103210130
			9263844823129915392: {1},    // 4/001013321032101300
			9263844823134109696: {1},    // 4/0010133210321013002
			9263844823130963968: {1},    // 4/00101332103210130020
			9263844823130701824: {1},    // 4/001013321032101300201
			9263844823130898432: {1},    // 4/0010133210321013002013
			9263844823130947584: {1},    // 4/00101332103210130020133
			9263844823130943488: {1},    // 4/001013321032101300201331
			9263844823130944512: {1},    // 4/0010133210321013002013312
			9263844823130943744: {1},    // 4/00101332103210130020133120
			9263844823130943680: {1},    // 4/001013321032101300201331201
			9263844823130943696: {1},    // 4/0010133210321013002013312012
			9263844823130943708: {1},    // 4/00101332103210130020133120123
			9263844823130943709: {1},    // 4/001013321032101300201331201232
		},
		pointID2Location: map[spatialindexer.PointID]spatialindexer.Location{
			1: spatialindexer.Location{
				Lat: 37.402701,
				Lon: -121.974096,
			},
			2: spatialindexer.Location{
				Lat: 37.403530,
				Lon: -121.969768,
			},
		},
	}

	// center in 4655 great america pkwy
	center := spatialindexer.Location{
		Lat: 37.402799,
		Lon: -121.969861,
	}

	expect := []spatialindexer.PointInfo{
		spatialindexer.PointInfo{
			ID: 1,
			Location: spatialindexer.Location{
				Lat: 37.402701,
				Lon: -121.974096,
			},
		},
		spatialindexer.PointInfo{
			ID: 2,
			Location: spatialindexer.Location{
				Lat: 37.40353,
				Lon: -121.969768,
			},
		},
	}

	actual := queryNearByPoints(&fakeIndexer1, center, 10000)

	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("Expect result is \n%v but got \n%v\n", actual, expect)
	}
}

//More information could go to here: https://github.com/Telenav/osrm-backend/issues/236#issuecomment-603533484
func TestQueryNearByS2Cells1(t *testing.T) {
	// center in 4655 great america pkwy
	center := spatialindexer.Location{
		Lat: 37.402799,
		Lon: -121.969861,
	}
	actualCellIDs := queryNearByS2Cells(center, 1600)
	glog.Infof("\nTest URL is %s\n", generateDebugURL(actualCellIDs))

	expectCellIDs := []s2.CellID{
		9263843046942834688, // token = 808fc82b54
		9263843047026720768, // token = 808fc82b59
		9263843047131578368, // token = 808fc82b5f4
		9263843050566713344, // token = 808fc82c2c
		9263843050700931072, // token = 808fc82c34
		9263843050835148800, // token = 808fc82c3c
		9263843051170693120, // token = 808fc82c5
		9263843051506237440, // token = 808fc82c64
		9263843051590123520, // token = 808fc82c69
		9263843053049741312, // token = 808fc82cc
		9263843055197224960, // token = 808fc82d4
		9263843057344708608, // token = 808fc82dc
		9263843058686885888, // token = 808fc82e1
		9263843060364607488, // token = 808fc82e74
		9263843060448493568, // token = 808fc82e79
		9263843072377094144, // token = 808fc8314
		9263843073719271424, // token = 808fc8319
		9263843074256142336, // token = 808fc831b
		9263843074591686656, // token = 808fc831c4
		9263843074994339840, // token = 808fc831dc
		9263843075329884160, // token = 808fc831f
		9263843079893286912, // token = 808fc833
		9263843088483221504, // token = 808fc835
		9263843093046624256, // token = 808fc8361
		9263843093365391360, // token = 808fc83623
		9263843093784821760, // token = 808fc8363c
		9263843093919039488, // token = 808fc83644
		9263843094321692672, // token = 808fc8365c
		9263843094657236992, // token = 808fc8367
		9263843095999414272, // token = 808fc836c
		9263843176525856768, // token = 808fc8497fc
		9263843176798486528, // token = 808fc8499
		9263843178409099264, // token = 808fc849f
		9263843182972502016, // token = 808fc84b
		9263843188810973184, // token = 808fc84c5c
		9263843189146517504, // token = 808fc84c7
		9263843190488694784, // token = 808fc84cc
		9263843192032198656, // token = 808fc84d1c
		9263843192367742976, // token = 808fc84d3
		9263843192904613888, // token = 808fc84d5
		9263844597090484224, // token = 808fc9944
		9263844598432661504, // token = 808fc9949
		9263844599036641280, // token = 808fc994b4
		9263844599170859008, // token = 808fc994bc
		9263844599506403328, // token = 808fc994d
		9263844600043274240, // token = 808fc994f
		9263844601385451520, // token = 808fc9954
		9263844603532935168, // token = 808fc995c
		9263844608901644288, // token = 808fc997
		9263844614270353408, // token = 808fc9984
		9263844615612530688, // token = 808fc9989
		9263844616149401600, // token = 808fc998b
		9263844616484945920, // token = 808fc998c4
		9263844617223143424, // token = 808fc998f
		9263844617558687744, // token = 808fc99904
		9263844620444368896, // token = 808fc999b
		9263844620981239808, // token = 808fc999d
		9263844621316784128, // token = 808fc999e4
		9263844621451001856, // token = 808fc999ec
		9263844653193494528, // token = 808fc9a15
		9263844653730365440, // token = 808fc9a17
		9263844654267236352, // token = 808fc9a19
		9263844655944957952, // token = 808fc9a1f4
		9263844656028844032, // token = 808fc9a1f9
		9263844660441251840, // token = 808fc9a3
		9263844669031186432, // token = 808fc9a5
		9263844673594589184, // token = 808fc9a61
		9263844675205201920, // token = 808fc9a67
		9263844675742072832, // token = 808fc9a69
		9263844676278943744, // token = 808fc9a6b
		9263844684801769472, // token = 808fc9a8ac
		9263844684935987200, // token = 808fc9a8b4
		9263844685316620288, // token = 808fc9a8cab
		9263844699901263872, // token = 808fc9ac3
		9263844703592251392, // token = 808fc9ad0c
		9263844704934428672, // token = 808fc9ad5c
		9263844705269972992, // token = 808fc9ad7
		9263844706612150272, // token = 808fc9adc
		9263844708759633920, // token = 808fc9ae4
		9263844710101811200, // token = 808fc9ae9
		9263844710840008704, // token = 808fc9aebc
		9263844711175553024, // token = 808fc9aed
		9263844711712423936, // token = 808fc9aef
		9263844713054601216, // token = 808fc9af4
		9263844715202084864, // token = 808fc9afc
		9263844733455695872, // token = 808fc9b4
		9263844767815434240, // token = 808fc9bc
		9263844789290270720, // token = 808fc9c1
		9263844794658979840, // token = 808fc9c24
		9263844795934048256, // token = 808fc9c28c
		9263844796851552256, // token = 808fc9c2c2b
		9263844797276225536, // token = 808fc9c2dc
		9263844797611769856, // token = 808fc9c2f
		9263844798953947136, // token = 808fc9c34
		9263844801101430784, // token = 808fc9c3c
		9263844806470139904, // token = 808fc9c5
		9263844815060074496, // token = 808fc9c7
		9263844836534910976, // token = 808fc9cc
		9263844858009747456, // token = 808fc9d1
		9263844866599682048, // token = 808fc9d3
		9263844871163084800, // token = 808fc9d41
		9263844872773697536, // token = 808fc9d47
		9263844873310568448, // token = 808fc9d49
		9263844873646112768, // token = 808fc9d4a4
		9263844873780330496, // token = 808fc9d4ac
		9263844882437373952, // token = 808fc9d6b
		9263844882907136000, // token = 808fc9d6cc
		9263844883041353728, // token = 808fc9d6d4
		9263844884517748736, // token = 808fc9d72c
		9263844896865779712, // token = 808fc9da0c
		9263844896999997440, // token = 808fc9da14
		9263844898006630400, // token = 808fc9da5
		9263844898543501312, // token = 808fc9da7
		9263844899885678592, // token = 808fc9dac
		9263844902033162240, // token = 808fc9db4
		9263844903710883840, // token = 808fc9dba4
		9263844903845101568, // token = 808fc9dbac
		9263844907921965056, // token = 808fc9dc9f
		9263844908207177728, // token = 808fc9dcb
		9263844908676939776, // token = 808fc9dccc
		9263844964779950080, // token = 808fc9e9dc
		9263844965115494400, // token = 808fc9e9f
		9263844966457671680, // token = 808fc9ea4
		9263844968605155328, // token = 808fc9eac
		9263844969947332608, // token = 808fc9eb1
		9263844970484203520, // token = 808fc9eb3
		9263844970937188352, // token = 808fc9eb4b
		9263844972564578304, // token = 808fc9ebac
	}

	if !reflect.DeepEqual(actualCellIDs, expectCellIDs) {
		t.Errorf("Expect result is \n%v but got \n%v\n", expectCellIDs, actualCellIDs)
	}
}
