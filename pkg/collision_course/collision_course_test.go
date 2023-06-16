package collision_course

import (
	"reflect"
	"testing"
)

func TestRun(t *testing.T) {
	input := `{"include": "dd0534e453ae91d2a821d588b8a50a68"}`
	expected := Output{
		Files: []string{
			"ZGQwNTM0ZTQ1M2FlOTFkMmE4MjFkNTg4YjhhNTBhNjgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAP9YdpxOV7/EOBwO7v7ytmTSk2kGp/UVPef23SGZPiM6sE6HKVTTUPJm/zAnAcjbbE5AUxuvwCwVe+ZwWnTb86PpeLVVKQGXWPeUlJ0lAwhG6YHwGmPmfRGn7N7AcC8rC8RhoETgT+qh9xz97MdmwideYxqmE0iFWby1dmNs/m",
			"ZGQwNTM0ZTQ1M2FlOTFkMmE4MjFkNTg4YjhhNTBhNjgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAP9YdpxOV7/EOBwO7v7ytmTSm2kGp/UVPef23SGZPiM6sE6HKVTTUPJm/zAvAcjbbE5AUxuvwCwVe+5wWnTb86PpeLVVKQGXWPeUlJ0lAwhG6YnwGmPmfRGn7N7AcC8rC8RhoETgT+qh9xz95MdmwideYxqmE0iFWby9dmNs/m",
		},
	}
	result, err := Run(input)
	if err != nil || !reflect.DeepEqual(*result, expected) {
		t.Fatalf(`Run("%s") = %+v, expected %+v`, input, result, expected)
	}
}
