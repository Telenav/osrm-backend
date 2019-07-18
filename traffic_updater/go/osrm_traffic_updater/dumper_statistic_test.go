package main

import (
	"testing"
	"sync"
)

func TestDumperStatistic(t *testing.T) {
	var d dumperStatistic
	var wg sync.WaitGroup

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go accumulateDumper(&d, &wg)
	}
	wg.Wait()

	if (d.wayCnt != 10) || (d.nodeCnt != 20) || (d.fwdRecordCnt != 30) || (d.bwdRecordCnt != 40) || (d.wayMatchedCnt != 50) || (d.nodeMatchedCnt != 60) {
			t.Error("TestDumperStatistic failed.\n")
		}

}

func accumulateDumper(d *dumperStatistic, wg *sync.WaitGroup) {
	d.update(1,2,3,4,5,6)
	wg.Done()
}

