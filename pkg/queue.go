package utils

import "fmt"

// Queue with circular fixed buffer
type Queue struct {
	Buf  []string
	Len  int
	Head int // is in the cell for future pop
	Tail int // is in the cell for future push
}

// CreateQueue returns an instance of Queue
func CreateQueue(size int) *Queue {
	return &Queue{
		Buf: make([]string, size),
	}
}

// Push puts element in the tail of Queue
func (d *Queue) Push(value string) {
	if d.Len == len(d.Buf) {
		if d.Tail == len(d.Buf)-1 {
			d.Buf[d.Tail] = value
			d.Tail = 0
			d.Head = 0
		} else {
			d.Buf[d.Tail] = value
			d.Tail++
			d.Head++
		}
	} else {
		if d.Tail == len(d.Buf)-1 {
			d.Buf[d.Tail] = value
			d.Tail = 0
		} else {
			d.Buf[d.Tail] = value
			d.Tail++
		}
		d.Len++
	}
}

// Pop takes element from the head of Queue
func (d *Queue) Pop() (string, error) {
	if d.Len == 0 {
		return "", fmt.Errorf("Queue is empty")
	}

	var result string
	if d.Head == len(d.Buf)-1 {
		result = d.Buf[d.Head]
		d.Buf[d.Head] = ""
		d.Head = 0
	} else {
		result = d.Buf[d.Head]
		d.Buf[d.Head] = ""
		d.Head++
	}
	d.Len--

	return result, nil
}

// FromHeadToTail returns list of all elements from Head to Tail
func (d *Queue) FromHeadToTail() []string {
	if d.Len == 0 {
		return nil
	}

	var elementsPastCounter int
	result := make([]string, 0)

	for i := d.Head; i != d.Tail || elementsPastCounter != d.Len; i++ {
		if i == len(d.Buf)-1 {
			result = append(result, d.Buf[i])
			i = -1
		} else {
			result = append(result, d.Buf[i])
		}
		elementsPastCounter++
	}

	return result
}
