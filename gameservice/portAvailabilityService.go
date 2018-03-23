/*

Static structure to manage port availability, no bearing on what is done,
should hand out ports in any order and recycle them when possible

is shared memory safe
 */

package main

import "sync"

const maxUnusedPortsSize = 120
const basePort = 6001

var unusedPorts [maxUnusedPortsSize]int
var unusedPortsSize int
var portArrayMut sync.Mutex

//true if checked out false or nil if still in use
var checkedOut = make(map[int]bool,maxUnusedPortsSize)
var checkedOutMut sync.Mutex

func initPortService() {
for i := 0; i < maxUnusedPortsSize; i++ {
unusedPorts[i] = basePort + i
}
unusedPortsSize = maxUnusedPortsSize
}

func numPortsAvailable() int {
portArrayMut.Lock()
size := unusedPortsSize
portArrayMut.Unlock()

return size
}

func requestPort() int {
portArrayMut.Lock()
checkedOutMut.Lock()

port := unusedPorts[unusedPortsSize-1]
unusedPortsSize--

checkedOut[port] = true

portArrayMut.Unlock()
checkedOutMut.Unlock()

return port
}

func freePort(portNumber int) {
portArrayMut.Lock()
checkedOutMut.Lock()
checkedOut[portNumber] = false
unusedPorts[unusedPortsSize] = portNumber
unusedPortsSize++
checkedOutMut.Unlock()
portArrayMut.Unlock()
}