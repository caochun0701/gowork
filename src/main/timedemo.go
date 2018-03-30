package main

import (
	"fmt"
	"time"
	"sync/atomic"
	"container/list"
	"sync"
	"github.com/go-lumber/log"
)
var sum int32 = 0
var N int32 = 300

var tt *Timer

const (
	TIME_NEAR_SHIFT  = 8
	TIME_NEAR        = 1 << TIME_NEAR_SHIFT
	TIME_LEVEL_SHIFT = 6
	TIME_LEVEL       = 1 << TIME_LEVEL_SHIFT
	TIME_NEAR_MASK   = TIME_NEAR - 1
	TIME_LEVEL_MASK  = TIME_LEVEL - 1
)

type Timer struct {
	near [TIME_NEAR]*list.List
	t    [4][TIME_LEVEL]*list.List
	sync.Mutex
	time uint32
	tick time.Duration
	quit chan struct{}
}
func (t *Timer) Stop() {
	close(t.quit)
}

func now()  {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	atomic.AddInt32(&sum, 1)
	v := atomic.LoadInt32(&sum)
	if v == 2*N{
		tt.Stop()
	}
}

func New(d time.Duration) *Timer {
	t := new(Timer)
	t.time = 0
	t.tick = d
	t.quit = make(chan struct{})
	var i, j int
	for i = 0; i < TIME_NEAR; i++ {
		t.near[i] = list.New()
	}

	for i = 0; i < 4; i++ {
		for j = 0; j < TIME_LEVEL; j++ {
			t.t[i][j] = list.New()
		}
	}
	return t
}

type Node struct {
	expire uint32
	f      func()
}

func (t *Timer) Start() {
	tick := time.NewTicker(t.tick)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			t.update()
		case <-t.quit:
			return
		}
	}
}
func (t *Timer) execute() {
	t.Lock()
	idx := t.time & TIME_NEAR_MASK
	vec := t.near[idx]
	if vec.Len() > 0 {
		front := vec.Front()
		vec.Init()
		t.Unlock()
		dispatchList(front)
		return
	}
	t.Unlock()
}
func (t *Timer) shift() {
	t.Lock()
	var mask uint32 = TIME_NEAR
	t.time++
	ct := t.time
	if ct == 0 {
		t.moveList(3, 0)
	} else {
		time := ct >> TIME_NEAR_SHIFT
		var i int = 0
		for (ct & (mask - 1)) == 0 {
			idx := int(time & TIME_LEVEL_MASK)
			if idx != 0 {
				t.moveList(i, idx)
				break
			}
			mask <<= TIME_LEVEL_SHIFT
			time >>= TIME_LEVEL_SHIFT
			i++
		}
	}
	t.Unlock()
}

func dispatchList(front *list.Element) {
	for e := front; e != nil; e = e.Next() {
		node := e.Value.(*Node)
		go node.f()
	}
}

func (t *Timer) moveList(level, idx int) {
	vec := t.t[level][idx]
	front := vec.Front()
	vec.Init()
	for e := front; e != nil; e = e.Next() {
		node := e.Value.(*Node)
		t.addNode(node)
	}
}
func (t *Timer) update() {
	t.execute()
	t.shift()
	t.execute()
}
func (t *Timer) NewTimer(d time.Duration, f func()) *Node {
	n := new(Node)
	n.f = f
	t.Lock()
	n.expire = uint32(d/t.tick) + t.time
	t.addNode(n)
	t.Unlock()
	return n
}
func (t *Timer) addNode(n *Node) {
	expire := n.expire
	current := t.time
	if (expire | TIME_NEAR_MASK) == (current | TIME_NEAR_MASK) {
		t.near[expire&TIME_NEAR_MASK].PushBack(n)
	} else {
		var i uint32
		var mask uint32 = TIME_NEAR << TIME_LEVEL_SHIFT
		for i = 0; i < 3; i++ {
			if (expire | (mask - 1)) == (current | (mask - 1)) {
				break
			}
			mask <<= TIME_LEVEL_SHIFT
		}

		t.t[i][(expire>>(TIME_NEAR_SHIFT+i*TIME_LEVEL_SHIFT))&TIME_LEVEL_MASK].PushBack(n)
	}

}
func main() {
	//定义一个十分之一毫秒
	timer := New(time.Millisecond * 10)
	//赋值
	tt = timer
	fmt.Println(timer)
	var i int32
	//循环300次
	for i = 0; i < N; i++ {
		timer.NewTimer(time.Millisecond * time.Duration(10*i), now)
		timer.NewTimer(time.Millisecond * time.Duration(10*i), now)
	}
	timer.Start()
	if sum != 2*N {
		log.Println("failed")
	}
	log.Println("this is sum:",sum)
}
