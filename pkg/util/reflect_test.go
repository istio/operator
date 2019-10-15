// Copyright 2019 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

// TODO: add missing unit tests (istio/istio#17246).

// to ptr conversion utility functions
func toInt8Ptr(i int8) *int8 { return &i }

var _ = Describe("IsValueNil", func() {
	It("Should be true when nil", func() {
		Expect(IsValueNil(nil)).To(BeTrue())
	})
	It("Should be true when nil", func() {
		Expect(IsValueNil(nil)).To(BeTrue())
	})
	It("Should be true when ptr", func() {
		Expect(IsValueNil((*int)(nil))).To(BeTrue())
	})
	It("Should be true when map", func() {
		Expect(IsValueNil(map[int]int(nil))).To(BeTrue())
	})
	It("Should be true when slice", func() {
		Expect(IsValueNil([]int(nil))).To(BeTrue())
	})
	It("Should be true when interface", func() {
		Expect(IsValueNil(interface{}(nil))).To(BeTrue())
	})

	It("Should be false when ptr", func() {
		Expect(IsValueNil(toInt8Ptr(42))).To(BeFalse())
	})
	It("Should be false when map", func() {
		Expect(IsValueNil(map[int]int{42: 42})).To(BeFalse())
	})
	It("Should be false when slice", func() {
		Expect(IsValueNil([]int{1, 2, 3})).To(BeFalse())
	})
	It("Should be false when interface", func() {
		Expect(IsValueNil(interface{}(42))).To(BeFalse())
	})
})

var _ = Describe("IsValueNilOrDefault", func() {
	It("Should be true when nil", func() {
		Expect(IsValueNilOrDefault(nil)).To(BeTrue())
	})
	It("Should be true when ptr", func() {
		Expect(IsValueNilOrDefault((*int)(nil))).To(BeTrue())
	})
	It("Should be true when map", func() {
		Expect(IsValueNilOrDefault(map[int]int(nil))).To(BeTrue())
	})
	It("Should be true when slice", func() {
		Expect(IsValueNilOrDefault([]int(nil))).To(BeTrue())
	})
	It("Should be true when interface", func() {
		Expect(IsValueNilOrDefault(interface{}(nil))).To(BeTrue())
	})
	It("Should be true when int(0)", func() {
		Expect(IsValueNilOrDefault(int(0))).To(BeTrue())
	})
	It("Should be true when \"\"", func() {
		Expect(IsValueNilOrDefault("")).To(BeTrue())
	})
	It("Should be true when false", func() {
		Expect(IsValueNilOrDefault(false)).To(BeTrue())
	})

	// It("Should be true when ptr to ptr", func() {
	// 	i := 32
	// 	ip := &i
	// 	Expect(IsValueNilOrDefault(&ip)).To(BeTrue())
	// })
})

var _ = Describe("IsValueFuncs", func() {
	testInt := int(42)
	testStruct := struct{}{}
	testSlice := []bool{}
	testMap := map[bool]bool{}
	var testNilSlice []bool
	var testNilMap map[bool]bool

	allValues := []interface{}{nil, testInt, &testInt, testStruct, &testStruct, testNilSlice, testSlice, &testSlice, testNilMap, testMap, &testMap}

	DescribeTable("function",
		func(function func(v reflect.Value) bool, okValues []interface{}) {
			for _, v := range allValues {
				Context(fmt.Sprintf("With %v", v), func() {
					Expect(function(reflect.ValueOf(v))).To(Equal(isInListOfInterface(okValues, v)))
				})
			}
		},
		Entry("IsValuePtr", IsValuePtr, []interface{}{&testInt, &testStruct, &testSlice, &testMap}),
		Entry("IsValueStruct", IsValueStruct, []interface{}{testStruct}),
		Entry("IsValueInterface", IsValueInterface, []interface{}{}),
		Entry("IsValueStructPtr", IsValueStructPtr, []interface{}{&testStruct}),
		Entry("IsValueMap", IsValueMap, []interface{}{testNilMap, testMap}),
		Entry("IsValueSlice", IsValueSlice, []interface{}{testNilSlice, testSlice}),
		Entry("IsValueScalar", IsValueScalar, []interface{}{testInt, &testInt}),
	)
})

var _ = Describe("ValuesAreSameType", func() {
	type EnumType int64
	DescribeTable("same types",
		func(inV1, inV2 interface{}, want bool) {
			got := ValuesAreSameType(reflect.ValueOf(inV1), reflect.ValueOf(inV2))
			Expect(got).To(Equal(want))
		},
		Entry("success both are int32 types", int32(42), int32(43), true),
		Entry("fail unmatching int types", int16(42), int32(43), false),
		Entry("fail unmatching int and string type", int32(42), "42", false),
		Entry("fail EnumType and int64 types", EnumType(42), int64(43), false),
	)
})

var _ = Describe("IsTypeFuncs", func() {
	testInt := int(42)
	testStruct := struct{}{}
	testSlice := []bool{}
	testSliceOfInterface := []interface{}{}
	testMap := map[bool]bool{}
	var testNilSlice []bool
	var testNilMap map[bool]bool

	allTypes := []interface{}{nil, testInt, &testInt, testStruct, &testStruct, testNilSlice,
		testSlice, &testSlice, testSliceOfInterface, testNilMap, testMap, &testMap}

	DescribeTable("function",
		func(function func(v reflect.Type) bool, okTypes []interface{}) {
			for _, v := range allTypes {
				Expect(function(reflect.TypeOf(v))).To(Equal(isInListOfInterface(okTypes, v)))
			}
		},
		Entry("IsTypeStructPtr", IsTypeStructPtr, []interface{}{&testStruct}),
		Entry("IsTypeSlicePtr", IsTypeSlicePtr, []interface{}{&testSlice}),
		Entry("IsTypeMap", IsTypeMap, []interface{}{testNilMap, testMap}),
		Entry("IsTypeInterface", IsTypeInterface, []interface{}{}),
		Entry("IsTypeSliceOfInterface", IsTypeSliceOfInterface, []interface{}{testSliceOfInterface}),
	)
})

type interfaceContainer struct {
	I anInterface
}

type anInterface interface {
	IsU()
}

type implementsInterface struct {
	A string
}

func (*implementsInterface) IsU() {}

var _ = Describe("IsValueInterface", func() {
	It("returns true when given an interface", func() {
		intf := &interfaceContainer{
			I: &implementsInterface{
				A: "a",
			},
		}
		iField := reflect.ValueOf(intf).Elem().FieldByName("I")
		Expect(IsValueInterface(iField)).To(BeTrue())
	})
})

var _ = Describe("IsTypeInterface", func() {
	It("returns true when given an interface type", func() {
		intf := &interfaceContainer{
			I: &implementsInterface{
				A: "a",
			},
		}
		testIfField := reflect.ValueOf(intf).Elem().Field(0)

		Expect(IsTypeInterface(testIfField.Type())).To(BeTrue())
	})
})

func isInListOfInterface(lv []interface{}, v interface{}) bool {
	for _, vv := range lv {
		if reflect.DeepEqual(vv, v) {
			return true
		}
	}
	return false
}

var _ = Describe("DeleteFromSlicePtr", func() {
	It("deletes from a slice ptr at a given index", func() {
		parentSlice := []int{42, 43, 44, 45}
		var parentSliceI interface{} = parentSlice
		Expect(DeleteFromSlicePtr(&parentSliceI, 1)).To(Succeed())

		wantSlice := []int{42, 44, 45}
		Expect(parentSliceI).To(Equal(wantSlice))
	})

	It("produces an error if the parent isn't a slice ptr", func() {
		badParent := struct{}{}
		wantErr := `deleteFromSlicePtr parent type is *struct {}, must be *[]interface{}`

		Expect(DeleteFromSlicePtr(&badParent, 1)).To(MatchError(wantErr))
	})
})

var _ = Describe("UpdateSlicePtr", func() {
	It("updates an entry at the a given index", func() {
		parentSlice := []int{42, 43, 44, 45}
		var parentSliceI interface{} = parentSlice
		Expect(UpdateSlicePtr(&parentSliceI, 1, 42)).To(Succeed())

		wantSlice := []int{42, 42, 44, 45}
		Expect(parentSliceI).To(Equal(wantSlice))
	})

	It("produces an error if the parent isn't a slice ptr", func() {
		badParent := struct{}{}
		wantErr := `updateSlicePtr parent type is *struct {}, must be *[]interface{}`
		Expect(UpdateSlicePtr(&badParent, 1, 42)).To(MatchError(wantErr))
	})
})

var _ = Describe("InsertIntoMap", func() {
	It("inserts a value into the parent map", func() {
		parentMap := map[int]string{42: "forty two", 43: "forty three"}
		key := 44
		value := "forty four"
		Expect(InsertIntoMap(parentMap, key, value)).To(Succeed())

		wantMap := map[int]string{42: "forty two", 43: "forty three", 44: "forty four"}
		Expect(parentMap).To(Equal(wantMap))
	})

	It("inserts a value into the parent map ptr", func() {
		parentMap := map[int]string{42: "forty two", 43: "forty three"}
		key := 44
		value := "forty four"
		Expect(InsertIntoMap(&parentMap, key, value)).To(Succeed())

		wantMap := map[int]string{42: "forty two", 43: "forty three", 44: "forty four"}
		Expect(parentMap).To(Equal(wantMap))
	})

	It("raises an error if the parent is an unusable type", func() {
		badParent := struct{}{}
		wantErr := `insertIntoMap parent type is *struct {}, must be map`
		Expect(InsertIntoMap(&badParent, 55, "schfifty five")).To(MatchError(wantErr))
	})
})

var _ = Describe("integer functions", func() {
	var (
		allIntTypes     = []interface{}{int(-42), int8(-43), int16(-44), int32(-45), int64(-46)}
		allUintTypes    = []interface{}{uint(42), uint8(43), uint16(44), uint32(45), uint64(46)}
		allIntegerTypes = append(allIntTypes, allUintTypes...)
		nonIntTypes     = []interface{}{nil, "", []int{}, map[string]bool{}}
		allTypes        = append(allIntegerTypes, nonIntTypes...)
	)

	DescribeTable("IsIntKind and IsUintKind",
		func(function func(v reflect.Kind) bool, want []interface{}) {
			var got []interface{}
			for _, v := range allTypes {
				if function(reflect.ValueOf(v).Kind()) {
					got = append(got, v)
				}
			}
			Expect(got).To(Equal(want))
		},
		Entry("ints", IsIntKind, allIntTypes),
		Entry("uints", IsUintKind, allUintTypes),
	)

	Context("ToIntValue", func() {
		It("returns the int value of the given arg, or false", func() {
			var got []int
			for _, v := range allTypes {
				if i, ok := ToIntValue(v); ok {
					got = append(got, i)
				}
			}
			want := []int{-42, -43, -44, -45, -46, 42, 43, 44, 45, 46}
			Expect(got).To(Equal(want))
		})
	})
})
