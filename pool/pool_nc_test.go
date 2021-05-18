package pool_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/happyxhw/gopkg/pool"
)

func dialer(_ context.Context) (net.Conn, error) {
	fmt.Println("dial")
	cn, err := net.Dial("tcp", "127.0.0.1:8001")
	return cn, err
}

// ncat -l --keep-open -p 8001
func TestNcat(t *testing.T) {
	cn, err := net.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		panic(err)
	}

	_, err = cn.Write([]byte("123\n"))
	if err != nil {
		panic(err)
	}

	_, err = cn.Write([]byte("456\n"))
	if err != nil {
		panic(err)
	}
}

// ncat -l --keep-open -p 8001
func TestNcatPool(t *testing.T) {
	connPool := pool.NewConnPool(&pool.Options{
		Dialer: dialer,
		OnClose: func(conn *pool.Conn) error {
			fmt.Println("close: ", conn.CreatedAt())
			return nil
		},
		PoolSize:           100,
		PoolTimeout:        time.Second * 30,
		IdleTimeout:        time.Second * 10,
		IdleCheckFrequency: time.Second * 5,
		MinIdleConns:       100,
		MaxConnAge:         time.Minute * 10,
	})
	go func() {
		for {
			fmt.Println(connPool.Len(), connPool.IdleLen())
			fmt.Printf("stats: %+v\n", connPool.States())
			time.Sleep(time.Second)
		}
	}()

	for i := 0; i < 1; i++ {
		go func() {
			for {
				cn, err := connPool.Get(context.Background())
				if err != nil {
					panic(err)
				}
				_, err = cn.Write([]byte("123\n"))
				if err != nil {
					panic(err)
				}
				time.Sleep(time.Second)
				if err := connPool.Put(context.Background(), cn); err != nil {
					panic(err)
				}
			}
		}()
	}

	time.Sleep(time.Second * 1000)
}
