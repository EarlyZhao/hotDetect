package hotDetect

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)



func Test_FreqCount(t *testing.T) {

	Convey("Test-FreqCount", t, func() {
		fc := newFc(4,"test")
		fc.Add("ID1", 1)
		fc.Add("ID1", 1)
		fc.Add("ID1", 1)
		fc.Add("ID2", 1)
		fc.Add("ID2", 1)
		fc.Add("ID3", 1)
		fc.Add("ID3", 1)
		fc.Add("ID3", 1)
		fc.Add("ID3", 1)

		So(len(fc.fu), ShouldEqual, 3)
		fc.Add("ID4", 1)
		So(fc.size, ShouldEqual, 3)

		tops := fc.TopKey(2)
		So(tops[0].Key, ShouldEqual, "ID3")
		So(tops[1].Key, ShouldEqual, "ID1")

	})
}
