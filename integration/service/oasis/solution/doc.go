/*
Package solution implements the interface of Generator
 #1. checks whether current energy is sufficient.  Its possible that user could
 directly arrive destination under his energy capacity.

 #2. checks whether could reach destination by one time charge.  The logic for
 calculating one time charge or multiple times charge could be the same.  But
 here we choose to implement specific logic for calculating single charge.
 For original point, with its current capacity, there is a certain range and
 covers x number of charge stations:
 https://user-images.githubusercontent.com/16873751/73227073-1b8d5a80-4127-11ea-81b2-b5cdadcbfff9.png
 And the same for destination, we could draw similar circle but with maximum
 energy capacity of the vehicle:
 https://user-images.githubusercontent.com/16873751/73227100-3790fc00-4127-11ea-813d-80472725bf71.png
 If there is any shared charge staions in both ranges, then we could find most
 optimal solution from them for charge one times

 #3 finds solution for optimal multiple charges based on graph layer
*/
package solution
