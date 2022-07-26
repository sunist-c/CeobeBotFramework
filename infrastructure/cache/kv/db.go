package kv

import (
	"sync"
	"time"
)

type IDatabase interface {
	Get(key interface{}) (value interface{}, timeout time.Duration, err error)
	Set(key interface{}, value interface{}, timeout time.Duration) (v interface{}, t time.Duration, e error)
	Delete(key interface{}) (ok bool, err error)
	Range(f func(key, value interface{}))
	Close()
	Clear()
}

func NewDatabase(workers, buffSize, dirtyCount uint, hash func(key interface{}) uint) IDatabase {
	// init kv-database
	db := &database{
		buff:    make(chan *dbRequest, buffSize),
		exit:    make(chan struct{}, 1),
		data:    make(map[uint]*dirtyMap, dirtyCount),
		size:    dirtyCount,
		workers: workers,
		hash:    hash,
	}

	// init dirty-map
	for i := uint(0); i < dirtyCount; i++ {
		db.data[i] = &dirtyMap{
			lock:  sync.RWMutex{},
			dirty: make(map[interface{}]*entity),
		}
	}

	// start up operators
	db.operator()
	return db
}

type database struct {
	size    uint
	workers uint
	buff    chan *dbRequest
	exit    chan struct{}
	data    map[uint]*dirtyMap
	hash    func(interface{}) uint
}

func (d *database) Get(key interface{}) (value interface{}, timeout time.Duration, err error) {
	req := newDbRequest(getValue)
	defer close(req.callback)
	req.args = map[requestFields]interface{}{
		requestKey: key,
	}
	d.buff <- req

	resp := <-req.callback
	switch resp.responseType {
	case withError:
		err := resp.returns[responseError].(error)
		return nil, time.Duration(0), err
	default:
		return resp.returns[responseValue], resp.returns[responseTimeout].(time.Duration), nil
	}
}

func (d *database) Set(key interface{}, value interface{}, timeout time.Duration) (v interface{}, t time.Duration, e error) {
	req := newDbRequest(setValue)
	defer close(req.callback)
	req.args = map[requestFields]interface{}{
		requestKey:     key,
		requestValue:   value,
		requestTimeout: timeout,
	}
	d.buff <- req

	resp := <-req.callback
	switch resp.responseType {
	case withError:
		err := resp.returns[responseError].(error)
		return nil, time.Duration(0), err
	default:
		return resp.returns[responseValue], resp.returns[responseTimeout].(time.Duration), nil
	}
}

func (d *database) Delete(key interface{}) (ok bool, err error) {
	req := newDbRequest(deleteKey)
	defer close(req.callback)
	req.args = map[requestFields]interface{}{
		requestKey: key,
	}
	d.buff <- req

	resp := <-req.callback
	switch resp.responseType {
	case withError:
		err := resp.returns[responseError].(error)
		return false, err
	default:
		return true, nil
	}
}

func (d *database) Range(f func(key interface{}, value interface{})) {
	req := newDbRequest(rangeDB)
	defer close(req.callback)
	req.args = map[requestFields]interface{}{
		requestFunc: f,
	}
	d.buff <- req

	<-req.callback
	return
}

func (d *database) Close() {
	d.exit <- struct{}{}
}

func (d *database) Clear() {
	for i := uint(0); i < d.size; i++ {
		d.data[i].lock.Lock()
		d.data[i].dirty = make(map[interface{}]*entity)
		d.data[i].lock.Unlock()
	}
}
