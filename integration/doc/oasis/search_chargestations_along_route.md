# Search charge stations along route

The flow of calculate optimum charge stations via search along route.


### 1. Pass in destination and user's current energy level

###  2. Calculate route for given origin point and dest point
![image](https://user-images.githubusercontent.com/16873751/73202948-fb42a900-40f0-11ea-91cc-4b8957b3cadc.png)


### 3. Pick segment gaps from route which is used for energy search
![image](https://user-images.githubusercontent.com/16873751/73202975-0564a780-40f1-11ea-842a-f3abba5dba38.png)


###  4. Query for search service for charge points

![image](https://user-images.githubusercontent.com/16873751/73202990-0bf31f00-40f1-11ea-9173-9487e0070161.png)


###  5. Query for table service for all distance&duration in each group.  Ie, orig -> all charge points in the first segment gaps., all charge points in the first segment gaps to second gaps, ... , to destination

![image](https://user-images.githubusercontent.com/16873751/73203010-17464a80-40f1-11ea-9e73-38ec896e92d5.png)

- Need to guarantee each range at least has one candidate has energy to reach next stage.
- Need to be aware that single charge station could exists in multiple stage's candidate.

###  6. **A high level graph is generated**, use Dynamic Programming to pick the best strategy

- Please note each charge station could have multiple charge time <-> energy level pair.

### 7. Generate result and make charge station candidates as way points.


## Architecture
![image](https://user-images.githubusercontent.com/16873751/73203274-9471bf80-40f1-11ea-93e6-e2c07dae03f3.png)



## Corner cases

### Parameter input
- Origin/Dest is not reachable
    + Origin/Dest don't have route in between
    + Origin/Dest's distance is large than POC's estimation
- Unreasonable initial energy level

### Communication with OSRM
- OSRM server is down
- OSRM table service couldn't calculate route for given [O, D] pairs.
   + For example, given one point is never reachable, what will happen

### Communication with Search
- Search service didn't response
- Search service respond very few charge stations

### Integration
- Charge stations for this step is never reachable from previous step.