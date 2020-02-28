package stationgraph

/*
- Node represents charge station and wraps energy operation.  It could dected whether
  a target node is reachable or not, or how far it could go after charging.
- Neighbor records information from current node to oppsite one, such as distance,
  duration, energy consumption, etc.
- Graph represents charge station based node graph and provide functionality to find
  best solution.  Graph will be constructed by sub-class, once it has been built then
  orig and dest has been decided.  When sub-class construct the graph, it will crate
  information about node, and each node will contains neighbors which could be explored.
*/
