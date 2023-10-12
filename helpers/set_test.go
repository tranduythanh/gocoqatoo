package helpers

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSetOperations(t *testing.T) {
	Convey("Given two sets", t, func() {

		setA := map[interface{}]struct{}{
			"one":   {},
			"two":   {},
			"three": {},
		}
		setB := map[interface{}]struct{}{
			"three": {},
			"four":  {},
		}

		Convey("When calculating their union", func() {
			result := Union(setA, setB)

			Convey("The result should contain all elements from both sets", func() {
				So(result, ShouldContainKey, "one")
				So(result, ShouldContainKey, "two")
				So(result, ShouldContainKey, "three")
				So(result, ShouldContainKey, "four")
			})
		})

		Convey("When calculating their intersection", func() {
			result := Intersection(setA, setB)

			Convey("The result should only contain elements that exist in both sets", func() {
				So(result, ShouldNotContainKey, "one")
				So(result, ShouldNotContainKey, "two")
				So(result, ShouldContainKey, "three")
				So(result, ShouldNotContainKey, "four")
			})
		})

		Convey("When calculating their difference (setA - setB)", func() {
			result := Difference(setA, setB)

			Convey("The result should contain elements that exist in setA but not in setB", func() {
				So(result, ShouldContainKey, "one")
				So(result, ShouldContainKey, "two")
				So(result, ShouldNotContainKey, "three")
				So(result, ShouldNotContainKey, "four")
			})
		})

		Convey("SymDifference", func() {
			setA := map[interface{}]struct{}{"one": {}, "two": {}}
			setB := map[interface{}]struct{}{"two": {}, "three": {}}

			Convey("Given two sets, it should return their symmetric difference", func() {
				result := SymDifference(setA, setB)

				So(result, ShouldContainKey, "one")
				So(result, ShouldNotContainKey, "two") // because it's common to both sets
				So(result, ShouldContainKey, "three")
			})
		})

		Convey("IsSubset", func() {
			setA := map[interface{}]struct{}{"one": {}}
			setB := map[interface{}]struct{}{"one": {}, "two": {}, "three": {}}

			Convey("Given two sets, it should determine if one is a subset of the other", func() {
				So(IsSubset(setA, setB), ShouldBeTrue)
				So(IsSubset(setB, setA), ShouldBeFalse)
			})
		})

		Convey("IsSuperset", func() {
			setA := map[interface{}]struct{}{"one": {}, "two": {}, "three": {}}
			setB := map[interface{}]struct{}{"two": {}}

			Convey("Given two sets, it should determine if one is a superset of the other", func() {
				So(IsSuperset(setA, setB), ShouldBeTrue)
				So(IsSuperset(setB, setA), ShouldBeFalse)
			})
		})

		Convey("ToString", func() {
			set := map[interface{}]struct{}{"apple": {}, "banana": {}, "cherry": {}}

			Convey("Given a set, it should convert the set to a string", func() {
				result := ToString(set)
				So(result, ShouldContainSubstring, "apple, ")
				So(result, ShouldContainSubstring, "banana and ")
				So(result, ShouldContainSubstring, "cherry")
			})
		})
	})
}
