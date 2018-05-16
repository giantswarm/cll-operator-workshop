package logger

import (
	"fmt"
	"time"

	"github.com/giantswarm/micrologger"
)

var Default micrologger.Logger

func init() {
	var err error

	c := micrologger.Config{
		TimestampFormatter: func() interface{} {
			return time.Now().UTC().Format("2006/01/02 15:04:05.000")
		},
	}

	Default, err = micrologger.New(c)
	if err != nil {
		panic(fmt.Sprintf("%#v", err))
	}
}
