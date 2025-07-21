package main

import (
	"fmt"
	"math/rand"
	"time"
)

func reflexAgent(location, state string) string {
	if state == "DIRTY" {
		return "CLEAN"
	} else if location == "A" {
		return "RIGHT"
	} else if location == "B" {
		return "LEFT"
	}
	return ""
}

func generateStateKey(states []string) string {
	return fmt.Sprintf("(%s,%s,%s)", states[0], states[1], states[2])
}

func maybeMakeDirty(states []string) {
	n := rand.Intn(10) + 1 // 1 al 10
	if n == 6 || n == 7 {
		states[1] = "DIRTY" // A se ensucia
	} else if n == 8 || n == 9 {
		states[2] = "DIRTY" // B se ensucia
	} else if n == 10 {
		states[1] = "DIRTY"
		states[2] = "DIRTY"
	}
}

func run(states []string) {
	visited := make(map[string]bool)

	for {
		// Guardar el estado actual
		key := generateStateKey(states)
		if !visited[key] {
			fmt.Println("Visitando nuevo estado:", key)
			visited[key] = true
		}

		// Verificar si ya se visitaron los 8 estados
		if len(visited) == 8 {
			fmt.Println("\n¡Todos los 8 estados han sido visitados!")
			break
		}

		// Determinar ubicación y estado actual
		location := states[0]
		var state string
		if location == "A" {
			state = states[1]
		} else {
			state = states[2]
		}

		// Acción del agente
		action := reflexAgent(location, state)
		fmt.Printf("Ubicación: %s | Acción: %s\n", location, action)

		if action == "CLEAN" {
			if location == "A" {
				states[1] = "CLEAN"
			} else {
				states[2] = "CLEAN"
			}
		} else if action == "RIGHT" {
			states[0] = "B"
		} else if action == "LEFT" {
			states[0] = "A"
		}

		// Ensuciar aleatoriamente después de cada paso
		maybeMakeDirty(states)

		time.Sleep(1 * time.Second)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())          // Semilla para aleatoriedad
	states := []string{"A", "DIRTY", "DIRTY"} // [location, A_state, B_state]
	run(states)
}
