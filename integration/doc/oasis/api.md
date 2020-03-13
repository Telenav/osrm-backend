- [OSAIS API](#osais-api)
  - [Input](#input)
    - [Input Example](#input-example)
  - [Response](#response)
    - [Response Example](#response-example)
  - [Reference](#reference)

# OSAIS API

Oasis API expect user input `orig point`, `destination point` and `vehicle information`, it will select needed charge station to achieve the trip and return them as `waypoints`.


## Input 

The input contains information related with user's vehicle status, where to go and specific settings.

- Orig point
- Dest point
- Max charge range
- Current electric range
- Preferred charge buffer level
- Safe charge level
- Curve

|Option   |Values   | Descriptions   |Comments|
|:-|:-|:-|:-|
|max_range   |float(meters)   |Max range if fully charged   | |
|curr_range   |float(meters)   |Distance represent current electric level   | |
|prefer_level   |float(meters)   |Preferred charge buffer level   |[more info](https://github.com/Telenav/osrm-backend/issues/128#issuecomment-573171852)   |
|safe_level   |float(meters)   |Safe charge level   | [more info](https://github.com/Telenav/osrm-backend/issues/128#issuecomment-573171852)  |
|curve   |string   |TBD   |   |


### Input Example  

```url
/oasis/v1/earliest/-82.058695,35.036645;-81.89309,34.97914?max_range=500000.0&curr_range=160000.0
```

## Response 

Response contains information for charge station needed to complete the route.

### Response Example 

```JSON
[
{
  "distance": 90.0,
  "duration": 300.0,
  "estimate_remaining_range":100000.0,
  "weight": 300.0,
  "weight_name": "duration",
  "charge_stations": [
    {
      "address" : [
               {
                      "geo_coordinates": {
                                "latitude": 37.78509,  
                                 "longitude": -122.41988
                       },
                      "nav_coordinates": [
                                  {
                                           "latitude": 37.78509,
                                           "longitude": -122.41988
                                  }
                       ]
               }
       ],
      "wait_time" : 30.0,
      "charge_time": 100.0,
      "charge_range": 100.0,
       "detail_url":"url from search component which could retrieve charge station's information"
    },
    {
      "address" : [
               {
                      "geo_coordinates": {
                                "latitude": 13.40677,  
                                 "longitude": 52.53333
                       },
                      "nav_coordinates": [
                                  {
                                           "latitude": 13.40677,
                                           "longitude": 52.53333
                                  }
                       ]
               }
       ],
      "wait_time": 100.0,
      "charge_time": 100.0,
      "charge_range": 100.0,
      "detail_url":"url from search component which could retrieve charge station's information"
    },
  ]
}

]
```

## Reference 
- [OSRM HTTP Document](https://github.com/Telenav/osrm-backend/blob/master/docs/http.md)
- [HERE's energy consumption model](https://developer.here.com/documentation/routing/dev_guide/topics/resource-param-type-custom-consumption-details.html)
- [Issue 120](https://github.com/Telenav/osrm-backend/issues/128)
