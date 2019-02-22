# GoTraverse
### Installation
Download a release or directly build the code from this repository
```
go get github.com/colossus21/gotraverse
```
Import the library in your project
```go
import "github.com/colossus21/gotraverse"
``` 

## Getting Started
### Create Node Map

To work with nodes you need to create a map of nodes with heuristic values.

```go 
nodesStr := "S 8 A 8 B 4 C 3 D inf E inf G 0"
Nodes := gotraverse.MakeNodes(inp)
```
MakeNodes takes a string with nodes separated by its heuristic values and returns a Node map.
Note that,heuristic values can be of any numeric positive value and a special value called "inf" which is accepted as infinity.
### Add Edges

Edges can be added very easily using MakeEdges function.

```go 
edgesStr := "S A 3 S B 1 S C 8 A D 3 A E 7 A G 15 B G 20 C G 5"
Edges := gotraverse.MakeEdges(edgesStr, Nodes, "S")
```
MakeEdges takes a string value with all the edge information, the Node map and a node to return. ie. "S A 3" creates a directed edge from S to A with distance 3. All connections are separated by space as above. 

### Run Algorithm and Capture Node

You can apply the provided graph algorithms and return the node when a goal node is stated.

```go 
nodeFound := gotraverse.AStar(Edges, Nodes["G"])
```
All algorithms provided takes the edges and a goal node as parameter. If you are willing to traverse all nodes in a graph, provide a blank string.
nodeFound will have the goal Node with the path it took to reach there and prints the following:
```
[S 0]
S
[B 5][A 11][C 11]
B
[C 11][A 11][G 21]
C
[A 11][G 21][G 13]
A
[G 13][G 21][D 2147483653][E 2147483657]
G
```
You can further check the paths traversed by nodeFound by using the provided PrintNodePaths function:
```go
gotraverse.PrintNodePaths(nodeFound)
```
The following is printed in the console:
```
Node: S G: 0 H: 8 F: 8
Node: C G: 8 H: 3 F: 3
Node: G G: 13 H: 0 F: 0
```
By using this table you can check all the paths the node took to reach its goal along with its heuristic and distance information.
## Demo
The Nodes and Edges information used above was taken from the following graph, let's check out how easily we traverse the following graph implementing the provided algorithms. Console will provide the trace notes:
![alt text](/img/Graph.png)

##### A* Search
```go
gotraverse.AStar(Edges, Nodes["G"])
```
```
[S 0]
S
[B 5][A 11][C 11]
B
[C 11][A 11][G 21]
C
[A 11][G 21][G 13]
A
[G 13][G 21][D 2147483653][E 2147483657]
G
```

##### Greedy Search
```go
gotraverse.GreedySearch(Edges, Nodes["G"])
```
```
[S 8]
S
[C 3][A 8][B 4]
C
[G 0][A 8][B 4]
G
```
##### Uniform Cost Search (UCS)
```go
gotraverse.UCS(Edges, Nodes["G"])
```
```
[S 0]
S
[B 1][A 3][C 8]
B
[A 3][C 8][G 21]
A
[D 6][E 10][C 8][G 21][G 18]
D
[C 8][E 10][G 18][G 21]
C
[E 10][G 13][G 18][G 21]
E
[G 13][G 21][G 18]
G
```
##### BFS
```go
gotraverse.BFS(Edges, Nodes["G"])
```
```
[S ]
Dequeue Node: S
[A B C ]
Dequeue Node: B
[B C D E G ]
Dequeue Node: C
[C D E G ]
Dequeue Node: D
[D E G ]
Dequeue Node: D
[E G ]
Dequeue Node: E
[G ]
Dequeue Node: G
```
##### DFS
```go
gotraverse.BFS(Edges, Nodes["G"])
```
```
[S ]
Pop Node: S
[A B C ]
Pop Node: C
[A B G ]
Pop Node: G
```


### And coding style tests


## Authors

* **Rafiul Alam** - *Initial work* - [Colossus21](https://github.com/colossus21)


## Acknowledgments

* PriorityQueue (https://golang.org/pkg/container/heap/#example__priorityQueue)
