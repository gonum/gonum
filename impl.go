package lapack

var impl Lapack

func Register(i Lapack) {
	impl = i
}
