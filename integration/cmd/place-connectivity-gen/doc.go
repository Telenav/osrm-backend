// Package main contains the tool of chargestation-connectivity generator
//
// stage 1:
// inputs is json file(Example: https://github.com/Telenav/osrm-backend/blob/2b28e5b62de7fbf35163698966ee005683b1c8a3/integration/service/spatialindexer/poiloader/sample_input.json#L1)
// => convert to slice of [id:string,location: lat,lon]
// => calculate cellids for each location(for all levels)
// => build revese index for cellID -> ids
//
// stage 2:
// => iterate each place id
// => generate a circle(s2::cap), find all cellids intersect with that circle
// => retrieve all place ids in those cellids
// => generate result of place id(from), place ids
//
// stage 3:
// => load previous result from file
// => rank place ids by calculating distance between center place id and all others
// => Write results to file
//
package main
