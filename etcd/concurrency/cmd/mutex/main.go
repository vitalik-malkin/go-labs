package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
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

	stop chan struct{}
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

	clr.Info.Printf("connecting to etcd...\n")
	cli, err = etcd_cli.New(
		etcd_cli.Config{
			Endpoints: []string{"localhost:2379"},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	cl1 := clientState{
		name: "C-1",
		ctx:  mainCtx,
		op:   make(chan opState),
		stop: make(chan struct{}, 1),
	}
	go lockLoop(cl1)

	cl2 := clientState{
		name: "C-2",
		ctx:  mainCtx,
		op:   make(chan opState),
		stop: make(chan struct{}, 1),
	}
	go lockLoop(cl2)

	go workLoop(cl1)
	go workLoop(cl2)

	<-mainCtx.Done()
	clr.FgBlue.Println("\nexiting...")
	<-cl1.stop
	clr.FgBlue.Printf("%s: done\n", cl1.name)
	<-cl2.stop
	clr.FgBlue.Printf("%s: done\n", cl2.name)
}

func workLoop(cl clientState) {
	lckName := "LCK_2"
	for {
		select {
		case <-cl.ctx.Done():
			return
		default:
			time.Sleep(time.Second)

			timeoutMs := int64(rand.Float32()*3000.0) + 1000
			err := placeLock(cl.ctx, cl, lckName, timeoutMs)
			if err != nil {
				clr.Error.Printf("%s (WORK): timeout=%d, %v\n", cl.name, timeoutMs, err)
				continue
			}

			timeoutMs = int64(rand.Float32()*7000.0) + 3000
			clr.Yellow.Printf("%s (WORK): duration=%d, doing work...\n", cl.name, timeoutMs)
			time.Sleep(10 * time.Second)

			err = releaseLock(cl, lckName)
			if err != nil {
				clr.Error.Printf("%s (WORK): %v\n", cl.name, err)
			}

		}
	}
}

func placeLock(ctx context.Context, cl clientState, name string, timeoutMs int64) (err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*time.Duration(timeoutMs))
	defer cancel()

	op := opState{
		code: lockOpCode,
		lock: &lockState{
			name: name,
		},
		ctx:  ctx,
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
		ctx:  context.Background(),
		resp: make(chan error, 1),
	}
	cl.op <- op
	err = <-op.resp

	return
}

func lockLoop(cl clientState) error {
	defer func() {
		cl.stop <- struct{}{}
	}()

	session, err := etcd_conc.NewSession(
		cli,
	)
	if err != nil {
		clr.Error.Printf("error while creating session for '%s': %v", cl.name, err)
		return err
	}
	defer func() {
		session.Close()
		clr.FgBlue.Printf("%s: session closed\n", cl.name)
	}()

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
				err := mu.Lock(op.ctx)
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
					err = mu.Unlock(op.ctx)
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
