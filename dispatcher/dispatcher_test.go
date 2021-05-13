package dispatcher

import (
	"fmt"
	"testing"
	"time"
)

func TestNewDispatcher(t *testing.T) {
	d := NewDispatcher(
		WithPoolSize(100),
		WithStopTimeout(time.Second*5),
	)

	go func() {
		for {
			x, ok := <-d.ResultCh()
			if ok {
				fmt.Println(x.Result)
				fmt.Println(x.Error)
			} else {
				fmt.Println("end")
				return
			}
		}
	}()

	for i := 0; i < 30; i++ {
		y := i
		job := func() (interface{}, error) {
			time.Sleep(time.Second * 3)
			return y, nil
		}
		err := d.Send(&Task{Job: job, Timeout: time.Second * 10})
		if err == ErrStopped {
			fmt.Println(err)
		}
	}

	fmt.Println("debug")
	time.Sleep(time.Second * 10)
	err := d.Stop()
	fmt.Println("stop: ", err)
}
