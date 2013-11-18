package goblas

import "math"

import "github.com/gonum/blas"

// Compute a modified Givens transformation
func (Blas) Drotmg(d1, d2, x1, y1 float64) (p *blas.DrotmParams, rd1, rd2, rx1 float64) {
	var p1, p2, q1, q2, u float64

	gam := 4096.0
	gamsq := 16777216.0
	rgamsq := 5.9604645e-8

	rd1 = d1
	rd2 = d2
	rx1 = x1

	if d1 < 0 {
		p.Flag = -1
	} else {
		p2 = rd2 * y1
		if p2 == 0 {
			p.Flag = -2
			return
		}
		p1 = rd1 * x1
		q2 = p2 * y1
		q1 = p1 * x1
		if math.Abs(q1) > math.Abs(q2) {
			p.H[1] = -y1 / x1
			p.H[2] = p2 / p1
			u = 1 - p.H[2]*p.H[1]
			if u > 0 {
				p.Flag = 0
				rd1 /= u
				rd2 /= u
				rx1 /= u
			}
		} else {
			if q2 < 0 {
				p.Flag = -1
				rd1 = 0
				rd2 = 0
				rx1 = 0
			} else {
				p.Flag = 1
				p.H[0] = p1 / p2
				p.H[3] = x1 / y1
				u = 1 + p.H[0] + p.H[3]
				rd1, rd2 = rd2/u, rd1/u
				rx1 = y1 / u
			}
		}
		if rd1 != 0 {
			for rd1 <= rgamsq || rd1 >= gamsq {
				if p.Flag == 0 {
					p.H[0] = 1
					p.H[3] = 1
					p.Flag = -1
				} else {
					p.H[1] = -1
					p.H[2] = 1
					p.Flag = -1
				}
				if rd1 <= rgamsq {
					rd1 *= gam * gam
					rx1 /= gam
					p.H[0] /= gam
					p.H[2] /= gam
				} else {
					rd1 /= gam * gam
					rx1 *= gam
					p.H[0] *= gam
					p.H[2] *= gam
				}
			}
		}
		if rd2 != 0 {
			for math.Abs(rd2) <= rgamsq || math.Abs(rd2) >= gamsq {
				if p.Flag == 0 {
					p.H[0] = 1
					p.H[3] = 1
					p.Flag = -1
				} else {
					p.H[1] = -1
					p.H[2] = 1
					p.Flag = -1
				}
				if math.Abs(rd2) <= rgamsq {
					rd2 *= gam * gam
					p.H[1] /= gam
					p.H[3] /= gam
				} else {
					rd2 /= gam * gam
					p.H[1] *= gam
					p.H[3] *= gam
				}
			}
		}
	}
	return
}
