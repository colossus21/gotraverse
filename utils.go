package gotraverse

import (
	"fmt"
	"strconv"
	"strings"
)

// creates a Map of nodes from the given string
// string is space separated by node name and heuristic value
func MakeNodes(inputStr string)  map[string]*Node {
	// Eg. input: S 8 A 8 B 4 C 3 D inf E inf G 0 ie. S is the node name, 8 is the heuristic value
	nodesMap := map[string]*Node {}

	//Equivalent to JS string.split(arr, " ')
	inputs := strings.Fields(inputStr)
	for i:=0; i<len(inputs); i+=2 {
		nodeName := string(inputs[i])
		node := new(Node)
		node.name = nodeName
		if inputs[i+1] != "inf" {
			node.heuristicValue, _ = strconv.Atoi(inputs[i+1])
		} else {
			node.heuristicValue = INF
		}
		nodesMap[nodeName] = node
	}

	return nodesMap
}
// Create the graph using the given Node Map having taken the edges
// Edges are taken as string. String contains nodes and weight separated by space.
func MakeEdges(inputStr string, nodes map[string]*Node, startNode string) *Node {
	// Eg input S A 3 S B 1 S C 3 A D 3 A E 7 A G 15 B G 20 C G 5 ie. from node S to node A with distance 3
	inputs := strings.Fields(inputStr)
	for i:=0; i<len(inputs)-2; i+=3 {
		edge := new(Edge)
		edge.from = nodes[string(inputs[i])]
		edge.to = nodes[string(inputs[i+1])]
		val, err := strconv.Atoi(inputs[i+2])
		if err!=nil {
			panic(err)
		}
		edge.distance = val
		nodes[string(inputs[i])].edges = append(nodes[string(inputs[i])].edges, edge)
	}
	return nodes[startNode]
}


// convert node item to priority queue item
func NodeItem(node *Node, priority int) *Item {
	newItem := new(Item)
	newItem.value = node
	newItem.priority = priority
	return newItem
}

// traverses the node and print from start node to goal node
// takes node and {prints msg [g, h, f, n]} (optional)
func PrintNodePaths(node *Node, msg... string)  {
	newNode := new(Node)
	newNode.name = node.name
	newNode.parent = node.parent
	newNode.heuristicValue = node.heuristicValue
	newNode.distanceFromStartNode = node.distanceFromStartNode

	g := "G:"
	h := "H:"
	f := "F:"
	n := "Node:"

	if (msg != nil && len(msg)>3) {
		n = msg[0]
		g = msg[1]
		h = msg[2]
		f = msg[3]
	}

	for newNode != nil {
		defer fmt.Println(n, newNode.name, g, newNode.distanceFromStartNode, h, newNode.heuristicValue, f, newNode.heuristicValue)
		newNode = newNode.parent
	}
}




