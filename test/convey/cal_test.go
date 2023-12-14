package convey

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestCal(t *testing.T) {
	convey.Convey("normal plus", t, func() {
		symbol, p1, p2 := "+", 1, 2
		sum, _ := Cal(symbol, p1, p2)
		convey.So(sum, convey.ShouldEqual, 3)
	})
	convey.Convey("normal sub", t, func() {
		symbol, p1, p2 := "-", 1, 2
		sub, _ := Cal(symbol, p1, p2)
		convey.So(sub, convey.ShouldEqual, -1)
	})
	convey.Convey("wrong multip", t, func() {
		symbol, p1, p2 := "*", 0, 1
		m, _ := Cal(symbol, p1, p2)
		convey.So(m, convey.ShouldEqual, 0)
	})
	convey.Convey("divide", t, func() {
		symbol, p1, p2 := "/", 1, 0
		_, err := Cal(symbol, p1, p2)
		convey.So(err.Error(), convey.ShouldEqual, "divide by zero")
	})
}
