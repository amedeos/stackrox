// Code generated by genny. DO NOT EDIT.
// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/mauricelam/genny

package centralsensor

import (
	"sort"
)

// If you want to add a set for your custom type, simply add another go generate line along with the
// existing ones. If you're creating a set for a primitive type, you can follow the example of "string"
// and create the generated file in this package.
// For non-primitive sets, please make the generated code files go outside this package.
// Sometimes, you might need to create it in the same package where it is defined to avoid import cycles.
// The permission set is an example of how to do that.
// You can also specify the -imp command to specify additional imports in your generated file, if required.

// SensorCapability represents a generic type that we want to have a set of.

// SensorCapabilitySet will get translated to generic sets.
type SensorCapabilitySet map[SensorCapability]struct{}

// Add adds an element of type SensorCapability.
func (k *SensorCapabilitySet) Add(i SensorCapability) bool {
	if *k == nil {
		*k = make(map[SensorCapability]struct{})
	}

	oldLen := len(*k)
	(*k)[i] = struct{}{}
	return len(*k) > oldLen
}

// AddMatching is a utility function that adds all the elements that match the given function to the set.
func (k *SensorCapabilitySet) AddMatching(matchFunc func(SensorCapability) bool, elems ...SensorCapability) bool {
	oldLen := len(*k)
	for _, elem := range elems {
		if !matchFunc(elem) {
			continue
		}
		if *k == nil {
			*k = make(map[SensorCapability]struct{})
		}
		(*k)[elem] = struct{}{}
	}
	return len(*k) > oldLen
}

// AddAll adds all elements of type SensorCapability. The return value is true if any new element
// was added.
func (k *SensorCapabilitySet) AddAll(is ...SensorCapability) bool {
	if len(is) == 0 {
		return false
	}
	if *k == nil {
		*k = make(map[SensorCapability]struct{})
	}

	oldLen := len(*k)
	for _, i := range is {
		(*k)[i] = struct{}{}
	}
	return len(*k) > oldLen
}

// Remove removes an element of type SensorCapability.
func (k *SensorCapabilitySet) Remove(i SensorCapability) bool {
	if len(*k) == 0 {
		return false
	}

	oldLen := len(*k)
	delete(*k, i)
	return len(*k) < oldLen
}

// RemoveAll removes the given elements.
func (k *SensorCapabilitySet) RemoveAll(is ...SensorCapability) bool {
	if len(*k) == 0 {
		return false
	}

	oldLen := len(*k)
	for _, i := range is {
		delete(*k, i)
	}
	return len(*k) < oldLen
}

// RemoveMatching removes all elements that match a given predicate.
func (k *SensorCapabilitySet) RemoveMatching(pred func(SensorCapability) bool) bool {
	if len(*k) == 0 {
		return false
	}

	oldLen := len(*k)
	for elem := range *k {
		if pred(elem) {
			delete(*k, elem)
		}
	}
	return len(*k) < oldLen
}

// Contains returns whether the set contains an element of type SensorCapability.
func (k SensorCapabilitySet) Contains(i SensorCapability) bool {
	_, ok := k[i]
	return ok
}

// Cardinality returns the number of elements in the set.
func (k SensorCapabilitySet) Cardinality() int {
	return len(k)
}

// IsEmpty returns whether the underlying set is empty (includes uninitialized).
func (k SensorCapabilitySet) IsEmpty() bool {
	return len(k) == 0
}

// Clone returns a copy of this set.
func (k SensorCapabilitySet) Clone() SensorCapabilitySet {
	if k == nil {
		return nil
	}
	cloned := make(map[SensorCapability]struct{}, len(k))
	for elem := range k {
		cloned[elem] = struct{}{}
	}
	return cloned
}

// Difference returns a new set with all elements of k not in other.
func (k SensorCapabilitySet) Difference(other SensorCapabilitySet) SensorCapabilitySet {
	if len(k) == 0 || len(other) == 0 {
		return k.Clone()
	}

	retained := make(map[SensorCapability]struct{}, len(k))
	for elem := range k {
		if !other.Contains(elem) {
			retained[elem] = struct{}{}
		}
	}
	return retained
}

// Intersect returns a new set with the intersection of the members of both sets.
func (k SensorCapabilitySet) Intersect(other SensorCapabilitySet) SensorCapabilitySet {
	maxIntLen := len(k)
	smaller, larger := k, other
	if l := len(other); l < maxIntLen {
		maxIntLen = l
		smaller, larger = larger, smaller
	}
	if maxIntLen == 0 {
		return nil
	}

	retained := make(map[SensorCapability]struct{}, maxIntLen)
	for elem := range smaller {
		if _, ok := larger[elem]; ok {
			retained[elem] = struct{}{}
		}
	}
	return retained
}

// Union returns a new set with the union of the members of both sets.
func (k SensorCapabilitySet) Union(other SensorCapabilitySet) SensorCapabilitySet {
	if len(k) == 0 {
		return other.Clone()
	} else if len(other) == 0 {
		return k.Clone()
	}

	underlying := make(map[SensorCapability]struct{}, len(k)+len(other))
	for elem := range k {
		underlying[elem] = struct{}{}
	}
	for elem := range other {
		underlying[elem] = struct{}{}
	}
	return underlying
}

// Equal returns a bool if the sets are equal
func (k SensorCapabilitySet) Equal(other SensorCapabilitySet) bool {
	thisL, otherL := len(k), len(other)
	if thisL == 0 && otherL == 0 {
		return true
	}
	if thisL != otherL {
		return false
	}
	for elem := range k {
		if _, ok := other[elem]; !ok {
			return false
		}
	}
	return true
}

// AsSlice returns a slice of the elements in the set. The order is unspecified.
func (k SensorCapabilitySet) AsSlice() []SensorCapability {
	if len(k) == 0 {
		return nil
	}
	elems := make([]SensorCapability, 0, len(k))
	for elem := range k {
		elems = append(elems, elem)
	}
	return elems
}

// GetArbitraryElem returns an arbitrary element from the set.
// This can be useful if, for example, you know the set has exactly one
// element, and you want to pull it out.
// If the set is empty, the zero value is returned.
func (k SensorCapabilitySet) GetArbitraryElem() (arbitraryElem SensorCapability) {
	for elem := range k {
		arbitraryElem = elem
		break
	}
	return arbitraryElem
}

// AsSortedSlice returns a slice of the elements in the set, sorted using the passed less function.
func (k SensorCapabilitySet) AsSortedSlice(less func(i, j SensorCapability) bool) []SensorCapability {
	slice := k.AsSlice()
	if len(slice) < 2 {
		return slice
	}
	// Since we're generating the code, we might as well use sort.Sort
	// and avoid paying the reflection penalty of sort.Slice.
	sortable := &sortableSensorCapabilitySlice{slice: slice, less: less}
	sort.Sort(sortable)
	return sortable.slice
}

// Clear empties the set
func (k *SensorCapabilitySet) Clear() {
	*k = nil
}

// Freeze returns a new, frozen version of the set.
func (k SensorCapabilitySet) Freeze() FrozenSensorCapabilitySet {
	return NewFrozenSensorCapabilitySetFromMap(k)
}

// NewSensorCapabilitySet returns a new thread unsafe set with the given key type.
func NewSensorCapabilitySet(initial ...SensorCapability) SensorCapabilitySet {
	underlying := make(map[SensorCapability]struct{}, len(initial))
	for _, elem := range initial {
		underlying[elem] = struct{}{}
	}
	return underlying
}

type sortableSensorCapabilitySlice struct {
	slice []SensorCapability
	less  func(i, j SensorCapability) bool
}

func (s *sortableSensorCapabilitySlice) Len() int {
	return len(s.slice)
}

func (s *sortableSensorCapabilitySlice) Less(i, j int) bool {
	return s.less(s.slice[i], s.slice[j])
}

func (s *sortableSensorCapabilitySlice) Swap(i, j int) {
	s.slice[j], s.slice[i] = s.slice[i], s.slice[j]
}

// A FrozenSensorCapabilitySet is a frozen set of SensorCapability elements, which
// cannot be modified after creation. This allows users to use it as if it were
// a "const" data structure, and also makes it slightly more optimal since
// we don't have to lock accesses to it.
type FrozenSensorCapabilitySet struct {
	underlying map[SensorCapability]struct{}
}

// NewFrozenSensorCapabilitySetFromMap returns a new frozen set from the set-style map.
func NewFrozenSensorCapabilitySetFromMap(m map[SensorCapability]struct{}) FrozenSensorCapabilitySet {
	if len(m) == 0 {
		return FrozenSensorCapabilitySet{}
	}
	underlying := make(map[SensorCapability]struct{}, len(m))
	for elem := range m {
		underlying[elem] = struct{}{}
	}
	return FrozenSensorCapabilitySet{
		underlying: underlying,
	}
}

// NewFrozenSensorCapabilitySet returns a new frozen set with the provided elements.
func NewFrozenSensorCapabilitySet(elements ...SensorCapability) FrozenSensorCapabilitySet {
	underlying := make(map[SensorCapability]struct{}, len(elements))
	for _, elem := range elements {
		underlying[elem] = struct{}{}
	}
	return FrozenSensorCapabilitySet{
		underlying: underlying,
	}
}

// Contains returns whether the set contains the element.
func (k FrozenSensorCapabilitySet) Contains(elem SensorCapability) bool {
	_, ok := k.underlying[elem]
	return ok
}

// Cardinality returns the cardinality of the set.
func (k FrozenSensorCapabilitySet) Cardinality() int {
	return len(k.underlying)
}

// IsEmpty returns whether the underlying set is empty (includes uninitialized).
func (k FrozenSensorCapabilitySet) IsEmpty() bool {
	return len(k.underlying) == 0
}

// AsSlice returns the elements of the set. The order is unspecified.
func (k FrozenSensorCapabilitySet) AsSlice() []SensorCapability {
	if len(k.underlying) == 0 {
		return nil
	}
	slice := make([]SensorCapability, 0, len(k.underlying))
	for elem := range k.underlying {
		slice = append(slice, elem)
	}
	return slice
}

// AsSortedSlice returns the elements of the set as a sorted slice.
func (k FrozenSensorCapabilitySet) AsSortedSlice(less func(i, j SensorCapability) bool) []SensorCapability {
	slice := k.AsSlice()
	if len(slice) < 2 {
		return slice
	}
	// Since we're generating the code, we might as well use sort.Sort
	// and avoid paying the reflection penalty of sort.Slice.
	sortable := &sortableSensorCapabilitySlice{slice: slice, less: less}
	sort.Sort(sortable)
	return sortable.slice
}
