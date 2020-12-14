package idgen

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TraceIdGenerator struct {
	lock      sync.Mutex
	increment int64

	incrementLen    int
	maxIncrement    int64
	incrementFormat string
}

func NewTraceIdGenerator(incrementLen int) *TraceIdGenerator {
	return new(TraceIdGenerator).SetIncrementLen(incrementLen)
}

var DefaultTraceIdGenerator = NewTraceIdGenerator(4)

func (t *TraceIdGenerator) SetIncrementLen(incrementLen int) *TraceIdGenerator {
	t.incrementLen = incrementLen
	t.maxIncrement = int64(math.Pow10(incrementLen))
	t.incrementFormat = "%0" + strconv.Itoa(incrementLen) + "d"

	return t
}

func (t *TraceIdGenerator) GenerateId(ip, port string) ([]byte, error) {
	var id string

	for _, item := range strings.Split(ip, ".") {
		v, err := strconv.Atoi(item)
		if err != nil {
			return nil, err
		}
		id += fmt.Sprintf("%02x", v)
		fmt.Println(id)
	}
	fmt.Println(id)
	id += fmt.Sprintf("%05s", port)
	fmt.Println(id)
	id += strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	fmt.Println(id)
	t.lock.Lock()
	increment := t.increment
	t.increment = (t.increment + 1) % t.maxIncrement
	t.lock.Unlock()

	id += fmt.Sprintf(t.incrementFormat, increment)
	fmt.Println(id)

	return []byte(id), nil
}
