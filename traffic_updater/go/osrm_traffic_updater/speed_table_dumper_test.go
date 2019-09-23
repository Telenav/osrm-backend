package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestSpeedTableDumper1(t *testing.T) {
	// load result into sources
	var sources [TASKNUM]chan string
	for i := range sources {
		sources[i] = make(chan string, 10000)
	}
	go loadWay2NodeidsTable("./testdata/id-mapping.csv.snappy", sources)

	// construct mock traffic
	wayid2speed := make(map[int64]int)
	loadMockTrafficFlow2Map(wayid2speed)

	var ds dumperStatistic
	ds.Init(TASKNUM)
	dumpSpeedTable4Customize(wayid2speed, sources, "./testdata/target.csv", &ds)

	compareFileContentUnstable("./testdata/target.csv", "./testdata/expect.csv", t)
	validateStatistic(&ds, t)
}

func TestGenerateSingleRecord1(t *testing.T) {
	str := generateSingleRecord(12345, 54321, 33, true)
	if strings.Compare(str, "12345,54321,33\n") != 0 {
		t.Error("Test GenerateSingleRecord failed.\n")
	}
}

func TestGenerateSingleRecord2(t *testing.T) {
	str := generateSingleRecord(12345, 54321, 33, false)
	if strings.Compare(str, "54321,12345,33\n") != 0 {
		t.Error("Test GenerateSingleRecord failed.\n")
	}
}

func validateStatistic(ds *dumperStatistic, t *testing.T) {
	sum := ds.Sum()
	if (sum.wayCnt != 4) || (sum.nodeCnt != 9) || (sum.fwdRecordCnt != 4) || (sum.bwdRecordCnt != 3) || (sum.wayMatchedCnt != 4) || (sum.nodeMatchedCnt != 9) {
		t.Error("TestLoadWay2Nodeids failed with incorrect statistic.\n")
	}
}

func loadMockTrafficFlow2Map(wayid2speed map[int64]int) {
	wayid2speed[24418325] = 81
	wayid2speed[-24418332] = 87
	wayid2speed[24418332] = 87
	wayid2speed[24418343] = 47
	wayid2speed[-24418344] = 59
}

type tNodePair struct {
	f, t uint64
}

func loadSpeedCsv(f string, m map[tNodePair]int) {
	// load mock traffic file
	mockfile, err := os.Open(f)
	defer mockfile.Close()
	if err != nil {
		log.Fatal(err)
		fmt.Printf("Open file of %v failed.\n", f)
		return
	}
	fmt.Printf("Open file of %s succeed.\n", f)

	csvr := csv.NewReader(mockfile)
	for {
		row, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			} else {
				fmt.Printf("Error during decoding file %s, err = %v\n", f, err)
				return
			}
		}

		var from, to uint64
		var speed int
		if from, err = strconv.ParseUint(row[0], 10, 64); err != nil {
			fmt.Printf("#Error during decoding, row = %v\n", row)
			return
		}
		if to, err = strconv.ParseUint(row[1], 10, 64); err != nil {
			fmt.Printf("#Error during decoding, row = %v\n", row)
			return
		}
		if speed, err = strconv.Atoi(row[2]); err != nil {
			fmt.Printf("#Error during decoding, row = %v\n", row)
			return
		}

		m[tNodePair{from, to}] = speed
	}
}

func compareFileContentStable(f1, f2 string, t *testing.T) {
	b1, err1 := ioutil.ReadFile(f1)
	if err1 != nil {
		fmt.Print(err1)
	}
	str1 := string(b1)

	b2, err2 := ioutil.ReadFile(f2)
	if err2 != nil {
		fmt.Print(err2)
	}
	str2 := string(b2)

	if strings.Compare(str1, str2) != 0 {
		t.Error("Compare file content failed\n")
	}
}

func compareFileContentUnstable(f1, f2 string, t *testing.T) {
	r1 := make(map[tNodePair]int)
	loadSpeedCsv(f1, r1)

	r2 := make(map[tNodePair]int)
	loadSpeedCsv(f2, r2)

	eq := reflect.DeepEqual(r1, r2)
	if !eq {
		t.Error("TestLoadWay2Nodeids failed to generate correct map\n")
	}
}
