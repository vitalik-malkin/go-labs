package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/stan.go"
	_ "github.com/nats-io/stan.go"
)

func main() {
	logger := log.Default()
	stopCtx, stopCtxF := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGABRT)
	defer stopCtxF()

	takeLoopExit, pubLoopExit := make(chan error, 1), make(chan error, 1)
	defer close(takeLoopExit)
	defer close(pubLoopExit)

	timeCtx, timeCtxCancelF := context.WithTimeout(stopCtx, time.Second*30)
	defer timeCtxCancelF()

	go func() {
		takeLoopExit <- takeLoop(timeCtx, logger, false, true)
	}()
	go func() {
		pubLoopExit <- publishLoop(timeCtx, logger, false, true)
	}()

	var (
		errTakeLoop, errPubLoop error
	)
	select {
	case <-stopCtx.Done():
		errTakeLoop = <-takeLoopExit
		errPubLoop = <-pubLoopExit
		break
	case errTakeLoop = <-takeLoopExit:
		timeCtxCancelF()
		stopCtxF()
		errPubLoop = <-pubLoopExit
		break
	case errPubLoop = <-pubLoopExit:
		timeCtxCancelF()
		stopCtxF()
		errTakeLoop = <-takeLoopExit
		break
	}

	if errTakeLoop != nil {
		logger.Printf("take-loop: error; err: %v", errTakeLoop)
	}
	if errPubLoop != nil {
		logger.Printf("pub-loop: error; err: %v", errPubLoop)
	}

	if errPubLoop != nil || errTakeLoop != nil {
		os.Exit(1)
	}
}

func takeLoop(ctx context.Context, l *log.Logger, dummy bool, noLogInfoLevel bool) error {
	sc, err := stan.Connect("cluster00", "client01", stan.NatsURL("nats://10.160.0.17:4222"))
	if err != nil {
		return err
	}
	defer func() {
		scCloseErr := sc.Close()
		if scCloseErr != nil {
			if err != nil {
				err = fmt.Errorf("%v; %w", err, scCloseErr)
			} else {
				err = scCloseErr
			}
		}
	}()

	interval := time.Millisecond * 1
	tmr := time.NewTimer(interval)
	defer tmr.Stop()

	subInboxChan := make(chan *stan.Msg, 1)
	subErrChan := make(chan error, 1)
	subAcks := make(chan struct{})
	ctxDone := ctx.Done()

	if !dummy {
		sub, err := sc.QueueSubscribe(
			"subject00",
			"subject00_group000",
			func(msg *stan.Msg) {
				select {
				case <-ctxDone:
					return
				case subInboxChan <- msg:
					break
				}

				select {
				case <-ctxDone:
					return
				case <-subAcks:
					err := msg.Ack()
					if err != nil {
						select {
						case <-ctxDone:
							break
						case subErrChan <- fmt.Errorf("ack error: %w", err):
							break
						}
						return
					}
					break
				}
			},
			stan.SetManualAckMode(),
			stan.AckWait(time.Second*10),
			stan.DurableName("subject00--AA"))
		if err != nil {
			return err
		}
		defer func() {
			subCloseErr := sub.Close()
			if subCloseErr != nil {
				if err != nil {
					err = fmt.Errorf("%v; %w", err, subCloseErr)
				} else {
					err = subCloseErr
				}
			}
		}()
	}

	ackedCount := 0
	startDT := time.Now()

	for {
		select {
		case <-ctxDone:
			l.Printf("take-loop: interrupted by context signal; acked: %d; duration: %s", ackedCount, time.Now().Sub(startDT))
			return err
		case msg := <-subInboxChan:
			if !noLogInfoLevel {
				l.Printf("take-loop: got data; data: %s; seq: %d", string(msg.Data), msg.Sequence)
			}
			select {
			case <-ctxDone:
				return err
			case subAcks <- struct{}{}:
				ackedCount++
				if !noLogInfoLevel {
					l.Printf("take-loop: ack data; seq: %d", msg.Sequence)
				}
				break
			}
		case subErrObj := <-subErrChan:
			err = fmt.Errorf("take-loop: subscription error received; err: %v", subErrObj)
			return err
		case <-tmr.C:
			if !noLogInfoLevel {
				l.Printf("take-loop: no data")
			}
			tmr.Reset(interval)
		}
	}
}

func publishLoop(ctx context.Context, l *log.Logger, dummy bool, noLogInfoLevel bool) (err error) {
	sc, err := stan.Connect("cluster00", "client00", stan.NatsURL("nats://10.160.0.17:4222"))
	if err != nil {
		return err
	}
	defer func() {
		scCloseErr := sc.Close()
		if scCloseErr != nil {
			if err != nil {
				err = fmt.Errorf("%v; %w", err, scCloseErr)
			} else {
				err = scCloseErr
			}
		}
	}()

	interval := time.Millisecond * 1
	tmr := time.NewTimer(interval)
	defer tmr.Stop()

	pubCount := 0
	startDT := time.Now()
	ctxDone := ctx.Done()

	for {
		select {
		case <-ctxDone:
			l.Printf("pub-loop: interrupted by context signal; published: %d; duration: %s", pubCount, time.Now().Sub(startDT))
			return nil
		default:
			if !dummy {
				pubDataBin := []byte(fmt.Sprintf("%s", time.Now()))
				err = sc.Publish("subject00", pubDataBin)
				if err != nil {
					return err
				}
				pubCount++
				if !noLogInfoLevel {
					l.Printf("pub-loop: published data; data: %s", string(pubDataBin))
				}
			}
			tmr.Reset(interval)
		}
	}
}
