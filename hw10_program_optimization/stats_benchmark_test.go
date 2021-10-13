package hw10programoptimization

import (
	"archive/zip"
	"testing"
)

// go test -v -bench=BenchmarkGetDomainStat -benchmem stats_benchmark_test.go -count=5 .
func BenchmarkGetDomainStat(b *testing.B) {
	b.StopTimer()

	r, _ := zip.OpenReader("../testdata/users.dat.zip")
	defer r.Close()

	data, _ := r.File[0].Open()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetDomainStat(data, "biz")
	}
	b.StopTimer()
}
