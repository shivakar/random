// Package mathutils defines math and statistics utilities for developing
// probability distributions
package mathutils

import "math"

// Several functions below are simplified go implementation of the original C
// code from Cephes Math library available http://www.netlib.org/cephes/index.html
//
// Copyright information for the original C code is as follows:
//
// Cephes Math Library Release 2.8:  June, 2000
// Copyright 1984, 1987, 1989, 1992, 2000 by Stephen L. Moshier
//
// Additionally, the readme file at http://www.netlib.org/cephes/readme says:
//    Some software in this archive may be from the book _Methods and
// Programs for Mathematical Functions_ (Prentice-Hall or Simon & Schuster
// International, 1989) or from the Cephes Mathematical Library, a
// commercial product. In either event, it is copyrighted by the author.
// What you see here may be used freely but it comes with no support or
// guarantee.
//
//    The two known misprints in the book are repaired here in the
// source listings for the gamma function and the incomplete beta
// integral.
//
//    Stephen L. Moshier
//    moshier@na-net.ornl.gov

const (
	// Sqrt2Pi - Sqrt(2*Pi)
	Sqrt2Pi = 2.506628274631000502415765284811045253006986740609938316629923576 // http://oeis.org/A019727
	// SqrtI2Pi - 1.0/Sqrt(2*Pi)
	SqrtI2Pi = 0.398942280401432677939946059934381868475858631164934657665925829 // http://oeis.org/A231863
	// SqrtI2 - 1/Sqrt(2)
	SqrtI2 = 0.707106781186547524400844362104849039284835937688474036588339868 // http://oeis.org/A010503
	// Epsilon - Calculated as math.NextAfter(1, 2)-1
	Epsilon = 2.220446049250313e-16
	// MaxLog - Calculated as math.Log(math.MaxFloat64)
	MaxLog = 709.782712893384
	// Ipi - 1.0/Pi
	Ipi = 0.3183098861837906715377675267450287240689192914809128974953 // https://oeis.org/A049541
)

// PolyEval evaluates a polynomial of degree N
//
//                      2          N
//  y  =  C  + C x + C x  +...+ C x
//         0    1     2          N
//
//  Coefficients are stored in reverse order:
//
//  coef[0] = C  , ..., coef[N] = C  .
//             N                   0
//
//  This makes it easy to evaluate the polynomia using Horner's methdod
//  https://en.wikipedia.org/wiki/Horner%27s_method
//
//  The degree N is calculated as len(coef)-1
//
func PolyEval(x float64, coef []float64) float64 {
	fx := float64(0)
	for _, v := range coef {
		fx = fx*x + v
	}
	return fx
}

// NormICDF returns Polynomial approximation of inverse of Normal CDF
//
//  Based on ndtri.c from http://www.netlib.org/cephes/index.html
//  See browseable source at https://github.com/scipy/scipy/blob/master/scipy/special/cephes/ndtri.c
func NormICDF(y0 float64) float64 {

	var (
		/* approximation for 0 <= |y - 0.5| <= 3/8 */
		p0 = []float64{
			-5.99633501014107895267E1, 9.80010754185999661536E1,
			-5.66762857469070293439E1, 1.39312609387279679503E1, -1.23916583867381258016E0,
		}
		q0 = []float64{
			1.0,
			1.95448858338141759834E0, 4.67627912898881538453E0,
			8.63602421390890590575E1, -2.25462687854119370527E2,
			2.00260212380060660359E2, -8.20372256168333339912E1,
			1.59056225126211695515E1, -1.18331621121330003142E0,
		}
		/* Approximation for interval z = sqrt(-2 log y ) between 2 and 8
		   |* i.e., y between exp(-2) = .135 and exp(-32) = 1.27e-14.
		   |*/
		p1 = []float64{
			4.05544892305962419923E0, 3.15251094599893866154E1,
			5.71628192246421288162E1, 4.40805073893200834700E1,
			1.46849561928858024014E1, 2.18663306850790267539E0,
			-1.40256079171354495875E-1, -3.50424626827848203418E-2,
			-8.57456785154685413611E-4,
		}
		q1 = []float64{
			1.0,
			1.57799883256466749731E1, 4.53907635128879210584E1,
			4.13172038254672030440E1, 1.50425385692907503408E1,
			2.50464946208309415979E0, -1.42182922854787788574E-1,
			-3.80806407691578277194E-2, -9.33259480895457427372E-4,
		}
		/* Approximation for interval z = sqrt(-2 log y ) between 8 and 64
		   |* i.e., y between exp(-32) = 1.27e-14 and exp(-2048) = 3.67e-890.
		   |*/
		p2 = []float64{
			3.23774891776946035970E0, 6.91522889068984211695E0,
			3.93881025292474443415E0, 1.33303460815807542389E0,
			2.01485389549179081538E-1, 1.23716634817820021358E-2,
			3.01581553508235416007E-4, 2.65806974686737550832E-6,
			6.23974539184983293730E-9,
		}
		q2 = []float64{
			1.0,
			6.02427039364742014255E0, 3.67983563856160859403E0,
			1.37702099489081330271E0, 2.16236993594496635890E-1,
			1.34204006088543189037E-2, 3.28014464682127739104E-4,
			2.89247864745380683936E-6, 6.79019408009981274425E-9,
		}
	)

	if y0 <= float64(0) {
		return math.SmallestNonzeroFloat64
	}
	if y0 >= float64(1) {
		return math.MaxFloat64
	}
	code := 1
	y := y0

	if y > (1.0 - 0.13533528323661269189) { /* 0.135... = exp(-2) */
		y = 1.0 - y
		code = 0
	}

	if y > 0.13533528323661269189 {
		y = y - 0.5
		y2 := y * y
		x := y + y*(y2*PolyEval(y2, p0)/PolyEval(y2, q0))
		x = x * Sqrt2Pi
		return (x)
	}

	x := math.Sqrt(float64(-2.0) * math.Log(y))
	x0 := x - math.Log(x)/x

	z := float64(1.0) / x
	x1 := float64(0)
	if x < float64(8.0) { /* y > exp(-32) = 1.2664165549e-14 */
		x1 = z * PolyEval(z, p1) / PolyEval(z, q1)
	} else {
		x1 = z * PolyEval(z, p2) / PolyEval(z, q2)
	}
	x = x0 - x1
	if code != 0 {
		x = -x
	}
	return float64(x)
}
