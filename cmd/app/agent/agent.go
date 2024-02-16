package agent

import (
	"container/list"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	server "1/cmd/app/server"

	_ "github.com/mattn/go-sqlite3"
)

// переменные имитации задержки
var (
	plus  int = 5
	minus int = 1
	div   int = 1
	mult  int = 1
)

var mapa = make(map[string]float64)

type worker struct {
	pool     *Pool
	jobsChan chan server.Expression
	quit     chan *sync.WaitGroup
}

type Operations struct {
	Operation string `json:"operation"`
	Time      string `json:"time"`
}

var (
	ErrPoolClosed = errors.New("pool is closed")
	ErrQueueFull  = errors.New("queue is full")
)

type Pool struct {
	// Amount of workers in pool
	Size int
	// Amount of jobs that can be in queue
	QueueSize int

	finish      bool
	jobsQueue   chan server.Expression
	freeWorkers chan *worker
	workers     *list.List
}

func dbAgent() *sql.DB {
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS operations (Operation TEXT, Time TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	db.Exec("PRAGMA journal_mode=WAL")

	return db
}

func UpdateOperations(w http.ResponseWriter, r *http.Request) {
	db := dbAgent()
	defer db.Close()
	var ops Operations
	err := json.NewDecoder(r.Body).Decode(&ops)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE operations SET time = ? WHERE operation = ?", ops.Time, ops.Operation)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT * FROM operations")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var operation, time string
		if err := rows.Scan(&operation, &time); err != nil {
			log.Fatal(err)
		}
		switch operation {
		case "plus":
			plus, _ = strconv.Atoi(time)
		case "minus":
			minus, _ = strconv.Atoi(time)
		case "mult":
			mult, _ = strconv.Atoi(time)
		case "div":
			div, _ = strconv.Atoi(time)
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func (w *worker) start() {
	w.jobsChan = make(chan server.Expression, 1)
	w.quit = make(chan *sync.WaitGroup, 1)

	go func() {
		for {
			w.pool.freeWorkers <- w

			select {
			case job := <-w.jobsChan:
				w.doJob(job)

			case wg := <-w.quit:
				wg.Done()
				return
			}
		}
	}()
}

func (w *worker) doJob(job server.Expression) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in job %s: %s", job, r)
		}
	}()

	mapa[job.MathExpr] = EvaluateExpression(job.MathExpr)
}

func (p *Pool) Init() {
	p.jobsQueue = make(chan server.Expression, p.QueueSize)
	p.freeWorkers = make(chan *worker, p.Size)
	p.workers = list.New()
}

func (p *Pool) Start() {
	for i := 0; i < p.Size; i++ {
		w := &worker{
			pool: p,
		}
		p.workers.PushFront(w)
		w.start()
	}

	go func() {
		for job := range p.jobsQueue {
			fmt.Println(job)
			// Wait for the free worker
			w := <-p.freeWorkers

			// Send job to worker
			w.jobsChan <- job
		}
	}()
}

func (p *Pool) AddJob(data server.Expression) error {
	if p.finish {
		return ErrPoolClosed
	}
	select {
	case p.jobsQueue <- data:
	default:
		return ErrQueueFull
	}
	return nil
}

func (p *Pool) GetQueueLength() int {
	return len(p.jobsQueue)
}

func (p *Pool) GetActiveWorkers() int {
	return p.Size - len(p.freeWorkers)
}

func (p *Pool) Finish() {
	log.Println("Finishing all jobs...")
	p.finish = true
	for len(p.jobsQueue) != 0 {
		time.Sleep(50 * time.Millisecond)
	}
	wg := &sync.WaitGroup{}
	wg.Add(p.Size)
	for e := p.workers.Front(); e != nil; e = e.Next() {
		e.Value.(*worker).quit <- wg
	}
	wg.Wait()
}

func EvaluateExpression(expr string) float64 {
	tokens := tokenize(expr)
	// Преобразование в обратную польскую нотацию
	rpn := shuntingYard(tokens)
	// Вычисление результата
	result, _ := evaluateRPN(rpn)
	return result
}

func tokenize(expr string) []string {
	// Разбиваем выражение на токены
	tokens := []string{}
	buffer := ""
	for _, char := range expr {
		if char == ' ' {
			continue
		} else if strings.Contains("+-*/()", string(char)) {
			if len(buffer) > 0 {
				tokens = append(tokens, buffer)
				buffer = ""
			}
			tokens = append(tokens, string(char))
		} else {
			buffer += string(char)
		}
	}
	if len(buffer) > 0 {
		tokens = append(tokens, buffer)
	}
	return tokens
}

func shuntingYard(tokens []string) []string {
	// Алгоритм преобразования в обратную польскую нотацию
	var rpn []string
	var stack []string
	precedence := map[string]int{"+": 1, "-": 1, "*": 2, "/": 2}
	for _, token := range tokens {
		switch {
		case token == "(":
			stack = append(stack, token)
		case token == ")":
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				rpn = append(rpn, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = stack[:len(stack)-1] // Удаление "("
		case precedence[token] > 0:
			for len(stack) > 0 && precedence[stack[len(stack)-1]] >= precedence[token] {
				rpn = append(rpn, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		default: // Операнд
			rpn = append(rpn, token)
		}
	}
	for len(stack) > 0 {
		rpn = append(rpn, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	return rpn
}

func evaluateRPN(tokens []string) (float64, error) {
	var stack []float64
	for _, token := range tokens {
		if !strings.Contains("+-*/", token) {
			value, err := strconv.Atoi(token)
			if err != nil {
				return 0, err
			}
			stack = append(stack, float64(value))
		} else {
			if len(stack) < 2 {
				return 0, fmt.Errorf("Invalid expression")
			}
			operand2, operand1 := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			switch token {
			case "+":
				time.Sleep(time.Duration(plus) * time.Second)
				stack = append(stack, operand1+operand2)
			case "-":
				time.Sleep(time.Duration(minus) * time.Second)
				stack = append(stack, operand1-operand2)
			case "*":
				time.Sleep(time.Duration(mult) * time.Second)
				stack = append(stack, operand1*operand2)
			case "/":
				time.Sleep(time.Duration(div) * time.Second)
				stack = append(stack, operand1/operand2)
			}
		}
	}
	if len(stack) == 1 {
		return stack[0], nil
	}
	return 0, fmt.Errorf("Invalid expression")
}

func Process(Expressions []server.Expression) {
	pool := &Pool{
		Size:      len(Expressions),
		QueueSize: len(Expressions),
	}

	pool.Init()

	//pool.Start()

	// Add some jobs to the pool
	for i := 0; i < len(Expressions); i++ {
		if err := pool.AddJob(Expressions[i]); err != nil {
			log.Printf("Error adding job: %v", err)
		}
	}

	pool.Start()
	// log.Printf("Queue length: %d", pool.GetQueueLength())
	// log.Printf("Active workers: %d", pool.GetActiveWorkers())

	// Finish all jobs and workers
	pool.Finish()
}

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	db := server.Database()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM Expressions WHERE Status='pending'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var res []server.Expression
	for rows.Next() {
		var u server.Expression
		if err := rows.Scan(&u.ID, &u.MathExpr, &u.Status, &u.Result); err != nil {
			log.Fatal(err)
		}
		res = append(res, u)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	Process(res)

	for _, expr := range res {
		tokens := tokenize(expr.MathExpr)
		rpnTokens := shuntingYard(tokens)
		result, err := evaluateRPN(rpnTokens)
		if err != nil {
			log.Printf("Error evaluating expression: %v", err)
			continue
		}

		_, err = db.Exec("UPDATE Expressions SET Status='complete', Result=$1 WHERE MathExpr=$2", result, expr.MathExpr)
		if err != nil {
			log.Fatal(err)
		}
	}
}
