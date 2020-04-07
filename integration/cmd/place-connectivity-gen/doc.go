// Package main contains the tool of chargestation-connectivity generator
//
// stage 1:
// inputs is json file
// => convert to slice of [id:string,location: lat,lon]
// => calculate cellids for each point(for all levels)
// => build revese index for cellID -> ids
//
// stage 2:
// => iterate each point
// => generate a circle(s2::cap), find all cellids intersect with that circle
// => retrieve all ids
// => generate result of id(from), ids(all ids in certain distance)
//
// stage 3:
// => load data from file
// => for each line, its formid and all other ids
// => calculate distance between fromid and all other ids
// => sort result based on distance and write back to file
//
package main
