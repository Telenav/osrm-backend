package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	proxy "github.com/Telenav/osrm-backend/traffic_updater/pkg/gen-trafficproxy"
)

func saveTrafficDataFromGRPC(targetPath string, trafficData proxy.TrafficResponse) {
	startTime := time.Now()
	defer func() {
		log.Printf("saveTrafficDataFromGRPC to file %s takes %f seconds\n", targetPath, time.Now().Sub(startTime).Seconds())
	}()

	if err := os.RemoveAll(targetPath); err != nil {
		log.Fatal(err)
		return
	}

	outfile, err := os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE, 0755)
	defer outfile.Close()
	defer outfile.Sync()
	if err != nil {
		log.Fatal(err)
		log.Printf("Open output file of %s failed.\n", targetPath)
		return
	}
	log.Printf("Open output file of %s succeed.\n", targetPath)

	w := bufio.NewWriter(outfile)
	defer w.Flush()
	for _, flow := range trafficData.FlowResponses {
		_, err := w.WriteString(flow.String() + "\n")
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	for _, incident := range trafficData.IncidentResponses {
		_, err := w.WriteString(incident.String() + "\n")
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

func TestGetAllTrafficDataByGRPC(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	trafficData, err := getTrafficFlowsIncidentsByGRPC(flags.trafficProxyFlags, nil)
	if err != nil {
		t.Error(err)
	}
	quickViewFlows(trafficData.FlowResponses, 10)         //quick view first 10 lines
	quickViewIncidents(trafficData.IncidentResponses, 10) //quick view first 10 lines

	saveTrafficDataFromGRPC("dump_alltrafficdata", *trafficData)
}

func TestGetTrafficDataForWaysByGRPC(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	var wayIds []int64
	wayIds = append(wayIds, 829733412, 104489539)

	trafficData, err := getTrafficFlowsIncidentsByGRPC(flags.trafficProxyFlags, wayIds)
	if err != nil {
		t.Error(err)
	}
	quickViewFlows(trafficData.FlowResponses, 10)         //quick view first 10 lines
	quickViewIncidents(trafficData.IncidentResponses, 10) //quick view first 10 lines
}

func TestGetDeltaTrafficDataByGRPCStreaming(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	trafficDataChan := make(chan proxy.TrafficResponse)

	go func() {
		err := getDeltaTrafficFlowsIncidentsByGRPCStreaming(flags.trafficProxyFlags, trafficDataChan)
		if err != nil {
			t.Errorf("getDeltaTrafficFlowsIncidentsByGRPCStreaming failed, err: %v", err)
		}
	}()

	startTime := time.Now()
	statisticsInterval := 120 //120 seconds
	intervalIndex := 0
	var currentIntervalTrafficData proxy.TrafficResponse
	var totalFlowsCount, maxFlowsCount, minFlowsCount int64
	var totalIncidentsCount, maxIncidentsCount, minIncidentsCount int64
	var recvCount int
	for trafficData := range trafficDataChan {
		recvCount++

		currFlowsCount := int64(len(trafficData.FlowResponses))
		totalFlowsCount += currFlowsCount
		if currFlowsCount > maxFlowsCount {
			maxFlowsCount = currFlowsCount
		}
		if minFlowsCount == 0 || currFlowsCount < minFlowsCount {
			minFlowsCount = currFlowsCount
		}

		currIncidentsCount := int64(len(trafficData.IncidentResponses))
		totalIncidentsCount += currIncidentsCount
		if currIncidentsCount > maxIncidentsCount {
			maxIncidentsCount = currIncidentsCount
		}
		if minIncidentsCount == 0 || currIncidentsCount < minIncidentsCount {
			minIncidentsCount = currIncidentsCount
		}

		currentIntervalTrafficData.FlowResponses = append(currentIntervalTrafficData.FlowResponses, trafficData.FlowResponses...)
		currentIntervalTrafficData.IncidentResponses = append(currentIntervalTrafficData.IncidentResponses, trafficData.IncidentResponses...)

		if time.Now().Sub(startTime).Seconds() >= float64(statisticsInterval) {
			log.Printf("interval %d received flows from grpc streaming in %f seconds, recv count %d, total got flows count: %d, max per recv: %d, min per recv: %d\n",
				intervalIndex, time.Now().Sub(startTime).Seconds(), recvCount, totalFlowsCount, maxFlowsCount, minFlowsCount)
			log.Printf("interval %d received incidents from grpc streaming in %f seconds, recv count %d, total got incidents count: %d, max per recv: %d, min per recv: %d\n",
				intervalIndex, time.Now().Sub(startTime).Seconds(), recvCount, totalIncidentsCount, maxIncidentsCount, minIncidentsCount)

			recvCount = 0
			totalFlowsCount = 0
			maxFlowsCount = 0
			minFlowsCount = 0
			totalIncidentsCount = 0
			maxIncidentsCount = 0
			minIncidentsCount = 0
			startTime = time.Now()

			saveTrafficDataFromGRPC("dump_deltatrafficdata_"+strconv.Itoa(intervalIndex), currentIntervalTrafficData)

			intervalIndex++
			currentIntervalTrafficData.FlowResponses = nil
			currentIntervalTrafficData.IncidentResponses = nil
		}
	}
}
