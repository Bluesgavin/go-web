package server

import (
	"fmt"
	"time"
)

type MiddleWare func(next HandlerFunc) HandlerFunc

/** a default middleWare ( record the duration ) **/
func MetricFilterBuilder(next HandlerFunc) HandlerFunc {
	return func(c *Context) {
		// 执行前的时间
		startTime := time.Now().UnixNano()
		next(c)
		// 执行后的时间
		endTime := time.Now().UnixNano()
		fmt.Printf("run time: %d \n", endTime-startTime)
	}
}
