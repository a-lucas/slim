package polyfit

import (
	"fmt"
	"strings"

	"gonum.org/v1/gonum/mat"
)

func PolyStr(poly []float64) string {

	elts := []string{}

	for i, coef := range poly {
		if coef == 0 {
			continue
		}

		cc := fmt.Sprintf("%.4f", coef)

		if i > 0 {
			cc += "x"
		}
		if i > 1 {
			cc += fmt.Sprintf("^%d", i)
		}

		elts = append(elts, cc)
	}

	return strings.Replace(strings.Join(elts, " + "), "+ -", "- ", -1)

}

// PolyFit models a polynomial y from sample points xs and ys, to minimizes the squared residuals.
// It returns coefficients of the polynomial y:
//
//    y = β₁ + β₂x + β₃x² + ...
//
// It use linear regression, which assumes y is in form of:
//        m
//    y = ∑ βⱼ Φⱼ(x)
//        j=1
//
// In our case:
//    Φⱼ(x) = x^(j-1)
//
// Then
//    (Xᵀ X) βⱼ = Xᵀ Y
//    Xᵢⱼ = [ Φⱼ(xᵢ) ]
//
// See https://en.wikipedia.org/wiki/Least_squares#Linear_least_squares
func PolyFit(xs, ys []float64, degree int) []float64 {

	// Number of sample points
	n := len(xs)

	// Number of βⱼ is degree+1
	m := degree + 1

	// build matrix: Xᵢⱼ = [ Φⱼ(xᵢ) ]
	d := make([]float64, n*m)
	for i := 0; i < n; i++ {
		x := xs[i]
		v := float64(1)
		for j := 0; j < m; j++ {
			d[i*m+j] = v
			v *= x
		}
	}

	mtr := mat.NewDense(n, m, d)

	var right mat.Dense
	var coef mat.Dense
	var beta mat.Dense

	// coef * beta = right
	coef.Mul(mtr.T(), mtr)
	right.Mul(mtr.T(), mat.NewDense(n, 1, ys))
	beta.Solve(&coef, &right)

	rst := make([]float64, m)
	for i := 0; i < m; i++ {
		rst[i] = beta.At(i, 0)
	}

	return rst
}

// eval evaluates polynomial at x
func eval(poly []float64, x float64) float64 {
	return poly[0] + poly[1]*x + poly[2]*x*x + poly[3]*x*x*x
}

// maxminResiduals finds max and min offset along a curve.
func maxminResiduals(poly, xs, ys []float64) (float64, float64) {

	max, min := float64(0), float64(0)

	for i, x := range xs {
		v := eval(poly, x)
		diff := ys[i] - v
		if diff > max {
			max = diff
		}
		if diff < min {
			min = diff
		}
	}

	return max, min
}

// FindPoly finds a polynomial curve that has as many points as possible so that
// their distant to the curve smaller than margin.
//
// It returns the coeffecients of the curve and how many points is covered.
func FindPoly(xs, ys []float64, degree int, margin float64) ([]float64, int) {

	l, r := 0, len(xs)+1

	for {
		for l < r-1 {
			mid := (l + r) / 2
			xx, yy := xs[:mid], ys[:mid]

			poly := PolyFit(xx, yy, degree)
			max, min := maxminResiduals(poly, xx, yy)
			if max-min <= margin {
				l = mid
			} else {
				r = mid
			}
		}

		xs, ys = xs[:l], ys[:l]
		poly := PolyFit(xs, ys, degree)
		max, min := maxminResiduals(poly, xs, ys)

		// max-min are not guaranteed to be incremental.
		// Thus if max-min exceed margin, reset r to l and re-run binary search.
		if max-min > margin {
			l, r = 0, l
			continue
		} else {
			// Makes every point be above the curve
			poly[0] += min
			return poly, l
		}
	}
}

type reg struct {
	Poly  []float64
	start int32
}

type XArray32 struct {
	polyDegree       byte
	eltWidth         byte
	eltPerWord       byte
	inWordIndexWidth uint8
	eltMask          uint64

	// Regions    []reg
	Poly  [][]float64
	start []int32

	Datas []uint64
}

func (x *XArray32) Get(i int32) int32 {
	var j int
	l := len(x.start)
	for j = 0; j < l && i >= x.start[j]; j++ {
	}

	r := x.Poly[j-1]

	v := eval(r, float64(i))

	// d := x.Datas[i/int32(x.eltPerWord)]
	d := x.Datas[i>>x.inWordIndexWidth]
	d = d >> (uint(i&int32((1<<x.inWordIndexWidth)-1)) * uint(x.eltWidth))
	return int32(v) + int32(d&uint64(x.eltMask))

}

func Resample32To4(nums []int32, eltWidth uint) *XArray32 {
	lg := map[uint]uint8{
		1:  6,
		2:  5,
		4:  4,
		8:  3,
		16: 2,
	}

	n := len(nums)
	xs := make([]float64, n)
	ys := make([]float64, n)

	for i, v := range nums {
		xs[i] = float64(i)
		ys[i] = float64(v)
	}

	// for 4bit int
	marginInt := (1 << eltWidth) - 1
	margin := float64(marginInt)

	eltPerWord := 64 / eltWidth
	nWords := (n + int(eltPerWord) - 1) / int(eltPerWord)

	rst := &XArray32{
		polyDegree:       3,
		eltWidth:         byte(eltWidth),
		eltMask:          (1 << eltWidth) - 1,
		inWordIndexWidth: uint8(lg[eltWidth]),
		eltPerWord:       byte(eltPerWord),
		Datas:            make([]uint64, nWords),
	}

	debugdata := []int32{}

	for start := 0; start < n; {
		poly, nn := FindPoly(xs[start:], ys[start:], int(rst.polyDegree), margin)

		rst.start = append(rst.start, int32(start))
		rst.Poly = append(rst.Poly, poly)

		for i := 0; i < nn; i++ {
			j := start + i

			v := eval(poly, xs[j])

			d := nums[j] - int32(v)
			if d > int32(marginInt) {
				panic(fmt.Sprintf("d must smaller than %d", marginInt))
			}

			rst.Datas[j/int(rst.eltPerWord)] |= uint64(d) << uint(int(eltWidth)*(j%int(rst.eltPerWord)))

			debugdata = append(debugdata, d)
		}

		start += nn
	}

	return rst
}
