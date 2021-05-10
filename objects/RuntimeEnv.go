package objects

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/rickonono3/m2obj"
)

var RuntimeEnv *m2obj.Object

func initRuntimeEnv() {
	rand.Seed(time.Now().Unix())
	RuntimeEnv = m2obj.New(m2obj.Group{
		"AdminHash": strconv.Itoa(rand.Int()),
	})
}
