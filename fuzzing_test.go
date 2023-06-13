package main

import (
	"testing"

	fuzz "github.com/google/gofuzz"
)

func generateSetOperations(n int) ([]TestOperation, map[uint64][]byte) {
	f := fuzz.New().NilChance(0).Funcs(func(elem *TestOperation, c fuzz.Continue) {
		elem.method = Set
		elem.key = uint64(c.Intn(8192) + 10)
		c.Fuzz(&elem.value)
	})
	opts := []TestOperation{}
	index := make(map[uint64][]byte)
	for i := 0; i < n; i++ {
		opt := TestOperation{}
		f.Fuzz(&opt)
		opts = append(opts, opt)
		index[opt.key] = []byte(opt.value)
	}
	return opts, index
}

func generateGetOperations(n int, index map[uint64][]byte) []TestOperation {
	keys := []uint64{}
	for k := range index {
		keys = append(keys, k)
	}
	f := fuzz.New().NilChance(0).Funcs(func(elem *TestOperation, c fuzz.Continue) {
		elem.method = Get
		elem.key = keys[c.Intn(len(keys))]
		elem.value = string(index[elem.key])
	})
	opts := []TestOperation{}
	for i := 0; i < n; i++ {
		opt := TestOperation{}
		f.Fuzz(&opt)
		opts = append(opts, opt)
	}
	return opts
}

func generateProofOperations(n int, index map[uint64][]byte) []TestOperation {
	keys := []uint64{}
	for k := range index {
		keys = append(keys, k)
	}
	f := fuzz.New().NilChance(0).Funcs(func(elem *TestOperation, c fuzz.Continue) {
		elem.method = GetProof
		elem.key = keys[c.Intn(len(keys))]
	})
	opts := []TestOperation{}
	for i := 0; i < n; i++ {
		opt := TestOperation{}
		f.Fuzz(&opt)
		opts = append(opts, opt)
	}
	return opts
}

// TestFuzzing run random operations
func TestFuzzingGetSet(t *testing.T) {
	for _, env := range prepareEnv() {
		smt, _ := initSMT(env)
		// generate lots of set/commit opt
		opts, index := generateSetOperations(2000)
		// end with a commit
		opts = append(opts, TestOperation{
			method: Commit,
		})
		if execOperations(smt, opts) != nil {
			t.Errorf("TestFuzzingSet failed")
		}
		// generate lots of get opt
		opts = generateGetOperations(1000, index)
		if execOperations(smt, opts) != nil {
			t.Errorf("TestFuzzingGet failed")
		}
	}
}

func TestFuzzingProof(t *testing.T) {
	for _, env := range prepareEnv() {
		smt, _ := initSMT(env)
		// generate lots of set/commit opt
		opts, index := generateSetOperations(2000)
		// end with a commit
		opts = append(opts, TestOperation{
			method: Commit,
		})
		if execOperations(smt, opts) != nil {
			t.Errorf("TestFuzzingSet failed")
		}
		// generate lots of proof opt
		opts = generateProofOperations(1000, index)
		if execOperations(smt, opts) != nil {
			t.Errorf("TestFuzzingProof failed")
		}
	}
}
