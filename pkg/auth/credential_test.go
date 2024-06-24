package auth

import "testing"

func BenchmarkGenerateRandomString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateRandomString(32)
	}
}

func BenchmarkGenerateRandomBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateRandomBytes(32)
	}
}

func BenchmarkGenerateRandomURLSafeString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateRandomStringURLSafe(32)
	}
}
