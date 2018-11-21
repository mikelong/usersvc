package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/mikelong/usersvc"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	errs := make(chan error)

	svc := usersvc.NewUserService()

	var h http.Handler
	{
		h = usersvc.MakeHttpHandler(svc)
	}

	fmt.Printf("Listening on: %s\n", *httpAddr)

	go func() {
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}
