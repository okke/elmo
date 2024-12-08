package slice

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMap(t *testing.T) {
	s := []int{1, 2, 3}
	m := Map(s, func(i int) string {
		return fmt.Sprintf("%d", i)
	})

	if expected := []string{"1", "2", "3"}; !reflect.DeepEqual(m, expected) {
		t.Errorf("expected %v but got %v", expected, m)
	}
}
