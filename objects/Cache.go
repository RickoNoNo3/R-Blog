package objects

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/rickonono3/m2obj"
)

var Cache *m2obj.Object

func initCache() {
	rand.Seed(time.Now().Unix())
	Cache = m2obj.New(m2obj.Group{
		"AdminHash": strconv.Itoa(rand.Int()),
	})
}
