/*
Package stationconnquerier provides connectivitymap.Querier interface based on pre-build connectivity data.
It needs to connect orig and destination point into station graph.

For orig/start point, based on electric vehicle's current energy level, it queries all possible reachable stations take start point as center.
For destination/end point, based on electric vehicle's max energy level, it queries all possible stations which could reach destination with maximum amount of charge.
For charge stations, it retrieves connectivity from pre-build data.
*/
package stationconnquerier
