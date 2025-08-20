package main

import (
	"fmt"
	"sort"
)

// Genera los sucesores según el nodo actual
func successors(node []int) [][]int {
	name := node[0]
	cost := node[1]

	switch name {
	case 1: // A
		return [][]int{
			{2, cost + 2}, // A->B
			{3, cost + 3}, // A->C
			{4, cost + 4}, // A->D
		}
	case 2: // B
		return [][]int{
			{3, cost + 1}, // B->C
			{4, cost + 3}, // B->D
		}
	case 3: // C
		return [][]int{
			{2, cost + 1}, // C->B
			{4, cost + 2}, // C->D
		}
	case 4: // D
		// D no genera sucesores porque es el objetivo
		return [][]int{}
	}
	return [][]int{}
}

// Búsqueda de costo uniforme
func uniformCost(begin, end int) {
	list := [][]int{{begin, 0}} // lista con nodo y costo acumulado
	visited := map[int]bool{}   // para evitar ciclos infinitos

	for len(list) > 0 {
		current := list[0]
		list = list[1:]

		if visited[current[0]] {
			continue
		}
		visited[current[0]] = true

		fmt.Println("Current Node:", current)

		if current[0] == end {
			fmt.Println("SOLUTION with cost:", current[1])
			return
		}

		tmp := successors(current)
		fmt.Println("Successors:", tmp)

		if len(tmp) > 0 {
			list = append(list, tmp...)
			sort.Slice(list, func(i, j int) bool {
				return list[i][1] < list[j][1] // ordenar por costo acumulado
			})
			fmt.Println("New List:", list)
		}
	}

	fmt.Println("NO-SOLUTION")
}

func main() {
	uniformCost(1, 4)
}
