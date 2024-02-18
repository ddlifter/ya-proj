package workers

import (
	calculate "1/cmd/app/calculate"
	"1/cmd/app/database"
	"container/list"
	"errors"
	"log"
	"sync"
	"time"
)

var mapa = make(map[string]float64) // Промежуточные данные
var workerID int = 1                // Для учета id агентов

// Описание работающей единицы
type Worker struct {
	LastPing int
	Id       int
	Status   string
	pool     *Pool
	jobsChan chan database.Expression
	quit     chan *sync.WaitGroup
}

var (
	ErrPoolClosed = errors.New("pool is closed")
	ErrQueueFull  = errors.New("queue is full")
)

// Описание пула задач
type Pool struct {
	// Количество доступных воркеров
	Size int
	// Количество задач которые могут быть в очереди
	QueueSize int

	finish      bool
	jobsQueue   chan database.Expression
	freeWorkers chan *Worker
	workers     *list.List
}

// Функиця запуска агента
func (w *Worker) Start() {
	w.jobsChan = make(chan database.Expression, 1)
	w.quit = make(chan *sync.WaitGroup, 1)

	go func() {
		for {
			// Добавляем воркера в пул
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

// Функция описывающая выполнение задачи
func (w *Worker) doJob(job database.Expression) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in job %s: %s", job, r)
		}
	}()

	// Обновление промежуточных данных
	mapa[job.MathExpr] = calculate.EvaluateExpression(job.MathExpr)
	w.Status = "completed"

}

// Инициализация пула задач
func (p *Pool) Init() {
	p.jobsQueue = make(chan database.Expression, p.QueueSize)
	p.freeWorkers = make(chan *Worker, p.Size)
	p.workers = list.New()
}

// Запуск пула
func (p *Pool) Start() {
	for i := 0; i < p.Size; i++ {
		w := &Worker{
			pool: p,
		}
		p.workers.PushFront(w)
		w.Start()
		w.Status = "working"
		w.Id = workerID
		workerID++
		db := database.DbAgent()
		defer db.Close()

		// Добавляем новых агентов для мониторинга
		_, err := db.Exec("INSERT INTO agents (Status, LastPing) VALUES (?, ?)", "working", 0)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Вынимаем свободных агентов
	go func() {
		for job := range p.jobsQueue {
			// Wait for the free worker
			w := <-p.freeWorkers

			// Send job to worker
			w.jobsChan <- job
		}
	}()
}

// Добавить задачу в пул
func (p *Pool) AddJob(data database.Expression) error {
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

// Завершение задач, закрытие пула
func (p *Pool) Finish() {
	log.Println("Finishing all jobs...")
	p.finish = true
	for len(p.jobsQueue) != 0 {
		time.Sleep(50 * time.Millisecond)
	}
	wg := &sync.WaitGroup{}
	wg.Add(p.Size)
	for e := p.workers.Front(); e != nil; e = e.Next() {
		e.Value.(*Worker).quit <- wg
	}
	wg.Wait()
}

// Выполнение задачи, взаимодействие с бд
func Process(Expressions []database.Expression) {
	db := database.DbAgent()
	defer db.Close()

	// Удаляем старых агентов
	rows, err := db.Query("SELECT id FROM agents")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var w Worker
		if err := rows.Scan(&w.Id); err != nil {
			log.Fatal(err)
		}
		_, err := db.Exec("DELETE FROM agents WHERE id = ?", w.Id)
		if err != nil {
			return
		}
	}

	// Инициализация пула
	pool := &Pool{
		Size:      len(Expressions),
		QueueSize: len(Expressions),
	}

	pool.Init()

	// Добавляем задачи в пул
	for i := 0; i < len(Expressions); i++ {
		if err := pool.AddJob(Expressions[i]); err != nil {
			log.Printf("Error adding job: %v", err)
		}
	}
	pool.Start()

	// Закрываем пул
	pool.Finish()
}
