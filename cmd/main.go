package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-co-op/gocron/v2"
)

const httpPort = ":3000"

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("welcome"))
		if err != nil {
			log.Println(err)
		}
	})

	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	jobs, err := initJobs(s)
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		log.Println("Starting server on port", httpPort)
		if err := http.ListenAndServe(":8080", nil); err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()

		fmt.Printf("Starting job: %d\n", jobs[0].ID())
		s.Start()
	}()

	wg.Wait()
}

func initJobs(sheduller gocron.Scheduler) ([]gocron.Job, error) {
	j, err := sheduller.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(
			func() {
				fmt.Println("Hello World")
			},
		),
	)
	if err != nil {
		return nil, err
	}

	return []gocron.Job{j}, nil
}
