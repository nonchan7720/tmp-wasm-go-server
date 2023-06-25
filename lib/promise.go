package lib

import "syscall/js"

type PromiseCallback func(value any)

func Promise() (js.Value, PromiseCallback, PromiseCallback) {
	var (
		cb      js.Func
		resolve PromiseCallback
		reject  PromiseCallback
	)
	cb = js.FuncOf(func(_ js.Value, args []js.Value) any {
		cb.Release()
		resolve = func(value any) {
			args[0].Invoke(value)
		}

		reject = func(value any) {
			args[1].Invoke(value)
		}

		return js.Undefined()
	})
	p := js.Global().Get("Promise").New(cb)
	return p, resolve, reject
}

func Await(p js.Value) (js.Value, error) {
	if !hasThenMethod(p) {
		return p, nil
	}

	resCh, onFulfilled := argChanFunc()
	defer onFulfilled.Release()
	errCh, onRejected := argChanFunc()
	defer onRejected.Release()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				errCh <- js.ValueOf(r)
			}
		}()
		p.Call("then", onFulfilled, onRejected)
	}()

	for {
		select {
		case res := <-resCh:
			if !hasThenMethod(res) {
				return res, nil
			}
			go res.Call("then", onFulfilled, onRejected)
		case err := <-errCh:
			return js.Undefined(), Reason(err)
		}
	}
}

type Reason js.Value

var _ error = Reason{}

func (r Reason) Error() string {
	v := js.Value(r)

	if v.Type() == js.TypeObject {
		if message := v.Get("message"); message.Type() == js.TypeString {
			return message.String()
		}
	}

	return js.Global().Call("String", v).String()
}

func hasThenMethod(v js.Value) bool {
	t := v.Type()
	return (t == js.TypeObject || t == js.TypeFunction) && v.Get("then").Type() == js.TypeFunction
}

func argChanFunc() (chan js.Value, js.Func) {
	ch := make(chan js.Value)
	return ch, js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		ch <- args[0]
		return nil
	})
}
