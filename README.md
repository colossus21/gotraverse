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
### Create Graph

Let's create a graph considering a parent node, "S"

```go 
nodesStr := "S 8 A 8 B 4 C 3 D inf E inf G 0"
edgesStr := "S A 3 S B 1 S C 8 A D 3 A E 7 A G 15 B G 20 C G 5"
g := gotraverse.MakeGraph(nodesStr, edgesStr, "S") 
```
Nodes are represented as name and heuristic value separated by space. ie. S 8 means node named S has heuristic value 8. Edges are represented as starting node, connected node and weight/distance seperated by spaces. ie. S A 3 indicates S is connected to A with a weight of 3. 

### Run Algorithm and Capture Node

All algorithms come with the following interface

```go
type GoalSearch interface {
	SetIterativeNode(node *Node) GoalSearch
	SetGoalNode(node *Node) GoalSearch
	Search() GoalSearch
	GetCapturedNode() *Node
}
```
Let's search for a goal node using A* algorithm
```go
// Create an instance of the algorithm
astar := new(gotraverse.AStar)
// Perform a search operation
g.InitGoalSearch(searchAlgorithm, "G").Search()
// Print distance taken
fmt.Println(astar.GetCapturedNode().GetTotalDistance())
// Traverse the node to check all it's path choices
astar.GetCapturedNode().Traverse()
```
Output
```
[S 8]
S
[B 5][A 11][C 11]
B
[C 11][A 11][G 21]
C
[A 11][G 21][G 13]
A
[G 13][G 21][D 2147483653][E 2147483657]
G

13

Node: S G: 0 H: 8 F: 8
Node: C G: 8 H: 3 F: 3
Node: G G: 13 H: 0 F: 0
```


## Demo
The Nodes and Edges information used above was taken from the following graph
![alt text](/img/Graph.png)

##### A* Search
```go
g.InitGoalSearch( new(gotraverse.AStar), "G").Search()
```
```
[S 8]
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
g.InitGoalSearch( new(gotraverse.GreedySearch), "G").Search()
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
g.InitGoalSearch( new(gotraverse.UCS), "G").Search()
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
g.InitGoalSearch( new(gotraverse.BFS), "G").Search()
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
g.InitGoalSearch( new(gotraverse.DFS), "G").Search()
```
```
[S ]
Pop Node: S
[A B C ]
Pop Node: C
[A B G ]
Pop Node: G
```


*Project is under development. Contributions are appreciated.*   




## Acknowledgments

* Official Golang Website (https://golang.org/pkg/container/heap/#example__priorityQueue)
