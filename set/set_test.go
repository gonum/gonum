package set_test

import (
	"github.com/gonum/graph/set"
	"testing"
)

func TestAdd(t *testing.T) {
	s := set.NewSet()
	if s == nil {
		t.Fatal("Set cannot be created successfully")
	}

	if s.Cardinality() != 0 {
		t.Error("Set somehow contains new elements upon creation")
	}

	s.Add(1)
	s.Add(3)
	s.Add(5)

	if s.Cardinality() != 3 {
		t.Error("Incorrect number of set elements after adding")
	}

	if !s.Contains(1) || !s.Contains(3) || !s.Contains(5) {
		t.Error("Set doesn't contain element that was added")
	}

	s.Add(1)

	if s.Cardinality() > 3 {
		t.Error("Set double-adds element (element not unique)")
	} else if s.Cardinality() < 3 {
		t.Error("Set double-add lowered cardinality")
	}

	if !s.Contains(1) {
		t.Error("Set doesn't contain double-added element")
	}

	if !s.Contains(3) || !s.Contains(5) {
		t.Error("Set removes element on double-add")
	}

}

func TestRemove(t *testing.T) {
	s := set.NewSet()

	s.Add(1)
	s.Add(3)
	s.Add(5)

	s.Remove(1)

	if s.Cardinality() != 2 {
		t.Error("Incorrect number of set elements after removing an element")
	}

	if s.Contains(1) {
		t.Error("Element present after removal")
	}

	if !s.Contains(3) || !s.Contains(5) {
		t.Error("Set remove removed wrong element")
	}

	s.Remove(1)

	if s.Cardinality() != 2 || s.Contains(1) {
		t.Error("Double set remove does something strange")
	}

	s.Add(1)

	if s.Cardinality() != 3 || !s.Contains(1) {
		t.Error("Cannot add element after removal")
	}

}

func TestElements(t *testing.T) {
	s := set.NewSet()
	el := s.Elements()
	if el == nil {
		t.Errorf("Elements of empty set incorrectly returns nil and not zero-length slice")
	}

	if len(el) != 0 {
		t.Errorf("Elements of empty set has len greater than 0")
	}

	s.Add(1)
	s.Add(2)

	el = s.Elements()
	if len(el) != 2 {
		t.Fatalf("Elements not of same cardinality as set that spawned it")
	}

	if e, ok := el[0].(int); !ok {
		t.Errorf("Element in elements not of right type, %v", e)
	} else if e != 1 && e != 2 {
		t.Errorf("Element in elements has incorrect value %d", e)
	}

	if e, ok := el[1].(int); !ok {
		t.Errorf("Element in elements not of right type, %v", e)
	} else if e != 1 && e != 2 {
		t.Errorf("Element in elements has incorrect value %d", e)
	}

	el[0] = 19
	el[1] = 52

	if !s.Contains(1) || !s.Contains(2) || s.Cardinality() != 2 {
		t.Error("Mutating elements slice mutates set")
	}
}

func TestClear(t *testing.T) {
	s := set.NewSet()

	s.Add(8)
	s.Add(9)
	s.Add(10)

	s.Clear()

	if s.Cardinality() != 0 {
		t.Error("Clear did not properly reset set to size 0")
	}
}

func TestSelfEqual(t *testing.T) {
	s := set.NewSet()

	if !set.Equal(s, s) {
		t.Error("Set is not equal to itself")
	}

	s.Add(1)

	if !set.Equal(s, s) {
		t.Error("Set ceases self equality after adding element")
	}
}

func TestEqual(t *testing.T) {
	s1 := set.NewSet()
	s2 := set.NewSet()

	if !set.Equal(s1, s2) {
		t.Error("Two different empty sets not equal")
	}

	s1.Add(1)
	if set.Equal(s1, s2) {
		t.Error("Two different sets with different elements not equal")
	}

	s2.Add(1)
	if !set.Equal(s1, s2) {
		t.Error("Two sets with same element not equal")
	}
}

func TestCopy(t *testing.T) {
	s1 := set.NewSet()
	s2 := set.NewSet()

	s1.Add(1)
	s1.Add(2)
	s1.Add(3)

	s2.Copy(s1)

	if !set.Equal(s1, s2) {
		t.Fatalf("Two sets not equal after copy")
	}

	s2.Remove(1)

	if set.Equal(s1, s2) {
		t.Errorf("Mutating one set mutated another after copy")
	}
}

func TestSelfCopy(t *testing.T) {
	s1 := set.NewSet()

	s1.Add(1)
	s1.Add(2)

	s1.Copy(s1)

	if s1.Cardinality() != 2 {
		t.Error("Something strange happened when copying into self")
	}
}

func TestUnionSame(t *testing.T) {
	s1 := set.NewSet()
	s2 := set.NewSet()
	s3 := set.NewSet()

	s1.Add(1)
	s1.Add(2)

	s2.Add(1)
	s2.Add(2)

	s3.Union(s1, s2)

	if s3.Cardinality() != 2 {
		t.Error("Union of same sets yields set with wrong cardinality")
	}

	if !s3.Contains(1) || !s3.Contains(2) {
		t.Error("Union of same sets yields wrong elements")
	}
}

func TestUnionDiff(t *testing.T) {
	s1 := set.NewSet()
	s2 := set.NewSet()
	s3 := set.NewSet()

	s1.Add(1)
	s1.Add(2)

	s2.Add(3)

	s3.Union(s1, s2)

	if s3.Cardinality() != 3 {
		t.Error("Union of different sets yields set with wrong cardinality")
	}

	if !s3.Contains(1) || !s3.Contains(2) || !s3.Contains(3) {
		t.Error("Union of different sets yields set with wrong elements")
	}

	if s1.Contains(3) || !s1.Contains(2) || !s1.Contains(1) || s1.Cardinality() != 2 {
		t.Error("Union of sets mutates non-destination set (argument 1)")
	}

	if !s2.Contains(3) || s2.Contains(1) || s2.Contains(2) || s2.Cardinality() != 1 {
		t.Error("Union of sets mutates non-destination set (argument 2)")
	}
}

func TestUnionOverlapping(t *testing.T) {
	s1 := set.NewSet()
	s2 := set.NewSet()
	s3 := set.NewSet()

	s1.Add(1)
	s1.Add(2)

	s2.Add(2)
	s2.Add(3)

	s3.Union(s1, s2)

	if s3.Cardinality() != 3 {
		t.Error("Union of overlapping sets yields set with wrong cardinality")
	}

	if !s3.Contains(1) || !s3.Contains(2) || !s3.Contains(3) {
		t.Error("Union of overlapping sets yields set with wrong elements")
	}

	if s1.Contains(3) || !s1.Contains(2) || !s1.Contains(1) || s1.Cardinality() != 2 {
		t.Error("Union of sets mutates non-destination set (argument 1)")
	}

	if !s2.Contains(3) || s2.Contains(1) || !s2.Contains(2) || s2.Cardinality() != 2 {
		t.Error("Union of sets mutates non-destination set (argument 2)")
	}
}

func TestIntersectionSame(t *testing.T) {
	s1 := set.NewSet()
	s2 := set.NewSet()
	s3 := set.NewSet()

	s1.Add(2)
	s1.Add(3)

	s2.Add(2)
	s2.Add(3)

	s3.Intersection(s1, s2)

	if card := s3.Cardinality(); card != 2 {
		t.Errorf("Intersection of identical sets yields set of wrong cardinality %d", card)
	}

	if !s3.Contains(2) || !s3.Contains(3) {
		t.Error("Intersection of identical sets yields set of wrong elements")
	}
}

func TestIntersectionDiff(t *testing.T) {
	s1 := set.NewSet()
	s2 := set.NewSet()
	s3 := set.NewSet()

	s1.Add(2)
	s1.Add(3)

	s2.Add(1)
	s2.Add(4)

	s3.Intersection(s1, s2)

	if card := s3.Cardinality(); card != 0 {
		t.Errorf("Intersection of different yields non-empty set %d", card)
	}

	if !s1.Contains(2) || !s1.Contains(3) || s1.Contains(1) || s1.Contains(4) || s1.Cardinality() != 2 {
		t.Error("Intersection of sets mutates non-destination set (argument 1)")
	}

	if s2.Contains(2) || s2.Contains(3) || !s2.Contains(1) || !s2.Contains(4) || s2.Cardinality() != 2 {
		t.Error("Intersection of sets mutates non-destination set (argument 1)")
	}
}

func TestIntersectionOverlapping(t *testing.T) {
	s1 := set.NewSet()
	s2 := set.NewSet()
	s3 := set.NewSet()

	s1.Add(2)
	s1.Add(3)

	s2.Add(3)
	s2.Add(4)

	s3.Intersection(s1, s2)

	if card := s3.Cardinality(); card != 1 {
		t.Errorf("Intersection of overlapping sets yields set of incorrect cardinality %d", card)
	}

	if !s3.Contains(3) {
		t.Errorf("Intersection of overlapping sets yields set with wrong element")
	}

	if !s1.Contains(2) || !s1.Contains(3) || s1.Contains(4) || s1.Cardinality() != 2 {
		t.Error("Intersection of sets mutates non-destination set (argument 1)")
	}

	if s2.Contains(2) || !s2.Contains(3) || !s2.Contains(4) || s2.Cardinality() != 2 {
		t.Error("Intersection of sets mutates non-destination set (argument 1)")
	}
}

func TestDiffSame(t *testing.T) {
	s1 := set.NewSet()
	s2 := set.NewSet()
	s3 := set.NewSet()

	s1.Add(1)
	s1.Add(2)

	s2.Add(1)
	s2.Add(2)

	s3.Diff(s1, s2)

	if card := s3.Cardinality(); card != 0 {
		t.Errorf("Difference of identical sets yields set with wrong cardinality %d", card)
	}

	if !s1.Contains(1) || !s1.Contains(2) || s1.Cardinality() != 2 {
		t.Error("Difference of identical sets mutates non-destination set (argument 1)")
	}

	if !s2.Contains(1) || !s2.Contains(2) || s2.Cardinality() != 2 {
		t.Error("Difference of identical sets mutates non-destination set (argument 1)")
	}
}

func TestDiffDiff(t *testing.T) {
	s1 := set.NewSet()
	s2 := set.NewSet()
	s3 := set.NewSet()

	s1.Add(1)
	s1.Add(2)

	s2.Add(3)
	s2.Add(4)

	s3.Diff(s1, s2)

	if card := s3.Cardinality(); card != 2 {
		t.Errorf("Difference of different sets yields set with wrong cardinality %d", card)
	}

	if !s3.Contains(1) || !s3.Contains(2) || s3.Contains(3) || s3.Contains(4) {
		t.Error("Difference of different sets yields set with wrong elements")
	}

	if !s1.Contains(1) || !s1.Contains(2) || s1.Contains(3) || s1.Contains(4) || s1.Cardinality() != 2 {
		t.Error("Difference of different sets mutates non-destination set (argument 1)")
	}

	if s2.Contains(1) || s2.Contains(2) || !s2.Contains(3) || !s2.Contains(4) || s2.Cardinality() != 2 {
		t.Error("Difference of different sets mutates non-destination set (argument 1)")
	}
}

func TestDiffOverlapping(t *testing.T) {
	s1 := set.NewSet()
	s2 := set.NewSet()
	s3 := set.NewSet()

	s1.Add(1)
	s1.Add(2)

	s2.Add(2)
	s2.Add(3)

	s3.Diff(s1, s2)

	if card := s3.Cardinality(); card != 1 {
		t.Errorf("Difference of overlapping sets yields set with wrong cardinality %d", card)
	}

	if !s3.Contains(1) || s3.Contains(2) || s3.Contains(3) {
		t.Error("Difference of overlapping sets yields set with wrong elements")
	}

	if !s1.Contains(1) || !s1.Contains(2) || s1.Contains(3) || s1.Cardinality() != 2 {
		t.Error("Difference of overlapping sets mutates non-destination set (argument 1)")
	}

	if s2.Contains(1) || !s2.Contains(2) || !s2.Contains(3) || s2.Cardinality() != 2 {
		t.Error("Difference of overlapping sets mutates non-destination set (argument 1)")
	}
}

func TestSubset(t *testing.T) {
	s1 := set.NewSet()
	s2 := set.NewSet()

	s2.Add(1)
	s2.Add(2)

	if !set.Subset(s1, s2) {
		t.Error("Empty set is not subset of another set")
	}

	s1.Add(1)
	if !set.Subset(s1, s2) {
		t.Errorf("Set {1} is not subset of {1,2}")
	}

	s1.Add(2)
	if !set.Subset(s1, s2) {
		t.Errorf("Set {1,2} is not subset of equal set {1,2}")
	}

	s1.Add(3)
	if set.Subset(s1, s2) {
		t.Error("Set {1,2,3} registers as subset of {1,2}")
	}

	s2.Add(4)
	if set.Subset(s1, s2) {
		t.Error("Set {1,2,3} registers as subset of {1,2,4}")
	}

	s2.Add(5)
	if set.Subset(s1, s2) {
		t.Error("Set {1,2,3} registers as subset of {1,2,4,5}")
	}

	if !set.Subset(s1, s1) || !set.Subset(s2, s2) {
		t.Error("Sets don't register as subsets of temselves")
	}
}

func TestProperSubset(t *testing.T) {
	s1 := set.NewSet()
	s2 := set.NewSet()

	s2.Add(1)
	s2.Add(2)

	if !set.ProperSubset(s1, s2) {
		t.Error("Empty set is not proper subset of another non-empty set")
	}

	s1.Add(1)
	if !set.ProperSubset(s1, s2) {
		t.Errorf("Set {1} is not proper subset of {1,2}")
	}

	s1.Add(2)
	if set.ProperSubset(s1, s2) {
		t.Errorf("Set {1,2} is proper subset of equal set {1,2}")
	}

	s1.Add(3)
	if set.ProperSubset(s1, s2) {
		t.Error("Set {1,2,3} registers as proper subset of {1,2}")
	}

	s2.Add(4)
	if set.ProperSubset(s1, s2) {
		t.Error("Set {1,2,3} registers as proper subset of {1,2,4}")
	}

	s2.Add(5)
	if set.ProperSubset(s1, s2) {
		t.Error("Set {1,2,3} registers as proper subset of {1,2,4,5}")
	}

	if set.ProperSubset(s1, s1) || set.ProperSubset(s2, s2) {
		t.Error("Sets register as proper subsets of themselves")
	}

	s3 := set.NewSet()

	if set.ProperSubset(s3, s3) {
		t.Error("Empty set registers as proper subset of itself")
	}
}
