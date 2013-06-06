package sliceops

import (
	"fmt"
	"math"
	"testing"
)

const (
	EQTOLERANCE     = 1E-14
	SmallBenchmark  = 10
	MediumBenchmark = 1000
	LargeBenchmark  = 100000
)

func AreSlicesEqual(t *testing.T, truth, comp []float64, str string) {
	if !Eq(comp, truth, EQTOLERANCE) {
		t.Errorf(str+"Expected %v, returned %v", truth, comp)
	}
}

func TestAdd(t *testing.T) {
	a := []float64{1, 2, 3}
	b := []float64{4, 5, 6}
	c := []float64{7, 8, 9}
	truth := []float64{12, 15, 18}
	n := make([]float64, len(a))
	Add(n, a, b, c)
	AreSlicesEqual(t, truth, n, "Wrong addition of slices new receiver")
	Add(a, b, c)
	AreSlicesEqual(t, truth, n, "Wrong addition of slices for no new receiver")
}

func TestAddconst(t *testing.T) {
	s := []float64{3, 4, 1, 7, 5}
	c := 6.0
	truth := []float64{9, 10, 7, 13, 11}
	AddConst(s, c)
	AreSlicesEqual(t, truth, s, "Wrong addition of constant")
}

func TestApplyFunc(t *testing.T) {
	s := []float64{3, 4, 1, 7, 5}
	f := math.Sin
	truth := make([]float64, len(s))
	for i, val := range s {
		truth[i] = math.Sin(val)
	}
	ApplyFunc(s, f)
	AreSlicesEqual(t, truth, s, "Wrong application of function")
}

func TestCumprod(t *testing.T) {
	s := []float64{3, 4, 1, 7, 5}
	receiver := make([]float64, len(s))
	Cumprod(receiver, s)
	truth := []float64{3, 12, 12, 84, 420}
	AreSlicesEqual(t, truth, receiver, "Wrong cumprod returned with new receiver")
	Cumprod(receiver, s)
	AreSlicesEqual(t, truth, receiver, "Wrong cumprod returned with reused receiver")
}

func TestCumSum(t *testing.T) {
	s := []float64{3, 4, 1, 7, 5}
	receiver := make([]float64, len(s))
	Cumsum(receiver, s)
	truth := []float64{3, 7, 8, 15, 20}
	AreSlicesEqual(t, truth, receiver, "Wrong cumsum returned with new receiver")
	Cumsum(receiver, s)
	AreSlicesEqual(t, truth, receiver, "Wrong cumsum returned with reused receiver")
}

func TestEq(t *testing.T) {
	s1 := []float64{1, 2, 3, 4}
	s2 := []float64{1, 2, 3, 4 + 1E-14}
	if !Eq(s1, s2, 1E-13) {
		t.Errorf("Equal slices returned as unequal")
	}
	if Eq(s1, s2, 1E-15) {
		t.Errorf("Unequal slices returned as equal")
	}
}

func TestEqLen(t *testing.T) {
	s1 := []float64{1, 2, 3, 4}
	s2 := []float64{1, 2, 3, 4}
	s3 := []float64{1, 2, 3}
	if !Eqlen(s1, s2) {
		t.Errorf("Equal lengths returned as unequal")
	}
	if Eqlen(s1, s3) {
		t.Errorf("Unequal lengths returned as equal")
	}
	if !Eqlen(s1) {
		t.Errorf("Single slice returned as unequal")
	}
	if !Eqlen() {
		t.Errorf("No slices returned as unequal")
	}
}

func TestFind(t *testing.T) {
	s := []float64{3, 4, 1, 7, 5}
	f := func(v float64) bool { return v > 3.5 }
	trueInds := []int{1, 3, 4}
	inds := Find(s, f)
	if len(inds) != len(trueInds) {
		t.Errorf("Wrong number of elements returned")
		return
	}
	for i, val := range trueInds {
		if inds[i] != val {
			t.Errorf("Index mismatch")
			fmt.Println(trueInds)
			fmt.Println(inds)
			return
		}
	}
}

func TestFindFirst(t *testing.T) {
	s := []float64{3, 4, 1, 7, 5}
	f := func(v float64) bool { return v > 3.5 }
	trueInds := []int{1, 3}
	k := 2
	inds, err := FindFirst(s, f, k)
	if err != nil {
		t.Errorf("Incorrectly did not find enough elements")
	}
	if len(inds) != len(trueInds) {
		t.Errorf("Wrong number of elements returned")
		return
	}
	for i, val := range trueInds {
		if inds[i] != val {
			t.Errorf("Index mismatch")
			fmt.Println(trueInds)
			fmt.Println(inds)
			return
		}
	}
	f = func(v float64) bool { return v > 6.0 }
	inds, err = FindFirst(s, f, k)
	if err == nil {
		t.Errorf("Incorrectly found enough elements")
	}
}

func TestLinspace(t *testing.T) {
	receiver := make([]float64, 5)
	truth := []float64{1, 2, 3, 4, 5}
	Linspace(receiver, 1, 5)
	AreSlicesEqual(t, truth, receiver, "Improper linspace")
	receiver = make([]float64, 6)
	truth = []float64{0, 0.2, 0.4, 0.6, 0.8, 1.0}
	Linspace(receiver, 0, 1)
	AreSlicesEqual(t, truth, receiver, "Improper linspace")
}

func TestLogspace(t *testing.T) {
	receiver := make([]float64, 6)
	truth := []float64{0.001, 0.01, 0.1, 1, 10, 100}
	Logspace(receiver, 0.001, 100)
	tst := make([]float64, 6)
	for i := range truth {
		tst[i] = receiver[i] / truth[i]
	}
	comp := make([]float64, 6)
	for i := range comp {
		comp[i] = 1
	}
	AreSlicesEqual(t, comp, tst, "Improper logspace")
}

func TestLogsumexp(t *testing.T) {
	s := []float64{1, 2, 3, 4, 5}
	val := Logsumexp(s)
	// http://www.wolframalpha.com/input/?i=log%28exp%281%29+%2B+exp%282%29+%2B+exp%283%29+%2B+exp%284%29+%2B+exp%285%29%29
	truth := 5.4519143959375933331957225109748087179338972737576824
	if math.Abs(val-truth) > EQTOLERANCE {
		t.Errorf("Wrong logsumexp for many values")
	}
	s = []float64{1, 2}
	// http://www.wolframalpha.com/input/?i=log%28exp%281%29+%2B+exp%282%29%29
	truth = 2.3132616875182228340489954949678556419152800856703483
	val = Logsumexp(s)
	if math.Abs(val-truth) > EQTOLERANCE {
		t.Errorf("Wrong logsumexp for two values. %v expected, %v found", truth, val)
	}
	// This case would normally underflow
	s = []float64{-1001, -1002, -1003, -1004, -1005}
	// http://www.wolframalpha.com/input/?i=log%28exp%28-1001%29%2Bexp%28-1002%29%2Bexp%28-1003%29%2Bexp%28-1004%29%2Bexp%28-1005%29%29
	truth = -1000.54808560406240666680427748902519128206610272624
	val = Logsumexp(s)
	if math.Abs(val-truth) > EQTOLERANCE {
		t.Errorf("Doesn't match for underflow case. %v expected, %v found", truth, val)
	}
}

func TestMax(t *testing.T) {
	s := []float64{}
	val, ind := Max(s)
	if val != 0 {
		t.Errorf("Val not returned as default when slice length is zero")
	}
	if ind != -1 {
		t.Errorf("Ind not returned as -1 for empty slice")
	}
	s = []float64{3, 4, 1, 7, 5}
	val, ind = Max(s)
	if val != 7 {
		t.Errorf("Wrong value returned")
	}
	if ind != 3 {
		t.Errorf("Wrong index returned")
	}
}

func TestMin(t *testing.T) {
	s := []float64{}
	val, ind := Min(s)
	if val != 0 {
		t.Errorf("Val not returned as default when slice length is zero")
	}
	if ind != -1 {
		t.Errorf("Ind not returned as -1 for empty slice")
	}
	s = []float64{3, 4, 1, 7, 5}
	val, ind = Min(s)
	if val != 1 {
		t.Errorf("Wrong value returned")
	}
	if ind != 2 {
		t.Errorf("Wrong index returned")
	}
}

func TestNorm(t *testing.T) {
	s := []float64{-1, -3.4, 5, 6}
	val := Norm(s, math.Inf(1))
	truth := 6.0
	if math.Abs(val-truth) > EQTOLERANCE {
		t.Errorf("Doesn't match for inf norm. %v expected, %v found", truth, val)
	}
	// http://www.wolframalpha.com/input/?i=%28%28-1%29%5E2+%2B++%28-3.4%29%5E2+%2B+5%5E2%2B++6%5E2%29%5E%281%2F2%29
	val = Norm(s, 2)
	truth = 8.5767126569566267590651614132751986658027271236078592
	if math.Abs(val-truth) > EQTOLERANCE {
		t.Errorf("Doesn't match for inf norm. %v expected, %v found", truth, val)
	}
	// http://www.wolframalpha.com/input/?i=%28%28%7C-1%7C%29%5E3+%2B++%28%7C-3.4%7C%29%5E3+%2B+%7C5%7C%5E3%2B++%7C6%7C%5E3%29%5E%281%2F3%29
	val = Norm(s, 3)
	truth = 7.2514321388020228478109121239004816430071237369356233
	if math.Abs(val-truth) > EQTOLERANCE {
		t.Errorf("Doesn't match for inf norm. %v expected, %v found", truth, val)
	}

	//http://www.wolframalpha.com/input/?i=%7C-1%7C+%2B+%7C-3.4%7C+%2B+%7C5%7C%2B++%7C6%7C
	val = Norm(s, 1)
	truth = 15.4
	if math.Abs(val-truth) > EQTOLERANCE {
		t.Errorf("Doesn't match for inf norm. %v expected, %v found", truth, val)
	}
}

func TestProd(t *testing.T) {
	s := []float64{}
	val := Prod(s)
	if val != 1 {
		t.Errorf("Val not returned as default when slice length is zero")
	}
	s = []float64{3, 4, 1, 7, 5}
	val = Prod(s)
	if val != 420 {
		t.Errorf("Wrong prod returned. Expected %v returned %v", 420, val)
	}
}

func TestScale(t *testing.T) {
	s := []float64{3, 4, 1, 7, 5}
	c := 5.0
	truth := []float64{15, 20, 5, 35, 25}
	Scale(s, c)
	AreSlicesEqual(t, truth, s, "Bad scaling")
}

func TestSub(t *testing.T) {
	s := []float64{3, 4, 1, 7, 5}
	v := []float64{1, 2, 3, 4, 5}
	truth := []float64{2, 2, -2, 3, 0}
	Sub(s, v)
	AreSlicesEqual(t, truth, s, "Bad subtract")
}

func TestSubDst(t *testing.T) {
	s := []float64{3, 4, 1, 7, 5}
	v := []float64{1, 2, 3, 4, 5}
	truth := []float64{2, 2, -2, 3, 0}
	dst := make([]float64, len(s))
	SubDst(dst, s, v)
	AreSlicesEqual(t, truth, dst, "Bad subtract")
}

func TestSum(t *testing.T) {
	s := []float64{}
	val := Sum(s)
	if val != 0 {
		t.Errorf("Val not returned as default when slice length is zero")
	}
	s = []float64{3, 4, 1, 7, 5}
	val = Sum(s)
	if val != 20 {
		t.Errorf("Wrong sum returned")
	}
}
