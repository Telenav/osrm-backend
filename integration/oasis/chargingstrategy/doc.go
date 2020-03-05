/*
package chargingstrategy contains domain knowledge related with electric vehicle charging.
It will consider:
- Different vehicle has different charging curve and battery capacity
- Different charging stations has different amount and type((L2, L3)) of chargers

It will return(For initial version):
- Charging candidates which could represent time used in charge station and additional energy got.
  Waiting time, cost could also be added here.
  Charging candidates will be converted into graph nodes which could be applied for different kind of algorithms,
  such as find best charging strategy
*/
package chargingstrategy
