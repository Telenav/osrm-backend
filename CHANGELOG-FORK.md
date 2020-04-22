
# Unreleased
Changes from v10.2.0      
- Feature:    
  - CHANGED: for internal refactoring, moved `unidbpatch` and `mapsource` packages into `integration/util` folder [#300](https://github.com/Telenav/osrm-backend/pull/300)
- Bugfix:    
- Performance:    
- Tools:    
- Docs:    


# v10.2.0
Changes from v10.1.0      
- Feature:    
  - ADDED ranker to rank near by places for `oasis` [#258](https://github.com/Telenav/osrm-backend/pull/258)
  - ADDED connectivity map for `oasis` [#259](https://github.com/Telenav/osrm-backend/pull/259)
  - ADDED `place-connectivity-generator` command line tool to generate connectivity map for places for `oasis` [#268](https://github.com/Telenav/osrm-backend/pull/268)
  - ADDED local station finder based on `google:s2` indexer for `oasis` [#278](https://github.com/Telenav/osrm-backend/pull/278)
  - CHANGED for internal refactoring, refactor station finder to support multiple strategy for `oasis` [#271](https://github.com/Telenav/osrm-backend/pull/271)
  - CHANGED for internal refactoring, use `struct embedding` to avoid interface forwarding in `oasis` [#280](https://github.com/Telenav/osrm-backend/pull/280)
  - ADDED historical Speed query interface on a specified time for way, for `osrm-ranking` [#256](https://github.com/Telenav/osrm-backend/pull/256)
  - CHANGED: for internal refactoring, move ranking service into service folder [#261](https://github.com/Telenav/osrm-backend/pull/261)
  - CHANGED: for internal refactoring, move live traffic related packges into `integration/traffic` [#264](https://github.com/Telenav/osrm-backend/pull/264)
  - ADDED new cmd tool `nodes2way-builder` to generate the db from `wayid2nodeids.csv` or `wayid2nodeids.csv.snappy` [#274](https://github.com/Telenav/osrm-backend/pull/274)
  - ADDED new cmd tool `nodes2way-cli` to able to query ways from nodes [#274](https://github.com/Telenav/osrm-backend/pull/274)
  - CHANGED for internal refactoring, move `pkg/wayidsflag` to `util/idsflag` to extend its scope [#277](https://github.com/Telenav/osrm-backend/pull/277)
  - CHANGED **HTTP API** `JSON` response `annotation/nodes` from `Number` to `string` to avoid conversion overflow [#285](https://github.com/Telenav/osrm-backend/pull/285)    
- Bugfix:    
- Performance:    
  - CHANGED `nodes2way-builder` to improve the db building performance [#281](https://github.com/Telenav/osrm-backend/pull/281)
- Tools:    
  - CHANGED pull request template for easy using [#283](https://github.com/Telenav/osrm-backend/pull/283)     
  - ADDED `nodes2way.db` build process when compile data [#286](https://github.com/Telenav/osrm-backend/pull/286)     
  - ADDED `nodes2way.db.snappy` into `osrm-backend-within-mapdata` docker image [#286](https://github.com/Telenav/osrm-backend/pull/286)    
- Docs:    
  - ADDED CHANGELOG-FORK.md for keeping changelogs of this forked repo [#283](https://github.com/Telenav/osrm-backend/pull/283)

# v10.1.0
Init release. Changes from Init fork(0a556fe45073e3c01d4ce90911017421c71129b3)
- Feature:    
  - ADDED `osrm-traffic-updater` command line tool for live traffic integration via GRPC [#18](https://github.com/Telenav/osrm-backend/pull/18)
  - ADDED `wayid2nodeid-extractor` command line tool to extract `wayid,nodeid,nodeid,...` mapping from PBF [#28](https://github.com/Telenav/osrm-backend/pull/28)    
  - ADDED cell update ratio and only customize cell be updated [#62](https://github.com/Telenav/osrm-backend/pull/62)
  - ADDED `osrm-ranking` microservice to post process response with traffic [#98](https://github.com/Telenav/osrm-backend/pull/98)
  - REMOVED debug log compile time control, then it's able to output debug log only by `-l DEBUG` on runtime [#113](https://github.com/Telenav/osrm-backend/pull/113)     
  - ADDED new car profile and libs for unidb by convention `xxx-unidb` [#214](https://github.com/Telenav/osrm-backend/pull/214)
  - ADDED `historicalspeed-timezone-builder` command line tool to build timezone information for historical speeds [#243](https://github.com/Telenav/osrm-backend/pull/243)
  - ADDED historical speed files loading in `osrm-ranking` [#235](https://github.com/Telenav/osrm-backend/pull/235)
  - ADDED `oasis` microservice for EV routing [#130](https://github.com/Telenav/osrm-backend/pull/130)
  - ADDED `poi-converter` command line tool to preprecess POI for `oasis` [#247](https://github.com/Telenav/osrm-backend/pull/247)
  - ADDED spatial index based on google::s2 for `oasis` [#248](https://github.com/Telenav/osrm-backend/pull/248)
- Bugfix:    
  - CHANGED `OSMNodeID` from `33bits` to `63bit` to compatible for PBF of Telenav UniDB [#15](https://github.com/Telenav/osrm-backend/pull/15)
  - ADDED `speed_unit` tag parsing to compatible for PBF of Telenav UniDB [#67](https://github.com/Telenav/osrm-backend/pull/67)
  - FIXED traffic signals parsing for UniDB PBF [#232](https://github.com/Telenav/osrm-backend/pull/232)
  - FIXED wayid overflow when UniDB PBF [#112](https://github.com/Telenav/osrm-backend/pull/112)
- Performance:    
- Tools:    
  - ADDED `docker-orchestration` folder for dockers and kubernetes setup and integration        
  - ADDED multiple options to startup `osrm-routed` when launch container: `routed_no_traffic_startup`, `routed_startup` and `routed_blocking_traffic_startup`    
  - ADDED `compile_mapdata` that able to compile mapdata in contianer      
  - ADDED compatible to compile mapdata from OSM or UniDB  
  - ADDED `trafficproxy-cli` command line tool for debugging live traffic easier [#87](https://github.com/Telenav/osrm-backend/pull/87)
  - ADDED `osrm-files-extractor` command line tool that able to parse OSRM compile files directly     
  - ADDED CI/CD based on GitHub Actions
- Docs:    
  - ADDED `osrm-ranking` design [#78](https://github.com/Telenav/osrm-backend/pull/78)
  - FIXED wrong comment for compressing edges if it crosses a traffic signal [#199](https://github.com/Telenav/osrm-backend/pull/199)

