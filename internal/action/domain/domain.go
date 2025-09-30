package domain

import "time"

const (
	ActionTypeReferUser = "REFER_USER"
)

type Action struct {
	ID         int
	Type       string
	UserID     int
	TargetUser int
	CreatedAt  time.Time
}

// Graph represents a directed graph where each node is a user and edges represent invitations.
// If simplification is needed, needs to change nodes map to map[int][]int
type Graph struct {
	nodes map[int]*Node
}

func NewGraph() *Graph {
	nodes := make(map[int]*Node)
	return &Graph{
		nodes: nodes,
	}
}

type Node struct {
	UserID   int
	ParentID *int
	Children []*Node
}

func (g *Graph) AddEdge(parentUserID, childUserID int) {
	if parentUserID == childUserID {
		return // avoid self-loop
	}

	parent := g.getOrCreate(parentUserID)
	child := g.getOrCreate(childUserID)

	// Assume a user can be invited only ONCE.
	if child.ParentID == nil {
		child.ParentID = &parentUserID
	}

	// Avoid duplicate child references.
	for _, c := range parent.Children {
		if c.UserID == childUserID {
			return // already linked no actioun
		}
	}

	parent.Children = append(parent.Children, child)
}

func (g *Graph) getOrCreate(userID int) *Node {
	if n, ok := g.nodes[userID]; ok {
		return n
	}

	n := &Node{UserID: userID}
	g.nodes[userID] = n

	return n
}

func (g *Graph) ReferralCount(userID int) int {
	root, ok := g.nodes[userID]
	if !ok {
		return 0
	}

	seen := make(map[int]bool)

	childreen := countChildrenRec(root, seen)

	return childreen
}

func countChildrenRec(n *Node, seen map[int]bool) int {
	if n == nil || seen[n.UserID] {
		return 0 // fast path to avoid cycles
	}
	seen[n.UserID] = true

	total := 0
	for _, c := range n.Children {
		total += 1 + countChildrenRec(c, seen)
	}

	return total
}
