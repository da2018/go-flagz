// Copyright 2015 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package main

import (
	"log"
	"os"
	"time"

	etcd "github.com/coreos/etcd/client"
	"github.com/da2018/go-flagz"
	"github.com/da2018/go-flagz/watcher"
	flag "github.com/spf13/pflag"
)

var (
	myFlagSet = flag.NewFlagSet("custom_flagset", flag.ContinueOnError)

	staticInt = myFlagSet.Int32("example_my_static_int", 8080, "some static int int")

	dynStr = flagz.DynString(myFlagSet, "example_my_dynamic_string", "initial_value", "someusage")
	dynInt = flagz.DynInt64(myFlagSet, "example_my_dynamic_int", 1337, "someusage int")
)

func main() {
	myFlagSet.Parse(os.Args[1:])
	logger := log.New(os.Stderr, "wr ", log.LstdFlags)

	client, err := etcd.New(etcd.Config{Endpoints: []string{"http://localhost:2379"}})
	if err != nil {
		logger.Fatalf("Failed setting up %v", err)
	}
	w, err := watcher.New(myFlagSet, etcd.NewKeysAPI(client), "/example/flagz", logger)
	if err != nil {
		logger.Fatalf("Failed setting up %v", err)
	}
	err = w.Initialize()
	if err != nil {
		logger.Fatalf("Failed setting up %v", err)
	}
	w.Start()

	for true {
		logger.Printf("staticint: %v dynint: %v dynstring: %v",
			*staticInt,
			dynInt.Get(),
			dynStr.Get())
		time.Sleep(1500 * time.Millisecond)
	}
}
