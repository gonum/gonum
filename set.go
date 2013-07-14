package discrete

// On one hand, using an interface{} as a key works on some levels.
// On the other hand, from experience, I can say that working with interface{} is a pain
// so I don't like it in an API. An alternate idea is to make Set an interface with a method that allows you to GRAB a map[interface{}]struct{} from
// the implementation, but that adds a lot of calls and needless operations, making the library slower
//
// Another point, using an interface{} may be pointless because a map key MUST have == and != defined, limiting the possible keys anyway (for instance, if you had a set of [3]floats I don't think it will do a deep
// comparison, making it rather pointless). Also, keying with a float will mean it does a strict == with the floats, possibly causing bad behavior. It may be best to just make it a map[int]struct{}. Thoughts?
type Set struct {
	data map[interface{}]struct{}
	id   uint
}

// I highly doubt we have to worry about running out of IDs, but we could add a little reclaimID function if we're worried
var globalid uint = 0

// For cleanliness
var flag struct{} = struct{}{}

func NewSet() Set {
	defer func() { globalid++ }()
	return Set{
		data: make(map[interface{}]struct{}, 0),
		id:   globalid,
	}
}

// Reverts the set to the empty set without reallocating
func (s1 Set) Clear() Set {
	for el, _ := range s1.data {
		delete(s1.data, el)
	}

	return s1
}

// Ensures a perfect copy from s1 to dest (meaning the sets will be equal)
func (s1 Set) CopyTo(dest Set) Set {
	if s1.id == dest.id {
		return dest
	}

	if len(dest.data) > 0 {
		for el := range dest.data {
			delete(dest.data, el)
		}
	}

	for el := range s1.data {
		dest.data[el] = flag
	}

	return dest
}

func (s1 Set) Equal(s2 Set) bool {
	if s1.id == s2.id {
		return true
	} else if len(s1.data) != len(s2.data) {
		return false
	}

	for el := range s1.data {
		if _, ok := s2.data[el]; !ok {
			return false
		}
	}

	return true
}

func (s1 Set) Union(dest, s2 Set) Set {
	if s1.id == s2.id {
		return s1.CopyTo(dest)
	}

	if s1.id != dest.id && s2.id != dest.id {
		dest.Clear()
	}

	if dest.id != s1.id {
		for el := range s1.data {
			dest.data[el] = flag
		}
	}

	if dest.id != s2.id {
		for el := range s2.data {
			dest.data[el] = flag
		}
	}

	return dest
}

func (s1 Set) Intersection(dest, s2 Set) Set {
	var swap Set

	if s1.id == s2.id {
		return s1.CopyTo(dest)
	} else if s1.id == dest.id {
		swap = s2
	} else if s2.id == dest.id {
		swap = s1
	} else {
		dest.Clear()

		if len(s1.data) > len(s2.data) {
			s1, s2 = s2, s1
		}
		for el := range s1.data {
			if _, ok := s2.data[el]; ok {
				dest.data[el] = flag
			}
		}

		return dest
	}

	for el := range dest.data {
		if _, ok := swap.data[el]; !ok {
			delete(dest.data, el)
		}
	}

	return dest
}

func (s1 Set) Diff(dest, s2 Set) Set {
	if s1.id == s2.id {
		return dest.Clear()
	} else if s2.id == dest.id {
		tmp := NewSet()

		return s1.Diff(tmp, s2).CopyTo(dest)
	} else if s1.id == dest.id {
		for el := range dest.data {
			if _, ok := s2.data[el]; ok {
				delete(dest.data, el)
			}
		}

	} else {
		dest.Clear()
		for el := range s1.data {
			if _, ok := s2.data[el]; !ok {
				dest.data[el] = flag
			}
		}
	}

	return dest
}

// Are Add/Remove necessary?
func (s1 Set) Add(element interface{}) {
	s1.data[element] = flag
}

func (s1 Set) Remove(element interface{}) {
	delete(s1.data, element)
}
