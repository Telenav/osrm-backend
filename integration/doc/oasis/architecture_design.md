# Architecture design

## Overview
<img src="https://user-images.githubusercontent.com/16873751/86186667-81bf4480-baee-11ea-9c41-cca538078232.png" alt="overview" width="600"/><br/>

## Service layer
<img src="https://user-images.githubusercontent.com/16873751/86185744-f1800000-baeb-11ea-8683-4cd8e38a7a4c.png" alt="overview" width="600"/><br/>

## Solution layer

<img src="https://user-images.githubusercontent.com/16873751/86185765-ffce1c00-baeb-11ea-9cf2-81c9ad9dbb5d.png" alt="overview" width="600"/><br/>

- `Solution` is the layer contains logic for how to select charge stations, such as
  * whether need charge or not
  * is destination reachable by single charge
  * multiple charge station solution finding, which could be search along route or charge station based routing, implemented by lower layer
- `GenerateSolution` abstract the interface for how to generate multiple charge station solution



## Graph layer

<img src="https://user-images.githubusercontent.com/16873751/86185777-0b214780-baec-11ea-986c-88ac16f3a2ed.png" alt="overview" width="600"/><br/>

- `stationgraph` implements the `GenerateSolution` interface, and this package represents algorithm
- `chargingstrategy` abstract logic of `charge` which supports calculation in `stationgraph`

## Place layer

<img src="https://user-images.githubusercontent.com/16873751/86185788-14121900-baec-11ea-9da1-051e61ca8c1d.png" alt="overview" width="600"/><br/>

A different view:  
<img src="https://user-images.githubusercontent.com/16873751/86186460-01004880-baee-11ea-8c1a-2d24268a002c.png" alt="overview" width="600"/><br/>

Definition of `Place`
```go
// Place records place(location, point of interest) related information such as
// ID and location
// Place represents charge stations for most of times for OASIS service, but it 
// could also represent for a user select location such as original location or
// destination location.
type Place struct {
   ID       PlaceID
   Location *nav.Location
}

// PlaceID defines ID for given place(location, point of interest)
// The data used for pre-processing must contain valid PlaceID, which means it
// either a int64 directly or be processed as int64
type PlaceID int64
```

Definition of `TopoQuerier` interface
```go
// TopoQuerier used to return topological information for places
type TopoQuerier interface {

   // GetNearByPlaces finds near by stations by given placeID and return them in recorded sequence
   // Returns nil if given placeID is not found or no connectivity
   GetNearByPlaces(placeID common.PlaceID) []*common.RankedPlaceInfo

   // LocationQuerier returns *nav.location for given placeID
   LocationQuerier
}
```

Definition of `SpatialQuerier` interface
```go
// SpatialQuerier answers special query
type SpatialQuerier interface {

   // GetNearByPlaceIDs returns a group of places near to given center location
   GetNearByPlaceIDs(center nav.Location, radius float64, limitCount int) []*common.PlaceInfo
}
```
Definition of `LocationQuerier` interface
```go
// LocationQuerier returns *nav.location for given placeID
type LocationQuerier interface {

   // GetLocation returns *nav.Location for given placeID
   // Returns nil if given placeID is not found
   GetLocation(placeID PlaceID) *nav.Location
}
```


## More info
- For discussions, please go to [#352](https://github.com/Telenav/osrm-backend/issues/352)