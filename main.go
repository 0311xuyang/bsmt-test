package main

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"time"

	bsmt "github.com/bnb-chain/zkbnb-smt"
	"github.com/bnb-chain/zkbnb-smt/database"
	"github.com/bnb-chain/zkbnb-smt/database/memory"
	"github.com/bnb-chain/zkbnb-smt/metrics"
)

type testEnv struct {
	tag    string
	hasher *bsmt.Hasher
	db     func() (database.TreeDB, error)
	depth  uint8
}

func prepareEnv() []testEnv {
	initMemoryDB := func() (database.TreeDB, error) {
		return memory.NewMemoryDB(), nil
	}
	return []testEnv{
		{
			tag:    "memory",
			hasher: bsmt.NewHasherPool(func() hash.Hash { return sha256.New() }),
			db:     initMemoryDB,
			depth:  16,
		},
	}
}

type Optype int

const (
	Set Optype = iota
	SetWithVersion
	Commit
	CommitWithNewVersion
	Get
	GetProof
	VerifyProof
)

// TestOperation desc
type TestOperation struct {
	key     uint64
	value   string
	version uint64
	method  Optype
	fail    bool
}

// Exec run the operations
func (op TestOperation) Exec(smt bsmt.SparseMerkleTree) error {
	switch op.method {
	case Set:
		return smt.Set(op.key, []byte(op.value))
	case SetWithVersion:
		return smt.SetWithVersion(op.key, []byte(op.value), bsmt.Version(op.version))
	case Commit:
		recentVer := smt.RecentVersion()
		ver, err := smt.Commit(&recentVer)
		fmt.Printf("Commited with version %v\n", ver)
		return err
	case CommitWithNewVersion:
		recentVer := smt.RecentVersion()
		ver, err := smt.CommitWithNewVersion(&recentVer, (*bsmt.Version)(&op.version))
		fmt.Printf("Commited with version %v\n", ver)
		return err
	case Get:
		val, err := smt.Get(op.key, nil)
		if string(val) != op.value {
			return fmt.Errorf("Get assertion failed `%s` != `%s`", string(val), op.value)
		}
		return err
	case GetProof:
		pf, err := smt.GetProof(op.key)
		if !smt.VerifyProof(op.key, pf) {
			return fmt.Errorf("Proof assertion failed %d", op.key)
		}
		return err
	}
	return nil
}

func execOperations(smt bsmt.SparseMerkleTree, operations []TestOperation) error {
	for _, opt := range operations {
		err := opt.Exec(smt)
		// fmt.Printf("Run %v -> err = %v\n", opt.method, err)
		if err != nil {
			return err
		}
	}
	return nil
}

func execOperationsWithSleep(smt bsmt.SparseMerkleTree, operations []TestOperation, d time.Duration) error {
	for _, opt := range operations {
		err := opt.Exec(smt)
		// fmt.Printf("Run %v -> err = %v\n", opt.method, err)
		if err != nil {
			return err
		}
		time.Sleep(d)
	}
	return nil
}

func initSMT(env testEnv) (bsmt.SparseMerkleTree, error) {
	db, err := env.db()
	if err != nil {
		return nil, err
	}

	nilHash := []byte{}
	return bsmt.NewBNBSparseMerkleTree(env.hasher, db, env.depth, nilHash)
}

func initSMTWithMetrics(env testEnv, metrics metrics.Metrics) (bsmt.SparseMerkleTree, error) {
	db, err := env.db()
	if err != nil {
		return nil, err
	}

	nilHash := []byte{}
	return bsmt.NewBNBSparseMerkleTree(env.hasher, db, env.depth, nilHash, bsmt.EnableMetrics(metrics))
}

func main() {
	fmt.Println("running ...")
}
