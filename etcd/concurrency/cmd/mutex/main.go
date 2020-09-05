package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	etcd_cli "github.com/coreos/etcd/clientv3"
	etcd_conc "github.com/coreos/etcd/clientv3/concurrency"
	clr "gopkg.in/gookit/color.v1"
)

type opCode int

const (
	lockOpCode   opCode = 1
	unlockOpCode opCode = 2
)

type lockState struct {
	name string
}

type opState struct {
	code opCode
	lock *lockState
	ctx  context.Context

	resp chan error
}

type clientState struct {
	name string

	ctx context.Context
	op  chan opState
}

var cli *etcd_cli.Client

func main() {
	mainCtx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sig
		cancel()
	}()
	var err error

	cli, err = etcd_cli.New(
		etcd_cli.Config{
			Endpoints: []string{"localhost:2379"},
			Context:   mainCtx,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	l1 := clientState{
		name: "LOOP-1",
		ctx:  mainCtx,
		op:   make(chan opState),
	}
	go lockLoop(l1)

	l2 := clientState{
		name: "LOOP-2",
		ctx:  mainCtx,
		op:   make(chan opState),
	}
	go lockLoop(l2)

	go workLoop(l1)
	go workLoop(l2)

	<-mainCtx.Done()
	fmt.Println("\nexiting...")
}

func workLoop(cl clientState) {
	for {
		select {
		case <-cl.ctx.Done():
			return
		default:
			time.Sleep(time.Second)

			err := placeLock(cl, "lck1")
			if err != nil {
				clr.Error.Printf("%s (WORK): %v\n", cl.name, err)
				continue
			}

			clr.Info.Printf("%s (WORK): doing work...\n", cl.name)
			time.Sleep(10 * time.Second)

			err = releaseLock(cl, "lck1")
			if err != nil {
				clr.Error.Printf("%s (WORK): %v\n", cl.name, err)
			}

		}
	}
}

func placeLock(cl clientState, name string) (err error) {
	op := opState{
		code: lockOpCode,
		lock: &lockState{
			name: name,
		},
		resp: make(chan error, 1),
	}
	cl.op <- op
	err = <-op.resp

	return
}

func releaseLock(cl clientState, name string) (err error) {
	op := opState{
		code: unlockOpCode,
		lock: &lockState{
			name: name,
		},
		resp: make(chan error, 1),
	}
	cl.op <- op
	err = <-op.resp

	return
}

func lockLoop(cl clientState) error {
	session, err := etcd_conc.NewSession(
		cli,
	)
	if err != nil {
		clr.Error.Printf("error while creating session for '%s': %v", cl.name, err)
		return err
	}
	defer session.Close()

	clr.Info.Printf("%s: started\n", cl.name)

	muses := make(map[string]*etcd_conc.Mutex)

	for {
		select {
		case <-cl.ctx.Done():
			return cl.ctx.Err()
		case op := <-cl.op:
			muPfx := fmt.Sprintf("/%s/", op.lock.name)

			switch op.code {
			case lockOpCode:
				mu, ok := muses[muPfx]
				if !ok {
					mu = etcd_conc.NewMutex(session, muPfx)
					muses[muPfx] = mu
				}
				err := mu.Lock(cl.ctx)
				if err == nil {
					clr.Info.Printf("%s: lock('%s'), SUCCESS\n", cl.name, op.lock.name)
				} else {
					clr.Error.Printf("%s: lock('%s'), FAIL: %v\n", cl.name, op.lock.name, err)
				}

				op.resp <- err
			case unlockOpCode:
				var err error
				mu, ok := muses[muPfx]
				if !ok {
					err = fmt.Errorf("consistency violation: no mutex found by name '%s'", muPfx)
				}
				if err == nil {
					err = mu.Unlock(cl.ctx)
				}
				if err == nil {
					clr.Info.Printf("%s: unlock('%s'), SUCCESS\n", cl.name, op.lock.name)
				} else {
					clr.Error.Printf("%s: unlock('%s'), FAIL: %v\n", cl.name, op.lock.name, err)
				}
				op.resp <- err
			default:
				op.resp <- fmt.Errorf("unknown op code")
			}
		}

	}

}
