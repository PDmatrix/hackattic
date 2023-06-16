package mini_miner

import "testing"

var input = `{"difficulty": 14, "block": {"nonce": null, "data": [["6f7313de03396b1d53a92d8a99c7f3de", -40], ["ea2c4d9ebabc4f761dd6cb59d115a91b", -37], ["60b0fad25101d2a3365e40113787bdb4", -13]]}}`

func BenchmarkMiniminer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Run(input)
	}
}

func TestMiniminer(t *testing.T) {
	expected := 24167
	result, err := Run(input)
	if err != nil || result.(int) != expected {
		t.Fatalf(`Run("%s") = %v, expected %v`, input, result.(int), expected)
	}
}
