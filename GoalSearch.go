package gotraverse

type GoalSearch interface {
	SetIterativeNode(node *Node) GoalSearch
	SetGoalNode(node *Node) GoalSearch
	Search() GoalSearch
	GetCapturedNode() *Node
}
