package main

import (
	"fmt"
)

// Nodo representa un nodo con su identificador y nivel de profundidad
type Nodo struct {
	id    int
	nivel int
}

// successors retorna los sucesores válidos del nodo en una matriz 4x4, incluyendo diagonales
func successors(n int) []int {
	var succ []int

	row := (n - 1) / 4
	col := (n - 1) % 4

	// Direcciones posibles: fila, columna
	dirs := [8][2]int{
		{-1, -1}, // ↖
		{-1, 0},  // ↑
		{-1, 1},  // ↗
		{0, -1},  // ←
		{0, 1},   // →
		{1, -1},  // ↙
		{1, 0},   // ↓
		{1, 1},   // ↘
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

// BFS limitada por nivel
func breadthFirstSearchLimitado(inicio, fin, limite int) {
	lista := []Nodo{{id: inicio, nivel: 0}}

	for len(lista) > 0 {
		actual := lista[0]
		lista = lista[1:]

		fmt.Printf("Visitando nodo %d en nivel %d\n", actual.id, actual.nivel)

		if actual.id == fin {
			fmt.Println("SOLUCIÓN ENCONTRADA!")
			return
		}

		if actual.nivel < limite {
			sucesores := successors(actual.id)
			fmt.Printf("Sucesores de %d: %v\n", actual.id, sucesores)
			for _, s := range sucesores {
				lista = append(lista, Nodo{id: s, nivel: actual.nivel + 1})
			}
		}
	}
	fmt.Println("NO SE ENCONTRÓ SOLUCIÓN DENTRO DEL LÍMITE")
}

func main() {
	fmt.Println("Búsqueda por Anchura Limitada (BFS-L):")
	inicio := 1
	fin := 11
	limite := 2
	breadthFirstSearchLimitado(inicio, fin, limite)
}
