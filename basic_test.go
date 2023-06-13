package main

import (
	"fmt"
	"testing"
)

func TestNormalGetSet(t *testing.T) {
	fmt.Println("Running TestNormalGetSet")
	for _, env := range prepareEnv() {
		smt, _ := initSMT(env)
		opts := []TestOperation{
			{
				key:    100,
				value:  "some text lol",
				method: Set,
			},
			{
				method: Commit,
			},
			{
				key:    100,
				method: Get,
				value:  "some text lol",
			},
		}
		if execOperations(smt, opts) != nil {
			t.Errorf("Assertion failed")
		}
	}
}

func TestVersionSet(t *testing.T) {
	fmt.Println("Running TestVersionSet")
	for _, env := range prepareEnv() {
		smt, _ := initSMT(env)
		opts := []TestOperation{
			{
				key:     10,
				value:   "some other text lol",
				version: 5,
				method:  SetWithVersion,
			},
			{
				method:  CommitWithNewVersion,
				version: 5,
			},
			{
				key:    10,
				method: Get,
				value:  "some other text lol",
			},
		}
		if execOperations(smt, opts) != nil {
			t.Errorf("Assertion failed")
		}
	}
}

func TestProof(t *testing.T) {
	fmt.Println("Running TestProof")
	for _, env := range prepareEnv() {
		smt, _ := initSMT(env)
		opts := []TestOperation{
			{
				key:    105,
				value:  "some other text lolol",
				method: Set,
			},
			{
				method: Commit,
			},
			{
				key:    105,
				method: GetProof,
			},
		}
		if execOperations(smt, opts) != nil {
			t.Errorf("Assertion failed")
		}
	}
}
