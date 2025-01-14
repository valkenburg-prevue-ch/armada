package schedulerobjects

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestAllocatableByPriorityAndResourceType(t *testing.T) {
	tests := map[string]struct {
		Priorities     []int32
		UsedAtPriority int32
		Resources      map[string]resource.Quantity
	}{
		"lowest priority": {
			Priorities:     []int32{1, 5, 10},
			UsedAtPriority: 1,
			Resources: map[string]resource.Quantity{
				"cpu": resource.MustParse("1"),
				"gpu": resource.MustParse("2"),
			},
		},
		"mid priority": {
			Priorities:     []int32{1, 5, 10},
			UsedAtPriority: 5,
			Resources: map[string]resource.Quantity{
				"cpu": resource.MustParse("1"),
				"gpu": resource.MustParse("2"),
			},
		},
		"highest priority": {
			Priorities:     []int32{1, 5, 10},
			UsedAtPriority: 10,
			Resources: map[string]resource.Quantity{
				"cpu": resource.MustParse("1"),
				"gpu": resource.MustParse("2"),
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := NewAllocatableByPriorityAndResourceType(tc.Priorities, tc.Resources)
			assert.Equal(t, len(tc.Priorities), len(m))

			m.MarkAllocated(tc.UsedAtPriority, ResourceList{Resources: tc.Resources})
			for resourceType, quantity := range tc.Resources {
				for _, p := range tc.Priorities {
					actual := m.Get(p, resourceType)
					if p > tc.UsedAtPriority {
						assert.Equal(t, 0, quantity.Cmp(actual))
					} else {
						expected := resource.MustParse("0")
						assert.Equal(t, 0, expected.Cmp(actual))
					}
				}
			}

			m.MarkAllocatable(tc.UsedAtPriority, ResourceList{Resources: tc.Resources})
			for resourceType, quantity := range tc.Resources {
				for _, p := range tc.Priorities {
					actual := m.Get(p, resourceType)
					assert.Equal(t, 0, quantity.Cmp(actual))
				}
			}
		})
	}
}

func TestAssignedByPriorityAndResourceType(t *testing.T) {
	tests := map[string]struct {
		Priorities     []int32
		UsedAtPriority int32
		Resources      map[string]resource.Quantity
	}{
		"lowest priority": {
			Priorities:     []int32{1, 5, 10},
			UsedAtPriority: 1,
			Resources: map[string]resource.Quantity{
				"cpu": resource.MustParse("1"),
				"gpu": resource.MustParse("2"),
			},
		},
		"mid priority": {
			Priorities:     []int32{1, 5, 10},
			UsedAtPriority: 5,
			Resources: map[string]resource.Quantity{
				"cpu": resource.MustParse("1"),
				"gpu": resource.MustParse("2"),
			},
		},
		"highest priority": {
			Priorities:     []int32{1, 5, 10},
			UsedAtPriority: 10,
			Resources: map[string]resource.Quantity{
				"cpu": resource.MustParse("1"),
				"gpu": resource.MustParse("2"),
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := NewAllocatedByPriorityAndResourceType(tc.Priorities)
			assert.Equal(t, len(tc.Priorities), len(m))

			m.MarkAllocated(tc.UsedAtPriority, ResourceList{Resources: tc.Resources})
			for resourceType, quantity := range tc.Resources {
				for _, p := range tc.Priorities {
					actual := m.Get(p, resourceType)
					if p <= tc.UsedAtPriority {
						assert.Equal(t, 0, quantity.Cmp(actual))
					} else {
						expected := resource.MustParse("0")
						assert.Equal(t, 0, expected.Cmp(actual))
					}
				}
			}

			m.MarkAllocatable(tc.UsedAtPriority, ResourceList{Resources: tc.Resources})
			for resourceType := range tc.Resources {
				for _, p := range tc.Priorities {
					actual := m.Get(p, resourceType)
					expected := resource.MustParse("0")
					assert.Equal(t, 0, expected.Cmp(actual))
				}
			}
		})
	}
}

func TestResourceListDeepCopy(t *testing.T) {
	rl := ResourceList{
		Resources: map[string]resource.Quantity{
			"foo": resource.MustParse("1"),
		},
	}
	rlCopy := rl.DeepCopy()
	rlCopy.Resources["bar"] = resource.MustParse("2")
	q := rlCopy.Resources["foo"]
	q.Add(resource.MustParse("10"))
	assert.True(
		t,
		rl.Equal(ResourceList{
			Resources: map[string]resource.Quantity{
				"foo": resource.MustParse("1"),
			},
		}),
	)
}

func TestV1ResourceListConversion(t *testing.T) {
	rl := ResourceList{
		Resources: map[string]resource.Quantity{
			"foo": resource.MustParse("1"),
		},
	}
	rlCopy := rl.DeepCopy()
	v1rl := V1ResourceListFromResourceList(rlCopy)
	rlCopy.Resources["bar"] = resource.MustParse("2")
	q := rlCopy.Resources["foo"]
	q.Add(resource.MustParse("10"))

	rl = ResourceListFromV1ResourceList(v1rl)
	assert.True(
		t,
		rl.Equal(ResourceList{
			Resources: map[string]resource.Quantity{
				"foo": resource.MustParse("1"),
			},
		}),
	)

	v1rlCopy := V1ResourceListFromResourceList(rl)
	assert.True(t, maps.Equal(v1rlCopy, v1rl))
}
