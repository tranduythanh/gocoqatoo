package helpers

// Union returns the union of two sets.
func Union(setA, setB map[interface{}]struct{}) map[interface{}]struct{} {
	unionSet := make(map[interface{}]struct{})
	for k := range setA {
		unionSet[k] = struct{}{}
	}
	for k := range setB {
		unionSet[k] = struct{}{}
	}
	return unionSet
}

// Intersection returns the intersection of two sets.
func Intersection(setA, setB map[interface{}]struct{}) map[interface{}]struct{} {
	intersectionSet := make(map[interface{}]struct{})
	for k := range setA {
		if _, exists := setB[k]; exists {
			intersectionSet[k] = struct{}{}
		}
	}
	return intersectionSet
}

// Difference returns the difference of two sets (setA - setB).
func Difference(setA, setB map[interface{}]struct{}) map[interface{}]struct{} {
	diffSet := make(map[interface{}]struct{})
	for k := range setA {
		if _, exists := setB[k]; !exists {
			diffSet[k] = struct{}{}
		}
	}
	return diffSet
}

// SymDifference returns the symmetric difference of two sets.
func SymDifference(setA, setB map[interface{}]struct{}) map[interface{}]struct{} {
	tmpA := Union(setA, setB)
	tmpB := Intersection(setA, setB)
	return Difference(tmpA, tmpB)
}

// IsSubset checks if setA is a subset of setB.
func IsSubset(setA, setB map[interface{}]struct{}) bool {
	for k := range setA {
		if _, exists := setB[k]; !exists {
			return false
		}
	}
	return true
}

// IsSuperset checks if setA is a superset of setB.
func IsSuperset(setA, setB map[interface{}]struct{}) bool {
	return IsSubset(setB, setA)
}

// ToString converts the set to a string.
func ToString(set map[interface{}]struct{}) string {
	var s string
	i := 0
	size := len(set)
	for k := range set {
		i++
		if i == size-1 {
			s += k.(string) + " and "
		} else if i == size {
			s += k.(string)
		} else {
			s += k.(string) + ", "
		}
	}
	return s
}
