package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/mikelong/usersvc"
	"golang.org/x/net/context"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", "localhost:8080", "HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)
		logger = log.NewContext(logger).With("caller", log.DefaultCaller)
	}

	errs := make(chan error)

	ctx := context.Background()
	svc := usersvc.NewUserService()

	var h http.Handler
	{
		h = usersvc.MakeHttpHandler(ctx, svc)
	}

	fmt.Printf(*httpAddr)

	go func() {
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("exit", <-errs)

	fmt.Printf("FOO")
}
