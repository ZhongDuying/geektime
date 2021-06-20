package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/*

测试1: 调用shutdown方式结束

[zhongdy@localhost week03]$ ./main
http handle: Hello Client!
http handle: Hello Client!
http handle: Shutdown Now!
http handle [signal] exit.
http handle [shutdown] exit.
main: http: Server closed

测试2: 通过接收信号方式结束

[zhongdy@localhost week03]$ ./main
http handle: Hello Client!
http handle: Hello Client!
^Csignal received:interrupt
http handle [shutdown] exit.
http handle [signal] exit.
main: http: Server closed
[zhongdy@localhost week03]$

*/


func main() {
	muxHttp := http.NewServeMux()

	// 处理hello请求
	muxHttp.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("http handle: Hello Client!")
		_, err := w.Write([]byte("Hello Client!"))
		if err != nil {
			return
		}
	})

	// 处理shutdown请求
	chShutdown := make(chan struct{})
	muxHttp.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("http handle: Shutdown Now!")
		chShutdown <- struct{}{}
	})

	// 创建http Server及context
	server := http.Server{
		Addr: ":8081",
		Handler: muxHttp,
	}
	g, ctx := errgroup.WithContext(context.Background())

	// 启动http Server
	g.Go(func() error {
		return server.ListenAndServe()
	})

	// 处理shutdown消息
	g.Go(func() error {
		defer fmt.Println("http handle [shutdown] exit.")
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-chShutdown:
			ctxGracefulExit, cancel := context.WithTimeout(context.Background(), time.Second * 2)
			defer cancel()
			return server.Shutdown(ctxGracefulExit)
		}
	})

	// 处理系统信号
	g.Go(func() error {
		defer fmt.Println("http handle [signal] exit.")
		chSignal := make(chan os.Signal)
		signal.Notify(chSignal, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)	// 响应信号：中断、程序退出
		select {
		case sig := <- chSignal:
			sigInfo := sig.String()
			fmt.Println("signal received:" + sigInfo)

			ctxGracefulExit, cancel := context.WithTimeout(context.Background(), time.Second * 2)
			defer cancel()
			return server.Shutdown(ctxGracefulExit)
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	// goroutine 生命周期管理
	if err := g.Wait(); err != nil {
		fmt.Printf("main: %v\n", err)
		return
	}
}

