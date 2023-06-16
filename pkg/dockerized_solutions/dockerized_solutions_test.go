package dockerized_solutions

// func BenchmarkRun(b *testing.B) {
// 	input := `{"password": "JOSHUA8505", "salt": "evA4dMr+aFrT8U1io7k=", "pbkdf2": {"rounds": 500000, "hash": "sha256"}, "scrypt": {"N": 262144, "r": 16, "p": 2, "buflen": 32, "_control": "b19a18ea8a50a861d08eb94be602f6cbfe67ab98d2021400a3b83fbe3b8ba698"}}`
// 	for n := 0; n < b.N; n++ {
// 		Run(input)
// 	}
// }

// func TestRun(t *testing.T) {
// 	input := `{"password": "JOSHUA8505", "salt": "evA4dMr+aFrT8U1io7k=", "pbkdf2": {"rounds": 500000, "hash": "sha256"}, "scrypt": {"N": 262144, "r": 16, "p": 2, "buflen": 32, "_control": "b19a18ea8a50a861d08eb94be602f6cbfe67ab98d2021400a3b83fbe3b8ba698"}}`
// 	result, err := Run(input)
// 	expected := Output{
// 		Sha256: "df014bae0f167f018ee6e3e9d1cf169fb2704fda21a5def87357e88ea346e7bc",
// 		Hmac:   "b75c76423eead7160dcf108fe88ead9f110944524667f5646f1a34329b7f0924",
// 		Pbkdf2: "4f2202310ebb8ec5e8d9dd161e7558ae9155b3c9de24494b40716b1a29f877c2",
// 		Scrypt: "348bb0993a69d57083865d9639df866b9c65fa2bbbb02aebd45af689f49cd34b",
// 	}
// 	if err != nil || !reflect.DeepEqual(*result, expected) {
// 		t.Fatalf(`Run("%s") = %+v, expected %+v`, input, result, expected)
// 	}
// }
