package main

import (
	"container/heap"
	"encoding/json"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// NOMBRE: JULIO ALEJANDRO ZALDAÑA RÍOS - CARNET: 202110206
// IA1 - practica 1

// ==========================
//  Representación del estado
// ==========================
// Estado del 8-puzzle: 9 casillas, 0 es el espacio en blanco
// Índices: 0..8 en orden de lectura (3x3)

type State [9]int

var goal = State{1, 2, 3, 4, 5, 6, 7, 8, 0}

func (s State) String() string {
	b := strings.Builder{}
	for i, v := range s {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(strconv.Itoa(v))
	}
	return b.String()
}

func indexOfZero(s State) int {
	for i, v := range s {
		if v == 0 {
			return i
		}
	}
	return -1
}

func (s State) neighbors() []State {
	// Movimientos válidos del 0 (arriba, abajo, izquierda, derecha)
	pos := indexOfZero(s)
	row, col := pos/3, pos%3
	dirs := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	res := make([]State, 0, 4)
	for _, d := range dirs {
		r, c := row+d[0], col+d[1]
		if r >= 0 && r < 3 && c >= 0 && c < 3 {
			npos := r*3 + c
			var ns State
			ns = s
			ns[pos], ns[npos] = ns[npos], ns[pos]
			res = append(res, ns)
		}
	}
	return res
}

// ==========================
//  Heurísticas
// ==========================

func manhattan(s State) int {
	dist := 0
	for i, v := range s {
		if v == 0 {
			continue
		}
		// pos objetivo de v (1..8 -> índice v-1; 0 -> índice 8)
		gi := v - 1
		gr, gc := gi/3, gi%3
		r, c := i/3, i%3
		if r < gr {
			dist += gr - r
		} else {
			dist += r - gr
		}
		if c < gc {
			dist += gc - c
		} else {
			dist += c - gc
		}
	}
	return dist
}

func misplaced(s State) int {
	cnt := 0
	for i, v := range s {
		if v == 0 {
			continue
		}
		if v != goal[i] {
			cnt++
		}
	}
	return cnt
}

// ==========================
//  A* Search
// ==========================

type node struct {
	state  State
	g, h   int
	f      int
	parent *node
	// heap index (para actualización si hiciera falta)
	index int
}

type priorityQueue []*node

func (pq priorityQueue) Len() int { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool {
	// Menor f primero; si empata, menor h (más prometedor)
	if pq[i].f == pq[j].f {
		return pq[i].h < pq[j].h
	}
	return pq[i].f < pq[j].f
}
func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}
func (pq *priorityQueue) Push(x any) {
	n := x.(*node)
	*pq = append(*pq, n)
	n.index = len(*pq) - 1
}
func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[:n-1]
	return item
}

func reconstructPath(n *node) []State {
	var rev []State
	for cur := n; cur != nil; cur = cur.parent {
		rev = append(rev, cur.state)
	}
	// invertir
	for i, j := 0, len(rev)-1; i < j; i, j = i+1, j-1 {
		rev[i], rev[j] = rev[j], rev[i]
	}
	return rev
}

func isSolvable(s State) bool {
	// Para 8-puzzle (3x3): solvable si número de inversiones es par
	arr := make([]int, 0, 9)
	for _, v := range s {
		if v != 0 {
			arr = append(arr, v)
		}
	}
	inv := 0
	for i := 0; i < len(arr); i++ {
		for j := i + 1; j < len(arr); j++ {
			if arr[i] > arr[j] {
				inv++
			}
		}
	}
	return inv%2 == 0
}

func aStar(start State, heuristicName string) (path []State, expanded int, ok bool) {
	var hfunc func(State) int
	switch heuristicName {
	case "misplaced":
		hfunc = misplaced
	default:
		hfunc = manhattan
	}

	startH := hfunc(start)
	startNode := &node{state: start, g: 0, h: startH, f: startH}

	open := &priorityQueue{}
	heap.Init(open)
	heap.Push(open, startNode)

	came := make(map[string]*node) // best known node per state
	came[start.String()] = startNode

	gBest := make(map[string]int)
	gBest[start.String()] = 0

	for open.Len() > 0 {
		cur := heap.Pop(open).(*node)
		expanded++
		if cur.state == goal {
			return reconstructPath(cur), expanded, true
		}

		for _, nb := range cur.state.neighbors() {
			ng := cur.g + 1
			key := nb.String()
			best, seen := gBest[key]
			if !seen || ng < best {
				h := hfunc(nb)
				n := &node{state: nb, g: ng, h: h, f: ng + h, parent: cur}
				gBest[key] = ng
				came[key] = n
				heap.Push(open, n)
			}
		}
	}
	return nil, expanded, false
}

// ==========================
//  Utilidades de mezclado
// ==========================

func randomShuffle(nSteps int) State {
	// Mezcla por paseo aleatorio desde el objetivo => siempre resoluble
	s := goal
	prev := State{}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < nSteps; i++ {
		nbs := s.neighbors()
		// Evitar deshacer el último movimiento si es posible
		candidates := make([]State, 0, len(nbs))
		for _, nb := range nbs {
			if nb != prev {
				candidates = append(candidates, nb)
			}
		}
		if len(candidates) == 0 {
			candidates = nbs
		}
		prev = s
		s = candidates[rand.Intn(len(candidates))]
	}
	return s
}

// ==========================
//  Servidor HTTP + Frontend
// ==========================

type solveRequest struct {
	State     State  `json:"state"`
	Heuristic string `json:"heuristic"`
}

type solveResponse struct {
	Steps    []State `json:"steps"`
	Expanded int     `json:"expanded"`
	Cost     int     `json:"cost"`
	Message  string  `json:"message"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := indexTmpl.Execute(w, nil); err != nil {
			log.Println("template error:", err)
		}
	})

	mux.HandleFunc("/init", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(goal)
	})

	mux.HandleFunc("/shuffle", func(w http.ResponseWriter, r *http.Request) {
		n := 30
		if v := r.URL.Query().Get("n"); v != "" {
			if vi, err := strconv.Atoi(v); err == nil && vi > 0 {
				n = vi
			}
		}
		s := randomShuffle(n)
		json.NewEncoder(w).Encode(s)
	})

	mux.HandleFunc("/solve", func(w http.ResponseWriter, r *http.Request) {
		var req solveRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "JSON inválido"})
			return
		}
		if !isSolvable(req.State) {
			json.NewEncoder(w).Encode(solveResponse{Steps: nil, Expanded: 0, Cost: 0, Message: "Estado no resoluble"})
			return
		}
		path, expanded, ok := aStar(req.State, req.Heuristic)
		if !ok {
			json.NewEncoder(w).Encode(solveResponse{Steps: nil, Expanded: expanded, Cost: 0, Message: "No se encontró solución"})
			return
		}
		json.NewEncoder(w).Encode(solveResponse{Steps: path, Expanded: expanded, Cost: len(path) - 1, Message: "OK"})
	})

	addr := ":8080"
	log.Println("Servidor listo en http://localhost" + addr)
	log.Fatal(http.ListenAndServe(addr, logRequest(mux)))
}

func logRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

var indexTmpl = template.Must(template.New("idx").Parse(`<!doctype html>
<html lang="es">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>8-Puzzle — Go - IA1 - 202110206</title>
  <style>
    :root { --gap: 10px; --size: 90vmin; --cell: calc((var(--size) - 2*var(--gap)) / 3);}
    body { font-family: system-ui, -apple-system, Segoe UI, Roboto, Ubuntu, Cantarell, 'Helvetica Neue', Arial, 'Noto Sans', sans-serif; display:flex; align-items:center; justify-content:center; min-height:100vh; margin:0; background:#0f172a; color:#e5e7eb; }
    .app { width:min(900px, 96vw); display:grid; grid-template-columns: 1fr 1fr; gap: 20px; }
    .board { width:var(--size); max-width: 96vw; aspect-ratio:1/1; display:grid; grid-template-columns: repeat(3, 1fr); gap: var(--gap); padding: var(--gap); background:#111827; border-radius: 18px; box-shadow: 0 10px 30px rgba(0,0,0,.35); }
    .tile { display:flex; align-items:center; justify-content:center; font-weight:700; font-size: clamp(28px, 7vmin, 54px); background:#1f2937; border-radius: 16px; user-select:none; transition: transform .12s ease; box-shadow: inset 0 1px 0 rgba(255,255,255,.06), 0 6px 18px rgba(0,0,0,.25); }
    .tile:not(.empty):active { transform: scale(.97); }
    .tile.empty { background: #0b1220; box-shadow: inset 0 0 0 2px rgba(148,163,184,.25); }
    .panel { background:#0b1220; border-radius: 18px; padding: 16px; box-shadow: 0 10px 30px rgba(0,0,0,.35); }
    .row { display:flex; gap: 10px; flex-wrap: wrap; align-items:center; }
    button { background:#3b82f6; color:white; border:0; padding:10px 14px; border-radius:12px; font-weight:600; cursor:pointer; box-shadow: 0 6px 18px rgba(59,130,246,.35); }
    button.secondary { background:#6b7280; box-shadow: 0 6px 18px rgba(107,114,128,.35); }
    button.danger { background:#ef4444; box-shadow: 0 6px 18px rgba(239,68,68,.35); }
    select, input[type="number"] { background:#111827; color:#e5e7eb; border:1px solid #374151; border-radius: 10px; padding: 8px 10px; }
    .stats { margin-top: 12px; font-size: 14px; color:#cbd5e1; }
    .hint { font-size: 12px; color:#94a3b8; }
    @media (max-width: 860px) { .app { grid-template-columns: 1fr; place-items:center; } .board { --size: 84vmin; } }
  </style>
</head>
<body>
  <div class="app">
    <div>
      <div id="board" class="board"></div>
    </div>
    <div class="panel">
      <h2>8‑Puzzle</h2>
      <div class="row" style="margin:10px 0 12px">
        <button id="btnInit" class="secondary">Iniciar</button>
        <input id="steps" type="number" min="1" max="200" value="30" style="width:86px" />
        <button id="btnShuffle">Desordenar</button>
      </div>
      <div class="row" style="margin:4px 0 12px">
        <select id="heuristic">
          <option value="manhattan" selected>Heurística: Manhattan</option>
          <option value="misplaced">Heurística: Fichas fuera de lugar</option>
        </select>
      </div>
      <div class="row" style="margin:4px 0 12px">
        <button id="btnSolve">Resolver</button>
        <button id="btnStep" class="danger">Resolver paso a paso</button>
      </div>
      <div class="stats" id="stats"></div>
      <div class="hint">Tip: Puedes cambiar la heurística. "Resolver" anima todo el camino; "Paso a paso" muestra el siguiente estado en cada clic.</div>
    </div>
  </div>

<script>
  const boardEl = document.getElementById('board');
  const btnInit = document.getElementById('btnInit');
  const btnShuffle = document.getElementById('btnShuffle');
  const btnSolve = document.getElementById('btnSolve');
  const btnStep = document.getElementById('btnStep');
  const stepsInput = document.getElementById('steps');
  const heuristicSel = document.getElementById('heuristic');
  const statsEl = document.getElementById('stats');

  let state = [];
  let solution = null; // array de estados
  let stepIndex = 0;
  let playing = false;

  function render(s) {
    boardEl.innerHTML = '';
    s.forEach((v, i) => {
      const d = document.createElement('div');
      d.className = 'tile' + (v === 0 ? ' empty' : '');
      d.textContent = v === 0 ? '' : v;
      boardEl.appendChild(d);
    });
  }

  async function initBoard() {
    const res = await fetch('/init');
    state = await res.json();
    solution = null; stepIndex = 0; playing = false; statsEl.textContent='';
    render(state);
  }

  async function shuffleBoard() {
    const n = Number(stepsInput.value) || 30;
    const res = await fetch('/shuffle?n=' + encodeURIComponent(n));
    state = await res.json();
    solution = null; stepIndex = 0; playing = false; statsEl.textContent = 'Estado mezclado con ' + n + ' movimientos.';
    render(state);
  }

  async function ensureSolution() {
    if (solution) return solution;
    const payload = { state, heuristic: heuristicSel.value };
    const res = await fetch('/solve', { method:'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(payload) });
    const data = await res.json();
    if (data.message !== 'OK') {
      statsEl.textContent = data.message || 'No se pudo resolver.';
      return null;
    }
    solution = data.steps;
    stepIndex = 0;
    statsEl.textContent = 'Longitud solución:' + data.cost + ' | Nodos expandidos: ' + data.expanded;
    return solution;
  }

  async function solveAnimated() {
    if (playing) return; 
    const sol = await ensureSolution();
    if (!sol) return;
    playing = true;
    let i = 0;
    const timer = setInterval(() => {
      state = sol[i];
      render(state);
      i++;
      if (i >= sol.length) {
        clearInterval(timer);
        playing = false;
      }
    }, 350);
  }

  async function solveStep() {
    const sol = await ensureSolution();
    if (!sol) return;
    if (stepIndex < sol.length) {
      state = sol[stepIndex];
      render(state);
      stepIndex++;
    }
  }

  btnInit.addEventListener('click', initBoard);
  btnShuffle.addEventListener('click', shuffleBoard);
  btnSolve.addEventListener('click', solveAnimated);
  btnStep.addEventListener('click', solveStep);

  // Inicio
  initBoard();
</script>
</body>
</html>
`))
