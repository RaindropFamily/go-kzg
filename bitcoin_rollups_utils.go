package kzg

import (
	"github.com/protolambda/go-kzg/bls"
	"testing"
	"os"
	"io"
	"io/ioutil"
)

beta := 1 << 33 // a blob is 2^33 bits
epsilon := uint8(12) // the degree of the function if 2^12, which requires 2^12 + 1 points to recover
epsilon_power := uint64(1) << uint64(epsilon) // the degree of the function if 2^12, which requires 2^12 + 1 points to recover
num_of_commitments := beta/(1 << uint64(epsilon))/256

fs := NewFFTSettings(epsilon)
s1, s2 := GenerateTestingSetup("1927409816240961209460912649124", epsilon_power)
ks := NewKZGSettings(fs, s1, s2)

func ProposerCommit() []string {
	// first generate num_of_commitments polynomials and get a total of beta bits
	all_commitments := make([]string, num_of_commitments)

	for j := 0; j < num_of_commitments; j++ {
		// generate random bytes
		poly_to_commit := make([]bls.Fr, epsilon_power)
		for i := 0; i < len(poly_to_commit); i++ {
			poly_to_commit[i] = *bls.RandomFr()
		}
		all_commitments[i] = bls.StrG1(ks.CommitToPoly(poly_to_commit))
	}
	return all_commitments
}

func ValidatorSampling() string {
	// sample a random point and return its string form
	point_to_evaluate_at := *bls.RandomFr()
	content_x := bls.FrStr(&point_to_evaluate_at)
	return content_x
}

func SaveToFile(toSave ...string, filename ...string) string{
	thefilename := ""
	if len(filename) == 0 {
		thefilename = "./tmp/data"
	} else {
		thefilename = filename[0]
	}

	f, _ := os.Create(thefilename)
	for i := 0; i < len(toSave); i++ {
		io.WriteString(f, toSave[i]+"\n")
	}
	f.Close()
}

func ReadFromFile(filename ...string) string{
	thefilename := ""
	if len(filename) == 0 {
		thefilename = "./tmp/data"
	} else {
		thefilename = filename[0]
	}

	toReturn, _:= ioutil.ReadFile(thefilename)
	return toReturn
}