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

func (s1 *Set) Equal(s2 *Set) bool {
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

func (dst *Set) Diff(s1, s2 *Set) *Set {
	if s1 == s2 {
		return dst.Clear()
	} else if s2 == dst {
		tmp := NewSet()

		return dst.Copy(tmp.Diff(s1, s2))
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

func (s *Set) Contains(el interface{}) bool {
	_, ok := (*s)[el]
	return ok
}

// Are Add/Remove necessary?
func (s1 *Set) Add(element interface{}) {
	(*s1)[element] = flag
}

func (s1 *Set) Remove(element interface{}) {
	delete(*s1, element)
}

func (s1 *Set) Cardinality() int {
	return len(*s1)
}

/* Should probably be re-implemented as a tree later */
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

func (ds *DisjointSet) MakeSet(el interface{}) {
	if ds.master.Contains(el) {
		return
	}
	ds.master.Add(el)
	ns := NewSet()
	ns.Add(el)
	ds.subsets = append(ds.subsets, ns)
}

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
