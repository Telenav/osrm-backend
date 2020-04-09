/*

Package stationfinder provide functionality to find nearby charge stations and
related algorithm.
It's find functionality could be achieved by:
- TNSearchFinder, which is implemented by Telenav's search web service.
  + The data source is from professional data
  + Based on Lucence and hilbert value
  + Could potentially supports dynamic information

- LocalIndexerFinder
  + Support any kind of data sources as long as following format defined in https://github.com/Telenav/osrm-backend/issues/237#issue-585201041
  + For Telenav usage, uses the same source with Telenav search service
  + Spatial index is build based on google s2

Finders:
- origStationFinder holds logic for how to find reachable charge stations
  based on current energy level.
- destStationFinder holds logic for how to find reachable charge stations
  based on safe energy level and distance to nearest charge station(todo).
- lowEnergyLocationStationFinder holds logic for how to find reachable
  charge station near certain location.
- orig_iterator and dest_iterator wraps single point of orig/dest which could
  be used for algorithms

Algorithm:
- Each finder provide iterator to iterate charge station candidates.
- The choice of channel as response makes algorithm could be asynchronous func.
- FindOverlapBetweenStations provide functionality to find overlap
  between two iterator.
- CalcWeightBetweenChargeStationsPair provide functionality to calculate
  cost between two group of charge stations which could construct a new
  graph as edges.

*/
package stationfinder
