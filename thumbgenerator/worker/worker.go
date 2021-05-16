package worker

import (
	"database/sql"
	"github.com/bearname/videohost/thumbgenerator/model"
	"github.com/bearname/videohost/thumbgenerator/processor"
	"github.com/bearname/videohost/thumbgenerator/provider"
	log "github.com/sirupsen/logrus"
	"sync"
)

func Worker(tasksChan <-chan *model.Task, db *sql.DB, name int) {
	log.Printf("start worker %v\n", name)
	for task := range tasksChan {
		log.Printf("start processing video with id %v on worker %v\n", task.Id, name)
		processor.ProcessTask(task, db)
		log.Printf("end processing video with id %v on worker %v\n", task.Id, name)
	}
	log.Printf("stop worker %v\n", name)
}

const WorkersCount = 3

func WorkerPool(stopChan chan struct{}, db *sql.DB) *sync.WaitGroup {
	var wg sync.WaitGroup
	tasksChan := provider.RunTaskProvider(stopChan, db)
	for i := 0; i < WorkersCount; i++ {
		go func(i int) {
			wg.Add(1)
			Worker(tasksChan, db, i)
			wg.Done()
		}(i)
	}
	return &wg
}