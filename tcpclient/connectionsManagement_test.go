package tcpclient

import (
	"strconv"
	"sync"
	"testing"
	"time"
	"github.com/dachad/check-max-tcp-connections/tcpserver"
	"io/ioutil"
)

// We really need to refactor this test. We should verify connections do become established,
// rather than just waiting for a second and finish
// We should also test "failing" connections, and ensure their status is reported properly
func TestTcpConnect(t *testing.T) {
	var numberConnections = 2
	var host = "127.0.0.1"
	var port = 55555

	dispatcher := &tcpserver.Dispatcher{make(map[string]*tcpserver.Handler)}

	run := func() {
		if err := dispatcher.ListenHandlers(port); err != nil {
			t.Error("Could not start the TCP server", err)
			return
		}
		t.Log("TCP server started")
	}
	go run()

	var wg sync.WaitGroup
	wg.Add(numberConnections)

	for runner := 1; runner <= numberConnections; runner++ {
		t.Log("Initiating runner # ", strconv.Itoa(runner))
		go TCPConnect(runner, host, port, &wg, ioutil.Discard, make(chan Connection, numberConnections), make(chan bool))
		t.Logf("Runner %s initated. Remaining: %s", strconv.Itoa(runner), strconv.Itoa(numberConnections-runner))
	}

	t.Log("Waiting runners to finish")
	time.Sleep(time.Second)

	// Marking wait groups as done after a second does not make sense...
	// and in case one client finishes before this loop being executed (because of an error?)
	// we will get an exception (panic: sync: negative WaitGroup counter [recovered])

	//for runner := 1; runner <= numberConnections; runner++ {
	//	t.Log("Closing runner #", strconv.Itoa(runner))
	//	wg.Done()
	//}

}