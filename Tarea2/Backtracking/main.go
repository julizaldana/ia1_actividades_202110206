package main

import (
	"fmt"
	"sort"
)

// Nodo representa un nodo con identificador y nivel
type Nodo struct {
	id    int
	nivel int
}

var backtrackCount int = 0

// successors con movimientos diagonales en una matriz 4x4
func successors(n int) []int {
	var succ []int

	row := (n - 1) / 4
	col := (n - 1) % 4

	dirs := [8][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for _, dir := range dirs {
		newRow := row + dir[0]
		newCol := col + dir[1]
		if newRow >= 0 && newRow < 4 && newCol >= 0 && newCol < 4 {
			succ = append(succ, newRow*4+newCol+1)
		}
	}
	return succ
}

// DFS limitado con backtracking explícito
func dfsLimitado(nodo Nodo, fin int, limite int, visitados map[int]bool) bool {
	fmt.Printf("Visitando nodo %d en nivel %d\n", nodo.id, nodo.nivel)

	// Marcar como visitado en esta rama
	visitados[nodo.id] = true

	// Si encontramos la solución
	if nodo.id == fin {
		fmt.Println("¡SOLUCIÓN ENCONTRADA!")
		return true
	}

	if nodo.nivel < limite {
		sucesores := successors(nodo.id)
		sort.Ints(sucesores)
		fmt.Printf("Sucesores de %d: %v\n", nodo.id, sucesores)

		for _, s := range sucesores {
			if !visitados[s] {
				nuevoVisitados := copiarMapa(visitados)
				encontrado := dfsLimitado(Nodo{id: s, nivel: nodo.nivel + 1}, fin, limite, nuevoVisitados)
				if encontrado {
					return true
				}
			}
		}
	}

	// Si llegamos aquí, hicimos backtracking
	backtrackCount++
	return false
}

// Copia el mapa para mantener visitados por rama
func copiarMapa(original map[int]bool) map[int]bool {
	nuevo := make(map[int]bool)
	for k, v := range original {
		nuevo[k] = v
	}
	return nuevo
}

// reverse invierte los sucesores (para DFS)
func reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func main() {
	fmt.Println("Búsqueda por Profundidad Limitada (DFS-L) con Backtracking:")
	inicio := 1
	fin := 11
	limite := 2

	visitados := make(map[int]bool)
	encontrado := dfsLimitado(Nodo{id: inicio, nivel: 0}, fin, limite, visitados)

	if !encontrado {
		fmt.Println("NO SE ENCONTRÓ SOLUCIÓN DENTRO DEL LÍMITE")
	}
	fmt.Printf("Total de backtrackings realizados: %d\n", backtrackCount)
}
