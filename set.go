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
	s := Set(make(map[interface{}]struct{}, 0))
	return &s
}

func (s1 *Set) Clear() *Set {
	if len(*s1) == 0 {
		return s1
	}

	(*s1) = *NewSet()

	return s1
}

// Ensures a perfect copy from s1 to dest (meaning the sets will be equal)
func (s1 *Set) CopyTo(dest *Set) *Set {
	if s1 == dest {
		return dest
	}

	if len(*dest) > 0 {
		*(dest) = *NewSet()
	}

	for el := range *s1 {
		(*dest)[el] = flag
	}

	return dest
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

func (dest *Set) Union(s1, s2 *Set) *Set {
	if s1 == s2 {
		return s1.CopyTo(dest)
	}

	if s1 != dest && s2 != dest {
		dest.Clear()
	}

	if dest != s1 {
		for el := range *s1 {
			(*dest)[el] = flag
		}
	}

	if dest != s2 {
		for el := range *s2 {
			(*dest)[el] = flag
		}
	}

	return dest
}

func (dest *Set) Intersection(s1, s2 *Set) *Set {
	var swap *Set

	if s1 == s2 {
		return s1.CopyTo(dest)
	} else if s1 == dest {
		swap = s2
	} else if s2 == dest {
		swap = s1
	} else {
		dest.Clear()

		if len(*s1) > len(*s2) {
			s1, s2 = s2, s1
		}
		for el := range *s1 {
			if _, ok := (*s2)[el]; ok {
				(*dest)[el] = flag
			}
		}

		return dest
	}

	for el := range *dest {
		if _, ok := (*swap)[el]; !ok {
			delete(*dest, el)
		}
	}

	return dest
}

func (dest *Set) Diff(s1, s2 *Set) *Set {
	if s1 == s2 {
		return dest.Clear()
	} else if s2 == dest {
		tmp := NewSet()

		return s1.Diff(tmp, s2).CopyTo(dest)
	} else if s1 == dest {
		for el := range *dest {
			if _, ok := (*s2)[el]; ok {
				delete(*dest, el)
			}
		}

	} else {
		dest.Clear()
		for el := range *s1 {
			if _, ok := (*s2)[el]; !ok {
				(*dest)[el] = flag
			}
		}
	}

	return dest
}

// Are Add/Remove necessary?
func (s1 *Set) Add(element interface{}) {
	(*s1)[element] = flag
}

func (s1 *Set) Remove(element interface{}) {
	delete(*s1, element)
}
