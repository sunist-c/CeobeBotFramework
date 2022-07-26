package kv

import (
	"time"
)

func (d *database) operator() {
	// start up exit listener
	exitChan := make(chan struct{}, d.workers+1)
	go listenExitChan(d.exit, exitChan, d.workers)

	// start up operators
	for i := uint(0); i < d.workers; i++ {
		go startupOperator(d, nil)
	}
}

func listenExitChan(from, to chan struct{}, workers uint) {
	select {
	case <-from:
		for i := uint(0); i < workers; i++ {
			to <- struct{}{}
		}
	}
}

func startupOperator(d *database, e chan struct{}) {
	for {
		select {
		case request := <-d.buff:
			switch request.requestType {
			case getValue:
				request.callback <- getValueOperate(d, request.args[requestKey])
			case setValue:
				request.callback <- setValueOperate(d, request.args[requestKey], request.args[requestValue], request.args[requestTimeout].(time.Duration))
			case deleteKey:
				request.callback <- deleteKeyOperate(d, request.args[requestKey])
			case rangeDB:
				request.callback <- rangeDbOperate(d, request.args[requestFunc].(func(interface{}, interface{})))
			}
		case <-e:
			return
		}
	}
}

func getValueOperate(d *database, key interface{}) (resp dbResponse) {
	// get dirty-map index and set read-lock
	index := d.hash(key) % d.size
	_map := d.data[index]
	_map.lock.RLock()
	defer _map.lock.RUnlock()

	if v, ok := _map.dirty[key]; !ok {
		// condition: cannot find key in dirty-map
		// return error
		resp = dbResponse{
			responseType: withError,
			returns: map[responseFields]interface{}{
				responseError: NotFoundError{key: key},
			},
		}
	} else {
		// condition: get entity succeed in dirty-map
		if v.timeout.Before(time.Now()) {
			// condition: key-value pair timeout
			// delete key-value pair
			_map.lock.Lock()
			defer _map.lock.Unlock()
			delete(_map.dirty, key)

			// return error
			resp = dbResponse{
				responseType: withError,
				returns: map[responseFields]interface{}{
					responseError: NotFoundError{key: key},
				},
			}
		} else {
			// condition: key-value pair not timeout
			// return value
			resp = dbResponse{
				responseType: withValue,
				returns: map[responseFields]interface{}{
					responseValue:   v.data,
					responseTimeout: v.timeout.Sub(time.Now()),
				},
			}
		}
	}

	return resp
}

func setValueOperate(d *database, key, value interface{}, timeout time.Duration) (resp dbResponse) {
	// get dirty-map index and set write-lock
	index := d.hash(key) % d.size
	_map := d.data[index]
	_map.lock.Lock()
	defer _map.lock.Unlock()

	// update dirty-map
	v := &entity{
		timeout: time.Now().Add(timeout),
		data:    value,
	}
	_map.dirty[key] = v

	// return value
	resp = dbResponse{
		responseType: withValue,
		returns: map[responseFields]interface{}{
			responseValue:   v.data,
			responseTimeout: v.timeout.Sub(time.Now()),
		},
	}

	return resp
}

func deleteKeyOperate(d *database, key interface{}) (resp dbResponse) {
	// get dirty-map index and set read-lock
	index := d.hash(key) % d.size
	_map := d.data[index]
	_map.lock.RLock()
	defer _map.lock.RUnlock()

	if _, ok := _map.dirty[key]; !ok {
		// condition: cannot find key in dirty-map
		// return value
		resp = dbResponse{
			responseType: withValue,
			returns:      nil,
		}
	} else {
		// delete key-value pair
		_map.lock.Lock()
		delete(_map.dirty, key)
		_map.lock.Unlock()

		// return value
		resp = dbResponse{
			responseType: withValue,
			returns:      nil,
		}
	}

	return resp
}

func rangeDbOperate(d *database, f func(key, value interface{})) (resp dbResponse) {
	for i := uint(0); i < d.size; i++ {
		// set write-lock
		current := d.data[i]
		current.lock.Lock()

		// execute function
		for k, v := range current.dirty {
			f(k, v.data)
		}
		current.lock.Unlock()
	}

	// return value
	resp = dbResponse{
		responseType: withValue,
		returns:      nil,
	}

	return resp
}
