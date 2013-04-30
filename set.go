package discrete

// On one hand, using an interface{} as a key works on some levels.
// On the other hand, from experience, I can say that working with interface{} is a pain
// so I don't like it in an API. An alternate idea is to make Set an interface with a method that allows you to GRAB a map[interface{}]struct{} from
// the implementation, but that adds a lot of calls and needless operations, making the library slower
type Set map[interface{}]struct{}

// For cleanliness
var flag struct{} = struct{}{}

/* All functions copy the sets. Shouldn't be too much of a problem. I could be persuaded to change thing to be in-place though I guess. Partially depends on what we decide for other repositories */

// Add elements of set 2 to set 1
func (s1 Set) Union(s2 Set) Set {
	
	for el := range s2 {
		s1[el] = flag
	}
	
	return s2
}

// Iterate over the smaller set, remove all elements not also in the bigger set
func (s1 Set) Intersection(s2 Set) Set {
	// Iterate over the smaller set
	if len(s1) > len(s2) {
		s1,s2 = s2,s1
	} //Q: Should this just be an if-else? Style-wise I like the swap, not sure if it's performance heavy

	for el := range s1 {
		if s2[el] != flag {
			delete(el, s1)
		}
	}
	
	return s1
}

// Del
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