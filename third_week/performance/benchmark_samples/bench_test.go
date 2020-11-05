package main

import (
	"strconv"
	"testing"
)

type testStruct struct {
	X int
	Y string
}

func (t *testStruct) ToJSON() ([]byte, error) {
	return []byte(`{"X": ` + strconv.Itoa(t.X) + `, "Y": "` + t.Y + `"}`), nil
}

func BenchmarkToJSON(b *testing.B) {
	tmp := &testStruct{X: 1, Y: "string"}
	js, err := tmp.ToJSON()
	if err != nil {
		b.Fatal(err)
	}
	b.SetBytes(int64(len(js)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := tmp.ToJSON(); err != nil {
			b.Fatal(err)
		}
	}
}
