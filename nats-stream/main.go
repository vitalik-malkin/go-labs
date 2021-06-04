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

	takeLoopCount, pubLoopCount := 20, 20
	takeLoopExit, pubLoopExit := make(chan error, takeLoopCount), make(chan error, pubLoopCount)
	defer close(takeLoopExit)
	defer close(pubLoopExit)

	timeCtx, timeCtxCancelF := context.WithTimeout(stopCtx, time.Second*10)
	defer timeCtxCancelF()

	for y := 0; y < takeLoopCount; y++ {
		go func(i int) {
			takeLoopExit <- takeLoop(timeCtx, i, logger, false, true)
		}(y)
	}
	for y := 0; y < pubLoopCount; y++ {
		go func(i int) {
			pubLoopExit <- publishLoop(timeCtx, i, logger, false, true)
		}(y)
	}

	errs := []error{}
	waitForMaxErrs := pubLoopCount + takeLoopCount
	select {
	case <-stopCtx.Done():
		break
	case err := <-takeLoopExit:
		if err != nil {
			errs = append(errs, err)
		}
		waitForMaxErrs--
		stopCtxF()
		break
	case err := <-pubLoopExit:
		if err != nil {
			errs = append(errs, err)
		}
		waitForMaxErrs--
		stopCtxF()
		break
	}

	for y := 0; y < waitForMaxErrs; y++ {
		select {
		case err := <-takeLoopExit:
			if err != nil {
				errs = append(errs, err)
			}
			break
		case err := <-pubLoopExit:
			if err != nil {
				errs = append(errs, err)
			}
			break
		}
	}

	exitCode := 0
	for i := 0; i < len(errs); i++ {
		logger.Printf("%v", errs[i])
		exitCode = 1
	}

	os.Exit(exitCode)
}

func takeLoop(ctx context.Context, id int, l *log.Logger, dummy bool, noLogInfoLevel bool) error {
	clientID := fmt.Sprintf("take-client-%d", id)

	sc, err := stan.Connect("cluster00", clientID, stan.NatsURL("nats://10.160.0.17:4222"))
	if err != nil {
		return fmt.Errorf("take-loop %d: error; err: %v", id, err)
	}
	defer func() {
		scCloseErr := sc.Close()
		if scCloseErr != nil {
			if err != nil {
				err = fmt.Errorf("take-loop %d: %v; %w", id, err, scCloseErr)
			} else {
				err = fmt.Errorf("take-loop %d: error; err: %v", id, scCloseErr)
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
			return fmt.Errorf("take-loop %d: error; err: %v", id, err)
		}
		defer func() {
			subCloseErr := sub.Close()
			if subCloseErr != nil {
				if err != nil {
					err = fmt.Errorf("take-loop %d: %v; %w", id, err, subCloseErr)
				} else {
					err = fmt.Errorf("take-loop %d: error; err: %v", id, subCloseErr)
				}
			}
		}()
	}

	ackedCount := 0
	startDT := time.Now()

	for {
		select {
		case <-ctxDone:
			l.Printf("take-loop %d: interrupted by context signal; acked: %d; duration: %s", id, ackedCount, time.Now().Sub(startDT))
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
			err = fmt.Errorf("take-loop %d: subscription error received; err: %v", id, subErrObj)
			return err
		case <-tmr.C:
			if !noLogInfoLevel {
				l.Printf("take-loop: no data")
			}
			tmr.Reset(interval)
		}
	}
}

func publishLoop(ctx context.Context, id int, l *log.Logger, dummy bool, noLogInfoLevel bool) (err error) {
	clientID := fmt.Sprintf("pub-client-%d", id)

	sc, err := stan.Connect("cluster00", clientID, stan.NatsURL("nats://10.160.0.17:4222"))
	if err != nil {
		return fmt.Errorf("pub-loop %d: error; err: %v", id, err)
	}
	defer func() {
		scCloseErr := sc.Close()
		if scCloseErr != nil {
			if err != nil {
				err = fmt.Errorf("pub-loop %d: %v; %w", id, err, scCloseErr)
			} else {
				err = fmt.Errorf("pub-loop %d: error; err: %v", id, scCloseErr)
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
			l.Printf("pub-loop %d: interrupted by context signal; published: %d; duration: %s", id, pubCount, time.Now().Sub(startDT))
			return err
		default:
			if !dummy {
				pubDataBin := []byte(fmt.Sprintf("%s", time.Now()))
				err = sc.Publish("subject00", pubDataBin)
				if err != nil {
					return fmt.Errorf("pub-loop %d: error; err: %v", id, err)
				}
				pubCount++
				if !noLogInfoLevel {
					l.Printf("pub-loop: published data; data: %s", string(pubDataBin))
				}
			}
			//tmr.Reset(interval)
		}
	}
}
