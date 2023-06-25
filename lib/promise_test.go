package lib

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSimple(t *testing.T) {
	p, resolve, reject := Promise()

	go func() {
		time.Sleep(100 * time.Millisecond)

		if err := error(nil); err != nil {
			reject(err)
			return
		}
		resolve("Done!")
	}()

	v, err := Await(p)
	if err != nil {
		fmt.Printf("error: %#v\n", err.Error())
		return
	}

	require.Equal(t, v.String(), "Done!")
}
