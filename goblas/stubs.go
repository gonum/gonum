package goblas

import "github.com/gonum/blas"

func (Blas) Drotm(n int, x []float64, incX int, y []float64, incY int, p *blas.DrotmParams) {
}

func (Blas) Dgbmv(o blas.Order, tA blas.Transpose, m, n, kL, kU int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int) {
}
func (Blas) Dtrmv(o blas.Order, ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, a []float64, lda int, x []float64, incX int) {
}
func (Blas) Dtbmv(o blas.Order, ul blas.Uplo, tA blas.Transpose, d blas.Diag, n, k int, a []float64, lda int, x []float64, incX int) {
}
func (Blas) Dtpmv(o blas.Order, ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, ap []float64, x []float64, incX int) {
}
func (Blas) Dtrsv(o blas.Order, ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, a []float64, lda int, x []float64, incX int) {
}
func (Blas) Dtbsv(o blas.Order, ul blas.Uplo, tA blas.Transpose, d blas.Diag, n, k int, a []float64, lda int, x []float64, incX int) {
}
func (Blas) Dtpsv(o blas.Order, ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, ap []float64, x []float64, incX int) {
}
func (Blas) Dsymv(o blas.Order, ul blas.Uplo, n int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int) {
}
func (Blas) Dsbmv(o blas.Order, ul blas.Uplo, n, k int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int) {
}
func (Blas) Dspmv(o blas.Order, ul blas.Uplo, n int, alpha float64, ap []float64, x []float64, incX int, beta float64, y []float64, incY int) {
}
func (Blas) Dger(o blas.Order, m, n int, alpha float64, x []float64, incX int, y []float64, incY int, a []float64, lda int) {
}
func (Blas) Dsyr(o blas.Order, ul blas.Uplo, n int, alpha float64, x []float64, incX int, a []float64, lda int) {
}
func (Blas) Dspr(o blas.Order, ul blas.Uplo, n int, alpha float64, x []float64, incX int, ap []float64) {
}
func (Blas) Dsyr2(o blas.Order, ul blas.Uplo, n int, alpha float64, x []float64, incX int, y []float64, incY int, a []float64, lda int) {
}
func (Blas) Dspr2(o blas.Order, ul blas.Uplo, n int, alpha float64, x []float64, incX int, y []float64, incY int, a []float64) {
}

// Level 3 routines.
func (Blas) Dsymm(o blas.Order, s blas.Side, ul blas.Uplo, m, n int, alpha float64, a []float64, lda int, b []float64, ldb int, beta float64, c []float64, ldc int) {
}
func (Blas) Dsyrk(o blas.Order, ul blas.Uplo, t blas.Transpose, n, k int, alpha float64, a []float64, lda int, beta float64, c []float64, ldc int) {
}
func (Blas) Dsyr2k(o blas.Order, ul blas.Uplo, t blas.Transpose, n, k int, alpha float64, a []float64, lda int, b []float64, ldb int, beta float64, c []float64, ldc int) {
}
func (Blas) Dtrmm(o blas.Order, s blas.Side, ul blas.Uplo, tA blas.Transpose, d blas.Diag, m, n int, alpha float64, a []float64, lda int, b []float64, ldb int) {
}
func (Blas) Dtrsm(o blas.Order, s blas.Side, ul blas.Uplo, tA blas.Transpose, d blas.Diag, m, n int, alpha float64, a []float64, lda int, b []float64, ldb int) {
}
