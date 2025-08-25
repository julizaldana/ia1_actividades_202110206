package main

import (
	"container/heap"
	"fmt"
	"strings"
)

type Node struct {
	i, j    int
	g, h, f float64
	parent  *Node
	index   int // necesario para la priority queue
}

// Priority queue para A*
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].f < pq[j].f }
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}
func (pq *PriorityQueue) Push(x interface{}) {
	n := x.(*Node)
	n.index = len(*pq)
	*pq = append(*pq, n)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	*pq = old[0 : n-1]
	return node
}

// heurística: gaps mínimos restantes
func heuristic(i, j, lenA, lenB int, gap float64) float64 {
	remainA := lenA - i
	remainB := lenB - j
	diff := remainA - remainB
	if diff < 0 {
		diff = -diff
	}
	return float64(diff) * gap
}

func alignAStar(seqA, seqB string, match, mismatch, gap float64) {
	A := strings.Split(seqA, "")
	B := strings.Split(seqB, "")
	lenA, lenB := len(A), len(B)

	start := &Node{i: 0, j: 0, g: 0, h: heuristic(0, 0, lenA, lenB, gap)}
	start.f = start.g + start.h

	open := &PriorityQueue{}
	heap.Init(open)
	heap.Push(open, start)

	visited := make(map[[2]int]bool)
	costs := make(map[[2]int]float64) // guardar costos parciales

	var goal *Node

	for open.Len() > 0 {
		current := heap.Pop(open).(*Node)
		if visited[[2]int{current.i, current.j}] {
			continue
		}
		visited[[2]int{current.i, current.j}] = true
		costs[[2]int{current.i, current.j}] = current.g

		// Meta alcanzada
		if current.i == lenA && current.j == lenB {
			goal = current
			break
		}

		// Generar sucesores
		// Diagonal
		if current.i < lenA && current.j < lenB {
			cost := mismatch
			if A[current.i] == B[current.j] {
				cost = -match
			}
			next := &Node{
				i:      current.i + 1,
				j:      current.j + 1,
				g:      current.g + cost,
				h:      heuristic(current.i+1, current.j+1, lenA, lenB, gap),
				parent: current,
			}
			next.f = next.g + next.h
			heap.Push(open, next)
		}
		// Gap en B
		if current.i < lenA {
			next := &Node{
				i:      current.i + 1,
				j:      current.j,
				g:      current.g + gap,
				h:      heuristic(current.i+1, current.j, lenA, lenB, gap),
				parent: current,
			}
			next.f = next.g + next.h
			heap.Push(open, next)
		}
		// Gap en A
		if current.j < lenB {
			next := &Node{
				i:      current.i,
				j:      current.j + 1,
				g:      current.g + gap,
				h:      heuristic(current.i, current.j+1, lenA, lenB, gap),
				parent: current,
			}
			next.f = next.g + next.h
			heap.Push(open, next)
		}
	}

	// reconstruir alineamiento
	if goal != nil {
		var alignedA, alignedB []string
		cur := goal
		for cur.parent != nil {
			pi, pj := cur.parent.i, cur.parent.j
			if cur.i == pi+1 && cur.j == pj+1 {
				alignedA = append([]string{A[pi]}, alignedA...)
				alignedB = append([]string{B[pj]}, alignedB...)
			} else if cur.i == pi+1 {
				alignedA = append([]string{A[pi]}, alignedA...)
				alignedB = append([]string{"-"}, alignedB...)
			} else {
				alignedA = append([]string{"-"}, alignedA...)
				alignedB = append([]string{B[pj]}, alignedB...)
			}
			cur = cur.parent
		}
		fmt.Println("A:", strings.Join(alignedA, ""))
		fmt.Println("B:", strings.Join(alignedB, ""))
		fmt.Printf("Costo total: %.1f\n", goal.g)
	}

	// imprimir la "matriz parcial" de A*
	fmt.Println("\nMatriz parcial de nodos visitados:")
	header := "    " + strings.Join(B, " ")
	fmt.Println(header)
	for i := 0; i <= lenA; i++ {
		line := ""
		if i == 0 {
			line += " "
		} else {
			line += A[i-1]
		}
		for j := 0; j <= lenB; j++ {
			key := [2]int{i, j}
			if val, ok := costs[key]; ok {
				line += fmt.Sprintf(" %4.1f*", val)
			} else {
				line += "    . "
			}
		}
		fmt.Println(line)
	}
}

func main() {
	alignAStar("ABCDE", "ABCED", 1.0, 1.0, 2.0)
}
