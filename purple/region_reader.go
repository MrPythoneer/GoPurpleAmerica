package purple

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

type regionReader struct {
	bufio.Scanner
}

func ReadState(r io.Reader) any {
	sc := newRegionReader(r)
	return sc.scanRegion()
}

func newRegionReader(r io.Reader) *regionReader {
	return &regionReader{*bufio.NewScanner(r)}
}

func (sc *regionReader) scanPoint() Point {
	xy := strings.Split(sc.scanString(), "   ")

	xs := strings.TrimSpace(xy[0])
	ys := strings.TrimSpace(xy[1])

	x, _ := strconv.ParseFloat(xs, 64)
	y, _ := strconv.ParseFloat(ys, 64)

	return Point{x, y}
}

func (sc *regionReader) scanBBox() BBox {
	p1 := sc.scanPoint()
	p2 := sc.scanPoint()

	return NewBBox(p1, p2)
}

func (sc *regionReader) scanInt() int {
	v, _ := strconv.Atoi(sc.scanString())
	return v
}

func (sc *regionReader) scanString() string {
	sc.Scan()
	return sc.Text()
}

func (sc *regionReader) scanRegion() *State {
	reg := new(State)
	reg.Bbox = sc.scanBBox()
	reg.CountiesN = sc.scanInt()

	reg.Counties = make([]County, reg.CountiesN)
	for i := 0; i < reg.CountiesN; i++ {
		sc.Scan()

		name := sc.scanString()
		regionName := sc.scanString()
		n := sc.scanInt()

		points := make([]Point, n)
		for j := 0; j < n; j++ {
			points[j] = sc.scanPoint()
		}

		reg.Counties[i] = County{
			name,
			regionName,
			n,
			points,
		}
	}
	return reg
}
