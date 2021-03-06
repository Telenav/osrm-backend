OASIS stands for Optimized Charge Station Selection Service, which mainly supports generating optimal charge station candidates for given query based on user's vehical information. It expects routing engine(like OSRM) to generate electronic specific cost and only focus on how to choose most optimal charge station combination which balances charging cost(charging time, payment, etc).  
For example, OASIS might suggest user spend 1 hour to charge at station A for x amount of energy then reach destination with 2 hours in total, or charge at station B and station C for 30 minutes each and reach destination with duration of 2 hours 10 minutes.

Documents
- [api document](./api.md)
- [architecture design](./architecture_design.md)
- [container design](./container_design.md)
- Low level design
  + [search charge stations along route](./search_chargestations_along_route.md)
  + [charge station based routing](./charge_station_based_routing.md)
