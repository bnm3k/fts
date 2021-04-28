package main

type intSet map[uint32]struct{}

func (s intSet) add(i uint32) {
	s[i] = struct{}{}
}

func (s intSet) len() int {
	return len(s)
}

func (s intSet) clone() intSet {
	set := make(map[uint32]struct{})
	for n := range s {
		set[n] = struct{}{}
	}
	return set
}

func (s intSet) filterOutDifference(other intSet) {
	for n := range s {
		if _, ok := other[n]; !ok {
			delete(s, n)
		}
	}
}

func intersectAll(sets ...intSet) intSet {
	if len(sets) == 0 {
		return nil
	}
	res := sets[0].clone()
	for _, set := range sets[1:] {
		res.filterOutDifference(set)
	}
	return res
}
