package help_me_unpack

import (
	"reflect"
	"testing"
)

func BenchmarkUnpack(b *testing.B) {
	hmu := new(HelpMeUnpack)
	input := `{"bytes": "v3rahFGNDeTEqAAA5T/kQ8RUT6jREUlAQEkR0ahPVMQ="}`
	for n := 0; n < b.N; n++ {
		hmu.Solve(input)
	}
}

func TestUnpack(t *testing.T) {
	hmu := new(HelpMeUnpack)
	input := `{"bytes": "v3rahFGNDeTEqAAA5T/kQ8RUT6jREUlAQEkR0ahPVMQ="}`
	result, err := hmu.Solve(input)
	expected := &Output{
		Int:             -2066056513,
		Uint:            3826093393,
		Short:           -22332,
		Float:           456.499176025390625,
		Double:          50.1392107379302,
		BigEndianDouble: 50.1392107379302,
	}
	if err != nil || !reflect.DeepEqual(result, expected) {
		t.Fatalf(`Run("%s") = %+v, expected %+v`, input, result, expected)
	}
}
