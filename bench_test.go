package main

import (
	"fmt"
	"testing"

	bsmt "github.com/bnb-chain/zkbnb-smt"
	fuzz "github.com/google/gofuzz"
)

func BenchmarkWrite1(b *testing.B) {
	fmt.Println("benchmark write 1")
	for _, env := range prepareEnv() {
		smt, _ := initSMT(env)
		var buf []byte
		f := fuzz.New().NilChance(0)

		for i := 0; i < b.N; i++ {
			f.Fuzz(&buf)
			smt.Set(10, buf)
			smt.Commit(nil)
		}
	}
}

// it may fail when concurrency = 8
func BenchmarkWrite2(b *testing.B) {
	fmt.Println("benchmark write 2")
	for _, env := range prepareEnv() {
		smt, _ := initSMT(env)
		f := fuzz.New().NilChance(0).Funcs(func(item *bsmt.Item, c fuzz.Continue) {
			item.Key = uint64(c.Intn(4096))
			c.Fuzz(&item.Val)
		})
		batch := 10
		items := []bsmt.Item{}

		for i := 0; i < b.N; i++ {
			buf := bsmt.Item{}
			f.Fuzz(&buf)
			items = append(items, buf)
			if len(items) >= batch || i == b.N-1 {
				smt.MultiSet(items)
				smt.Commit(nil)
				items = []bsmt.Item{}
			}
		}
	}
}

func BenchmarkRead1(b *testing.B) {
	fmt.Println("benchmark read 1")
	for _, env := range prepareEnv() {
		smt, _ := initSMT(env)
		opts, index := generateSetOperations(4096)
		execOperations(smt, opts)

		opts = generateGetOperations(b.N, index)

		for i := 0; i < b.N; i++ {
			smt.Get(opts[i].key, nil)
		}
	}
}

func BenchmarkRead2(b *testing.B) {
	fmt.Println("benchmark read 2")
	for _, env := range prepareEnv() {
		smt, _ := initSMT(env)
		opts, index := generateSetOperations(4096)
		execOperations(smt, opts)

		opts = generateGetOperations(b.N, index)

		for i := 0; i < b.N; i++ {
			smt.GetProof(opts[i].key)
		}
	}
}

func BenchmarkCalc1(b *testing.B) {
	fmt.Println("benchmark calc 1")
	for _, env := range prepareEnv() {
		smt, _ := initSMT(env)
		opts, index := generateSetOperations(4096)
		execOperations(smt, opts)

		opts = generateGetOperations(b.N, index)

		for i := 0; i < b.N; i++ {
			pf, _ := smt.GetProof(opts[i].key)
			smt.VerifyProof(opts[i].key, pf)
		}
	}
}
