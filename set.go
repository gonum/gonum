package discrete

import ()

// On one hand, using an interface{} as a key works on some levels.
// On the other hand, from experience, I can say that working with interface{} is a pain
// so I don't like it in an API. An alternate idea is to make Set an interface with a method that allows you to GRAB a map[interface{}]struct{} from
// the implementation, but that adds a lot of calls and needless operations, making the library slower
//
// Another point, using an interface{} may be pointless because a map key MUST have == and != defined, limiting the possible keys anyway (for instance, if you had a set of [3]floats I don't think it will do a deep
// comparison, making it rather pointless). Also, keying with a float will mean it does a strict == with the floats, possibly causing bad behavior. It may be best to just make it a map[int]struct{}. Thoughts?
type Set map[interface{}]struct{}

// I highly doubt we have to worry about running out of IDs, but we could add a little reclaimID function if we're worried
var globalid uint64 = 0

// For cleanliness
var flag struct{} = struct{}{}

func NewSet() *Set {
	s := make(Set)
	return &s
}

func (s1 *Set) Clear() *Set {
	if len(*s1) == 0 {
		return s1
	}

	(*s1) = make(Set)

	return s1
}

// Ensures a perfect copy from s1 to dst (meaning the sets will be equal)
func (dst *Set) Copy(src *Set) *Set {
	if src == dst {
		return dst
	}

	if len(*dst) > 0 {
		*(dst) = *NewSet()
	}

	for el := range *src {
		(*dst)[el] = flag
	}

	return dst
}

// If every element in s1 is also in s2 (and vice versa), the sets are deemed equal
func Equal(s1, s2 *Set) bool {
	if s1 == s2 {
		return true
	} else if len(*s1) != len(*s2) {
		return false
	}

	for el := range *s1 {
		if _, ok := (*s2)[el]; !ok {
			return false
		}
	}

	return true
}

// Takes the union of s1 and s2, and stores it in dst.
//
// The union of two sets, s1 and s2, is the set containing all the elements of each, for instance:
//
//     {a,b,c} UNION {d,e,f} = {a,b,c,d,e,f}
//
// Since sets may not have repetition, unions of two sets that overlap do not contain repeat elements, that is:
//
//     {a,b,c} UNION {b,c,d} = {a,b,c,d}
func (dst *Set) Union(s1, s2 *Set) *Set {
	if s1 == s2 {
		return dst.Copy(s1)
	}

	if s1 != dst && s2 != dst {
		dst.Clear()
	}

	if dst != s1 {
		for el := range *s1 {
			(*dst)[el] = flag
		}
	}

	if dst != s2 {
		for el := range *s2 {
			(*dst)[el] = flag
		}
	}

	return dst
}

// Takes the intersection of s1 and s2, and stores it in dst
//
// The intersection of two sets, s1 and s2, is the set containing all the elements shared between the two sets, for instance
//
//     {a,b,c} INTERSECT {b,c,d} = {b,c}
//
// The intersection between a set and itself is itself, and thus effectively a copy operation:
//
//     {a,b,c} INTERSECT {a,b,c} = {a,b,c}
//
// The intersection between two sets that share no elements is the empty set:
//
//     {a,b,c} INTERSECT {d,e,f} = {}
func (dst *Set) Intersection(s1, s2 *Set) *Set {
	var swap *Set

	if s1 == s2 {
		return dst.Copy(s1)
	} else if s1 == dst {
		swap = s2
	} else if s2 == dst {
		swap = s1
	} else {
		dst.Clear()

		if len(*s1) > len(*s2) {
			s1, s2 = s2, s1
		}
		for el := range *s1 {
			if _, ok := (*s2)[el]; ok {
				(*dst)[el] = flag
			}
		}

		return dst
	}

	for el := range *dst {
		if _, ok := (*swap)[el]; !ok {
			delete(*dst, el)
		}
	}

	return dst
}

// Takes the difference (-) of s1 and s2 and stores it in dst.
//
// The difference (-) between two sets, s1 and s2, is all the elements in s1 that are NOT also in s2.
//
//     {a,b,c} - {b,c,d} = {a}
//
// The difference between two identical sets is the empty set:
//
//     {a,b,c} - {a,b,c} = {a}
//
// The difference between two sets with no overlapping elements is s1
//
//     {a,b,c} - {d,e,f} = {a,b,c}
//
// Implementation note: if dst == s2 (meaning they have identical pointers), a temporary set must be used to store the data
// and then copied over, thus s2.Diff(s1,s2) has an extra allocation and may cause worse performance in some cases.
func (dst *Set) Diff(s1, s2 *Set) *Set {
	if s1 == s2 {
		return dst.Clear()
	} else if s2 == dst {
		tmp := NewSet()

		tmp.Diff(s1, s2)
		*dst = *tmp
	} else if s1 == dst {
		for el := range *dst {
			if _, ok := (*s2)[el]; ok {
				delete(*dst, el)
			}
		}

	} else {
		dst.Clear()
		for el := range *s1 {
			if _, ok := (*s2)[el]; !ok {
				(*dst)[el] = flag
			}
		}
	}

	return dst
}

// Returns true if s1 is an improper subset of s2.
//
// An improper subset occurs when every element in s1 is also in s2 OR s1 and s2 are equal:
//
//     {a,b,c}   SUBSET {a,b,c} = true
//     {a,b}     SUBSET {a,b,c} = true
//     {c,d}     SUBSET {a,b,c} = false
//     {a,b,c,d} SUBSET {a,b,c} = false
//
// Special case: The empty set is a subset of everything
//
// 	   {} SUBSET {a,b} = true
//     {} SUBSET {}    = true
//
// In the case where one needs to test if s1 is smaller than s2, but not equal, use ProperSubset
func Subset(s1, s2 *Set) bool {
	if len(*s1) > len(*s2) {
		return false
	} else if s1 == s2 {
		return true
	} else if len(*s1) == 0 {
		return true
	}

	for _, el := range *s1 {
		if _, ok := (*s2)[el]; !ok {
			return false
		}
	}

	return true
}

// Returns true if s1 is a proper subset of s2.
// A proper subset is when every element of s1 is in s2, but s1 is smaller than s2 (i.e. they are not equal):
//
//     {a,b,c}   PROPER SUBSET {a,b,c} = false
//     {a,b}     PROPER SUBSET {a,b,c} = true
//     {c,d}     PROPER SUBSET {a,b,c} = false
//     {a,b,c,d} PROPER SUBSET {a,b,c} = false
//
// Special case: The empty set is a proper subset of everything (except itself):
//
//      {} PROPER SUBSET {a,b} = true
//      {} PROPER SUBSET {}    = false
//
// When equality is allowed, use Subset
func ProperSubset(s1, s2 *Set) bool {
	if len(*s1) >= len(*s2) {
		return false
	} else if len(*s1) == 0 {
		return true
	} // We can eschew the s1 == s2 because if they are the same their lens are equal anyway

	for _, el := range *s1 {
		if _, ok := (*s2)[el]; !ok {
			return false
		}
	}

	return true
}

// Returns true if el is an element of s.
func (s *Set) Contains(el interface{}) bool {
	_, ok := (*s)[el]
	return ok
}

// Adds the element el to s1
func (s1 *Set) Add(element interface{}) {
	(*s1)[element] = flag
}

// Removes the element el from s1
func (s1 *Set) Remove(element interface{}) {
	delete(*s1, element)
}

// Returns the number of elements in s1
func (s1 *Set) Cardinality() int {
	return len(*s1)
}

func (s1 *Set) Elements() (els []interface{}) {
	els = make([]interface{}, 0, len(*s1))
	for _, el := range *s1 {
		els = append(els, el)
	}

	return els
}

/* Should probably be re-implemented as a tree later */

// A disjoint set is a collection of non-overlapping sets. That is, for any two sets in the disjoint set, their intersection is the empty set
//
// A disjoint set has three principle operations: Make Set, Find, and Union.
//
// Make set creates a new set for an element (presuming it does not already exist in any set in the disjoint set), Find finds the set containing that element (if any),
// and Union merges two sets in the disjoint set. In general, algorithms operating on disjoint sets are "union-find" algorithms, where two sets are found with Find, and then joined with Union.
//
// A concrete example of a union-find algorithm can be found as discrete.Kruskal -- which unions two sets when an edge is created between two vertices, and refuses to make an edge between two vertices if they're part of the same set.
type DisjointSet struct {
	master  *Set
	subsets []*Set
}

func NewDisjointSet() *DisjointSet {
	return &DisjointSet{NewSet(), make([]*Set, 0)}
}

func (ds *DisjointSet) MasterSet() *Set {
	return ds.master
}

// If the element isn't already somewhere in there, adds it to the master set and its own tiny set
func (ds *DisjointSet) MakeSet(el interface{}) {
	if ds.master.Contains(el) {
		return
	}
	ds.master.Add(el)
	ns := NewSet()
	ns.Add(el)
	ds.subsets = append(ds.subsets, ns)
}

// Returns the set the element belongs to, or nil if none
func (ds *DisjointSet) Find(el interface{}) *Set {
	if !ds.master.Contains(el) {
		return nil
	}

	for _, subset := range ds.subsets {
		if subset.Contains(el) {
			return subset
		}
	}

	return nil
}

// Unions two subsets within the DisjointSet
//
// If either s1 or s2 do not appear in the disjoint set (meaning their pointers, deep equality is not tested),
// the function will return without doing anything. Finding sets to perform a union on is typically done with Find.
func (ds *DisjointSet) Union(s1, s2 *Set) {
	if s1 == s2 {
		return
	}
	s1Found, s2Found := false, false

	newSubsetList := make([]*Set, 0, len(ds.subsets)-1)

	for _, subset := range ds.subsets {
		if s1 == subset {
			s1Found = true
			continue
		} else if s2 == subset {
			s2Found = true
			continue
		}

		newSubsetList = append(newSubsetList, subset)
	}

	if s1Found && s2Found {
		newSubsetList = append(newSubsetList, s1.Union(s1, s2))
		ds.subsets = newSubsetList
	}
}
