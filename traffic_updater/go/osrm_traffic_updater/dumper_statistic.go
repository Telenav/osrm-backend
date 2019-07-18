package main

import (
	"sync"
	"fmt"
)

type dumperStatistic struct {
	wayCnt uint64
	nodeCnt uint64
	fwdRecordCnt uint64
	bwdRecordCnt uint64
	wayMatchedCnt uint64
	nodeMatchedCnt uint64
	mu sync.Mutex
}

func (d* dumperStatistic) update(wayCnt uint64, nodeCnt uint64, fwdRecordCnt uint64, 
	bwdRecordCnt uint64, wayMatchedCnt uint64, nodeMatchedCnt uint64) {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	d.wayCnt += wayCnt
	d.nodeCnt += nodeCnt
	d.fwdRecordCnt += fwdRecordCnt
	d.bwdRecordCnt += bwdRecordCnt
	d.wayMatchedCnt += wayMatchedCnt
	d.nodeMatchedCnt += nodeMatchedCnt
}

func (d* dumperStatistic) output() {
	d.mu.Lock()
	defer d.mu.Unlock()

	fmt.Printf("Statistic: \n")
	fmt.Printf("Load %d way from data with %d nodes.\n", d.wayCnt, d.nodeCnt)
	fmt.Printf("%d way with %d nodes matched with traffic record.\n", 
		d.wayMatchedCnt, d.nodeMatchedCnt)
	fmt.Printf("Generate %d records in final result with %d of them from forward traffic and %d from backword.\n", 
		d.fwdRecordCnt+ d.bwdRecordCnt, d.fwdRecordCnt, d.bwdRecordCnt)
}

