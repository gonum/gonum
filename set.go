package discrete

// On one hand, using an interface{} as a key works on some levels.
// On the other hand, from experience, I can say that working with interface{} is a pain
// so I don't like it in an API. An alternate idea is to make Set an interface with a method that allows you to GRAB a map[interface{}]struct{} from
// the implementation, but that adds a lot of calls and needless operations, making the library slower
type Set map[interface{}]struct{}

var flag struct{} = struct{}{}

func (s1 Set) Union(s2 Set) Set {
	
	for el := range s2 {
		s1[el] = flag
	}
	
	return s2
}

func (s1 Set) Intersection(s2 Set) (out Set) {
	// Iterate over the smaller set
	if len(s1) > len(s2) {
		s1,s2 = s2,s1
	} //Q: Should this just be an if-else? Style-wise I like the swap, not sure if it's performance heavy
	
	out = make(map[interface{}]struct{}, len(s1))
	
	for el := range s1 {
		if s2[el] == flag {
			out[el] = flag
		}
	}
	
	return out
}

func (s1 Set) Diff(s2 Set) Set {
	// Usually modifying an element being iterated over is a no-no, but go copies the argument to range
	// so this is safe
	for el := range s1 {
		if s2[el] == flag {
			delete(s1, el)
		}
	}
	
	return s1
}

// Are Add/Remove necessary?
func (s1 *Set) Add(element interface{}) {
	(*s1)[element] = flag
}

func (s1 *Set) Remove(element interface{}) {
	delete(*s1, element)
}