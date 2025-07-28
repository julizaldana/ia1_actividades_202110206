package main

import (
	"fmt"
)

// Nodo representa un nodo con su identificador y nivel de profundidad
type Nodo struct {
	id    int
	nivel int
	padre int // para registrar de dónde venimos
}

// successors retorna los sucesores válidos del nodo en una matriz 4x4 (incluyendo diagonales)
func successors(n int) []int {
	var succ []int
	row := (n - 1) / 4
	col := (n - 1) % 4

	// Direcciones posibles (8): ↖ ↑ ↗ ← → ↙ ↓ ↘
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

// DFS limitado con backjumping
func depthFirstSearchBackjumping(inicio, fin, limite int) {
	lista := []Nodo{{id: inicio, nivel: 0, padre: -1}}
	visitados := make(map[int]bool)
	nodosEliminados := 0
	backjumps := []string{}

	for len(lista) > 0 {
		actual := lista[0]
		lista = lista[1:]

		fmt.Printf("Visitando nodo %d en nivel %d\n", actual.id, actual.nivel)
		visitados[actual.id] = true

		if actual.id == fin {
			fmt.Println("¡SOLUCIÓN ENCONTRADA!")
			return
		}

		if actual.nivel < limite {
			sucesores := successors(actual.id)
			reverse(sucesores)
			for _, s := range sucesores {
				if !visitados[s] {
					nuevo := Nodo{id: s, nivel: actual.nivel + 1, padre: actual.id}
					lista = append([]Nodo{nuevo}, lista...)
				}
			}
		} else {
			// Se alcanzó el límite de profundidad
			// Eliminamos todos los nodos con nivel > 1 (backjump)
			restante := []Nodo{}
			for _, n := range lista {
				if n.nivel <= 1 {
					restante = append(restante, n)
				} else {
					nodosEliminados++
					if n.padre != -1 {
						backjumps = append(backjumps, fmt.Sprintf("Backjump de %d a %d", n.padre, inicio))
					}
				}
			}
			lista = restante
		}
	}

	fmt.Println("NO SE ENCONTRÓ SOLUCIÓN DENTRO DEL LÍMITE")
	fmt.Printf("Nodos eliminados por backjumping: %d\n", nodosEliminados)
	for _, b := range backjumps {
		fmt.Println(b)
	}
}

// reverse invierte el orden de los sucesores (para DFS)
func reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func main() {
	fmt.Println("Búsqueda por Profundidad Limitada con Backjumping:")
	inicio := 1
	fin := 11
	limite := 2
	depthFirstSearchBackjumping(inicio, fin, limite)
}
