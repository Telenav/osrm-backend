/*
package chargingstrategy contains EV related domain knowledge.
It will handle:
- Different vehicle has different charging curve and battery capacity
- Different charging stations has different amount of chargers and different type(L2, L3)

It will return(For initial version):
- Charging candidates which could represent time used in charge station and additional energy got.
  Waiting time, cost could also be added here.
  Charging candidates will be converted into graph nodes which could be applied for different kind of algorithm,
  such as find best charging strategy
*/
package chargingstrategy
