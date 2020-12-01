package dispatcher

import (
	"fmt"
	"testing"
	"time"
)

func TestNewDispatcher(t *testing.T) {
	d := NewDispatcher(10, -1)

	go func() {
		for {
			x, ok := <-d.ResultCh()
			if ok {
				fmt.Println(x)
			} else {
				fmt.Println("end")
				return
			}
		}
	}()

	for i := 0; i < 30; i++ {
		y := i
		err := d.Send(func() (interface{}, error) {
			time.Sleep(time.Second * 3)
			return y, nil
		})
		if err == ErrStopped {
			fmt.Println(err)
		}
	}

	fmt.Println("debug")

	d.Stop()
}
