/*

Package stationgraph builds charge station based graph and provide optimum charge solutions for given orig/dest.
- To create station graph
    + currEnergyLevel: current energy level for given electrical vehicle
    + maxEnergyLevel: maximum energy capacity for given electrical vehicle
    + strategy: Charge strategy be used for each charge station, it hides logic of different strategy of {time used for charging/wating, new energy got}.
    + querier: topological querier which generates connectivity for building graph
         * querier could be implemented by querying OSRM or pre-build conductivities
         * querier must contains stationfindertype.OrigLocationIDStr and stationfindertype.DestLocationIDStr
         * querier must provide ways to integrate stationfindertype.OrigLocationIDStr and stationfindertype.DestLocationIDStr into connectivity graph

- station graph returns
    + solution.Solution via GenerateChargeSolutions()
    + station graph will call algorithm and built certain type of Graph to calculate result

Data structure
- Graph defines interface of a graph.
- nodeGraph and mockGraph implements interface of graph
- For mockGraph, it mainly used for algorithm testing.  For all the nodes, topological is hard coded at initial time.
- For nodeGraph, which is built topological on the fly during running time.
    + Its start and node need to be setted.
    + physical node: 1-to-1 match to charge stations and orig/dest point.  For example, the following graph has 7 physical nodes
    + logical node: for each charge station, we might exist charge station with different charge status, depend on time used in charge stations.  For example, the following graph has 17 logical nodes.
         * given charge station 1, we have different choices
         * station 1, charge to 60% of maximum energy then left
         * station 1, charge to 80% of maximum energy then left
         * station 1, charge to 100% of maximum energy then left
         * Logical node is identified by unique station ID + chargingStatus
    + connectivity: first built physical connectivity, then built logical
         * physical connectivity is ansowered by interface of connectivitymap.Querier, eg, station 1 connects station 4 and station 5
         * logical connectivity is build considering physical connection + final charging status
         * logical connectivity's result is the final graph which algorithm runs on


                    station 1
               /       \    \
              /         \    \
             /          _\_____station 4
            /          /  \     /        \
           /          /    \   /          \
start  -------   station 2  \ /           end
           \          \     / \          /
            \          \   /   \        /
             \          \_/_____station 5
              \          /     /
               \        /     /
                    station 3
more information of this example graph could go to testion_graph_test


Algorithm
- dijkstra implements single direction dijkstra based on Graph's interface
- Basically, classic graph algorithm could find path with minimum cost, several modifications with electric vehicle case:
    + Whether two charge station is reachable or not.
         * For example, we have a node charge to 60% of total energy, let's say maximum is 100 unit, so this node could reach other nodes which is reachable by 60 energy unit.
    + Arrival energy.  For each node, we have a energy status when left previous charge station and will spend energy for reaching current node, arrival energy is calculated by that.
    + Charging time needed.  Charging time depends on two things: arrival energy and the status of charging when you left
         * For example, when reaching charge station, you have energy of 30, and you want to get 60% of maximum energy when left charge station for given node
         * chargingstrategy.Strategy will evaluate the time needed for charging, which hides the logic of charging curve for specific vehicle.

*/
package stationgraph
