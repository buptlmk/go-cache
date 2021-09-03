package timewheel

import (
	"errors"
	"fmt"
	"go-cache/internal/linkedlist"
	"go-cache/internal/syncx"
	"sync"
	"time"
)

const buffer = 16

type TimeWheel struct {
	interval time.Duration
	ticker   *time.Ticker
	slots    []*linkedlist.List
	// 这个string得保持唯一，对于task的key 在添加前先进行判断
	tasks map[string]*linkedlist.Item
	mutex sync.RWMutex

	currentPos int
	slotNum    int

	addTaskChanel     chan *task
	removeTaskChannel chan string
	stopChannel       chan struct{}
}

type task struct {
	delay time.Duration
	key   string
	job   func()

	circle int
	pos    int
}

func NewTimeWheel(interval time.Duration, slotNum int) *TimeWheel {

	slots := make([]*linkedlist.List, slotNum)

	for i := 0; i < slotNum; i++ {
		slots[i] = linkedlist.NewList()
	}

	return &TimeWheel{
		interval:          interval,
		slots:             slots,
		tasks:             make(map[string]*linkedlist.Item),
		slotNum:           slotNum,
		addTaskChanel:     make(chan *task, buffer),
		removeTaskChannel: make(chan string, buffer),
		stopChannel:       make(chan struct{}),
	}
}

func (t *TimeWheel) Start() {
	t.ticker = time.NewTicker(t.interval)
	go func() {
		for {
			select {
			case <-t.ticker.C:
				// 定时处理
				t.handleTask()
			case v, ok := <-t.removeTaskChannel:
				// 删除某些任务
				if ok {
					t.removeTask(v)
				} else {
					return
				}
			case v, ok := <-t.addTaskChanel:
				// 添加定时任务
				if ok {
					t.addTask(v)
				} else {
					return
				}

			case <-t.stopChannel:
				t.ticker.Stop()
				return
			}
		}
	}()
}
func (t *TimeWheel) AddJob(key string, delayed time.Duration, fn func()) error {
	t.mutex.RLock()
	_, ok := t.tasks[key]
	t.mutex.RUnlock()
	if ok {
		return errors.New("the " + key + " job is exist,please replace your job name")
	}

	task := &task{
		key:   key,
		delay: delayed,
		job:   fn,
	}

	t.addTaskChanel <- task
	return nil
}

func (t *TimeWheel) RemoveJob(key string) {
	t.mutex.RLock()
	item, ok := t.tasks[key]
	t.mutex.RUnlock()
	if !ok {
		fmt.Println("dd")
		return
	}
	tas := item.Value.(*task)
	t.slots[tas.pos].Remove(item)
}

// 计时器定时stop
func (t *TimeWheel) Stop() {
	t.stopChannel <- struct{}{}
}

func (t *TimeWheel) addTask(v *task) {

	// 首先根据当前位置确定应该把task 放在那个时间槽中
	pos, circle := t.getPosition(v.delay)
	v.circle = circle
	v.pos = pos
	item := &linkedlist.Item{
		Value: v,
	}
	l := t.slots[pos]
	l.Push(item)
	t.mutex.Lock()
	t.tasks[v.key] = item
	t.mutex.Unlock()
}

func (t *TimeWheel) removeTask(key string) {
	t.mutex.RLock()
	item, ok := t.tasks[key]
	t.mutex.RUnlock()

	if !ok {
		return
	}
	taskValue := item.Value.(*task)
	l := t.slots[taskValue.pos]
	l.Remove(item)
}

func (t *TimeWheel) handleTask() {
	t.execTask(t.slots[t.getCurrentPos()])

	if t.currentPos == t.slotNum-1 {
		t.currentPos = 0
	} else {
		t.currentPos++
	}
}

func (t *TimeWheel) execTask(l *linkedlist.List) {

	for e := l.Front(); e != nil; {
		temp := e.Value.(*task)
		if temp.circle > 0 {
			temp.circle--
			e = e.Next
			continue
		}

		syncx.SafedGroutine(temp.job, nil)
		l.Remove(e)
		e = e.Next
		t.mutex.Lock()
		delete(t.tasks, temp.key)
		t.mutex.Unlock()
	}
}

func (t *TimeWheel) getPosition(d time.Duration) (pos int, circle int) {
	n := d.Milliseconds()

	millSec := t.interval.Milliseconds()

	circle = int(n/millSec) / t.slotNum
	pos = (int(n/millSec) + t.getCurrentPos()) % t.slotNum
	return

}

func (t *TimeWheel) getCurrentPos() int {
	// 最好的应该是 atomic.load;
	return t.currentPos
}
