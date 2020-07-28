/*
Package topograph provides connectivity information for place graphs.

For example, given following graph

  -   -   -   - 4
|   |   |   |   |
  - 2 -   -   -
|   |   |   |   |
  -   -   - 3 -
|   |   |   |   |
1 -   -   -   -

Place connectivity map will pre-process each points and generate following result:
While query connectivity for place 1, will return
     (place 2, 3),   //  the shortest path between place 1 and place 2 is 3
     (place 3, 4),   //  the shortest path between place 1 and place 2 is 4
     (place 4, 7),   //  the shortest path between place 1 and place 2 is 7
The result is sorted by shortest path distance(or other user defined strategy)

When query for connectivity, user could also pass in limitation option, such as distance limitation.
For example, when query connectivity for place 3
With limitation = -1, it will return
     (place 2, 3),   //  the shortest path between place 3 and place 2 is 3
     (place 4, 3),   //  the shortest path between place 3 and place 4 is 3
     (place 1, 4),   //  the shortest path between place 3 and place 1 is 4
With limitation = 3, it will return
     (place 2, 3),   //  the shortest path between place 3 and place 2 is 3
     (place 4, 3),   //  the shortest path between place 3 and place 4 is 3

*/
package topograph
