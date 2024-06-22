package kzg

import (
	// "github.com/protolambda/go-kzg/bls"
	"testing"
	"fmt"
	// "os"
	// "io"
	// "io/ioutil"
	"time"
	"strconv"
)

func TestBitcoinRollups(t *testing.T){
	var beta uint64
	beta = 1 << 23
	var epsilon uint8
	epsilon = 12
	ks := initialzation(beta, epsilon)

	all_commitments, all_poly := ProposerCommit(ks)

	for j := 0; j < len(all_commitments); j++ {
		SaveToFile(all_commitments[j], "./tmp/data" + strconv.Itoa(j))
	}

	content_x := ValidatorSampling()

	start := time.Now()
	all_evaluation, all_proof := ProposerReply(ks, all_poly, content_x, all_commitments)
	duration := time.Since(start)

	fmt.Println(duration)

	all_commitments_read := make([]string, num_of_commitments)
	for j := 0; j < len(all_commitments); j++ {
		all_commitments_read[j] = ReadFromFile("./tmp/data" + strconv.Itoa(j))
	}

	res := ValidatorCheck(ks, content_x, all_commitments_read, all_evaluation, all_proof)

	if !res {
		t.Fatal("test fail")
	}
}