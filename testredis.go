// Package testredis provides an embedded Redis instance for unit testing
package testredis

import (
	"fmt"
	"io/ioutil"
	"os"

	lediscfg "github.com/siddontang/ledisdb/config"
	"github.com/siddontang/ledisdb/server"

	"gopkg.in/redis.v5"
)

// Redis is a testing Redis
type Redis interface {
	// Addr gets the address of this Redis
	Addr() string

	// Client opens a new client to this Redis
	Client() *redis.Client

	// Start starts an unstarted Redis instance
	Start()

	// Close closes this Redis and associated resources
	Close() error
}

type testRedis struct {
	dir string
	app *server.App
}

func (tr *testRedis) Addr() string {
	return tr.app.Address()
}

func (tr *testRedis) Start() {
	go tr.app.Run()
}

func (tr *testRedis) Client() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: tr.app.Address(),
	})
}

func (tr *testRedis) Close() error {
	tr.app.Close()
	return os.RemoveAll(tr.dir)
}

// Open opens a new test Redis
func Open() (Redis, error) {
	redis, err := OpenUnstarted()
	if err == nil {
		redis.Start()
	}
	return redis, err
}

// OpenUnstarted opens a new test Redis but not start it.
func OpenUnstarted() (Redis, error) {
	tmpDir, err := ioutil.TempDir("", "redis")
	if err != nil {
		return nil, fmt.Errorf("Unable to create temp dir: %v", err)
	}

	cfg := lediscfg.NewConfigDefault()
	cfg.Addr = "" // use auto-assigned address
	cfg.DataDir = tmpDir

	app, err := server.NewApp(cfg)
	if err != nil {
		return nil, err
	}
	return &testRedis{tmpDir, app}, nil
}
