// Package spatialindexer answers query of nearest places(place, point of interest) for conditions
// such as center location, radius, etc
//
// Sample Scenario 1: Build connectivity for charge stations during pro-processing
// indexer := NewS2Indexer().Build(poiCsvFile)
// for _, stationPoint := range chargeStations {
// 	nearbyStations := indexer.FindNearByIDs(stationPoint, 800km, -1)
// 	rankedStations := indexer.RankingIDsByShortestDistance(stationPoint, nearbyStations)
// }
//
//
// Sample Scenario 2: Dump S2Indexer's content to folder
// indexer.Dump(folderPath)
//
//
// Sample Scenario 3: Query reachable charge stations with current energy level
// indexer := NewS2Indexer().Load(folderPath)
// nearbyStations := indexer.FindNearByIDs(currentPoint, currentEnergyLevel, -1)
//
package spatialindexer
