package sliceops

import (
	"testing"
)

const (
	SmallBenchmark  = 10
	MediumBenchmark = 1000
	LargeBenchmark  = 100000
)

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

func TestHasEqLen(t *testing.T) {
	s1 := []float64{1, 2, 3, 4}
	s2 := []float64{1, 2, 3, 4}
	s3 := []float64{1, 2, 3}
	if !HasEqLen(s1, s2) {
		t.Errorf("Equal lengths returned as unequal")
	}
	if HasEqLen(s1, s3) {
		t.Errorf("Unequal lengths returned as equal")
	}
	if !HasEqLen(s1) {
		t.Errorf("Single slice returned as unequal")
	}
	if !HasEqLen() {
		t.Errorf("No slices returned as unequal")
	}
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

func TestCumSum(t *testing.T) {
	s := []float64{3, 4, 1, 7, 5}
	val := CumSum(nil, s)
	truth := []float64{3, 7, 8, 15, 20}
	if !Eq(val, truth, 1E-15) {
		t.Errorf("Wrong cumsum returned with nil receiver. Expected %v, returned %v", truth, val)
	}
	val = CumSum(val, s)
	if !Eq(val, truth, 1E-15) {
		t.Errorf("Wrong cumsum returned with non-nil receiver. Expected %v, returned %v", truth, val)
	}
}

func TestCumProd(t *testing.T) {
	s := []float64{3, 4, 1, 7, 5}
	val := CumProd(nil, s)
	truth := []float64{3, 12, 12, 84, 420}
	if !Eq(val, truth, 1E-15) {
		t.Errorf("Wrong cumprod returned with nil receiver. Expected %v, returned %v", truth, val)
	}
	val = CumProd(val, s)
	if !Eq(val, truth, 1E-15) {
		t.Errorf("Wrong cumprod returned with non-nil receiver. Expected %v, returned %v", truth, val)
	}
}
