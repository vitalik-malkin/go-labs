package main

import (
	"context"
	crnd "crypto/rand"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"unsafe"

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
	name       string
	lockTTLSec int

	ctx context.Context
	op  chan opState

	stop chan struct{}
}

var cli *etcd_cli.Client

func main() {
	var (
		clientNameFlag string
		lockTTLSecFlag int
	)
	flag.StringVar(&clientNameFlag, "client-name", "", "client name")
	flag.IntVar(&lockTTLSecFlag, "lock-ttl", -1, "lock time to live period")
	flag.Parse()

	clientName := strings.ToUpper(clientNameFlag)
	lockTTLSec := lockTTLSecFlag

	if len(clientName) < 1 {
		flag.Usage()
		log.Fatal("no client name specified")
	} else if lockTTLSec < 0 {
		flag.Usage()
		log.Fatal("lock time to live period cannot be less than 0")
	}

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

	cl := clientState{
		name:       clientName,
		lockTTLSec: lockTTLSec,
		ctx:        mainCtx,
		op:         make(chan opState),
		stop:       make(chan struct{}, 1),
	}
	go lockLoop(cl)

	go workLoop(cl)

	<-mainCtx.Done()
	clr.FgBlue.Println("\nexiting...")
	<-cl.stop
	clr.FgBlue.Printf("%s: done\n", cl.name)

	clr.Info.Println("Bye...")
}

func nextRnd(min, maxExclusive int32) (val int32) {
	rangeSize := int64(maxExclusive) - int64(min)
	if rangeSize < 1 {
		panic(fmt.Errorf("invalid range [%d; %d)", min, maxExclusive))
	} else if rangeSize == 1 {
		return min
	}

	bufArr := [4]byte{}
	buf := bufArr[:]
	_, err := crnd.Read(buf)
	if err != nil {
		panic(fmt.Errorf("error while using crypto-rand: %w", err))
	}
	genVal := *(*uint32)(unsafe.Pointer(&bufArr))
	genLen := int(math.Ceil(math.Log10(float64(genVal))))
	mul := float64(genVal) * math.Pow10(-genLen)
	addition := int32(float64(rangeSize) * mul)
	return min + addition
}

func workLoop(cl clientState) {
	lckName := "LCK_2"
	for {
		select {
		case <-cl.ctx.Done():
			return
		default:
			time.Sleep(time.Second)

			timeoutMs := int64(nextRnd(1000, 5000))
			err := placeLock(cl.ctx, cl, lckName, timeoutMs)
			if err != nil {
				clr.Error.Printf("%s (WORK): timeout=%d, %v\n", cl.name, timeoutMs, err)
				continue
			}

			timeoutMs = int64(nextRnd(12000, 20000))
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
		etcd_conc.WithTTL(cl.lockTTLSec),
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

	isSessionInvalid := false
	muses := make(map[string]*etcd_conc.Mutex)

	for {
		select {
		case <-cl.ctx.Done():
			return cl.ctx.Err()
		case <-session.Done():
			isSessionInvalid = true
		case op := <-cl.op:
			opErr := error(nil)
			if isSessionInvalid {
				opErr = fmt.Errorf("invalid session")
			} else {
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
				default:
					opErr = fmt.Errorf("unknown op code")
				}
			}
			op.resp <- opErr
		}

	}
}
