package dispatcher

import (
	"fmt"
	"testing"
	"time"
)

func TestNewDispatcher(t *testing.T) {
	d := NewDispatcher(10, time.Second*15)

	go func() {
		for i := 0; i < 80; i++ {
			y := i
			err := d.Send(func() (interface{}, error) {
				time.Sleep(time.Second * 1)
				return y, nil
			})
			if err == ErrStopped {
				fmt.Println(err)
			}
		}
	}()

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
	time.Sleep(time.Second * 30)
	d.Stop()
	time.Sleep(time.Second * 300)

}
