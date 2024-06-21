package kzg

import (
	"github.com/protolambda/go-kzg/bls"
	// "testing"
	"os"
	"io"
	"io/ioutil"
)

var num_of_commitments uint64
var epsilon_power uint64

func initialzation(beta uint64, epsilon uint8) *KZGSettings {
	epsilon_power = uint64(1) << uint64(epsilon)
	num_of_commitments = beta/(1 << uint64(epsilon))/256
    var fs *FFTSettings = NewFFTSettings(epsilon)
	var (
		s1 []bls.G1Point
		s2 []bls.G2Point
	)
    s1, s2 = GenerateTestingSetup("1927409816240961209460912649124", epsilon_power)
	var ks *KZGSettings = NewKZGSettings(fs, s1, s2)
	return ks
}

func ProposerCommit(ks *KZGSettings) ([]string, [][]bls.Fr) {
	// first generate num_of_commitments polynomials and get a total of beta bits
	all_commitments := make([]string, num_of_commitments)
	all_poly := make([][]bls.Fr, num_of_commitments)

	for j := 0; j < int(num_of_commitments); j++ {
		// generate random bytes
		all_poly[j] = make([]bls.Fr, epsilon_power)
		for i := 0; i < len(all_poly[j]); i++ {
			all_poly[j][i] = *bls.RandomFr()
			all_commitments[i] = bls.StrG1(ks.CommitToPoly(all_poly[j]))
		}
	}
	return all_commitments, all_poly
}

func ValidatorSampling() string {
	// sample a random point and return its string form
	point_to_evaluate_at := *bls.RandomFr()
	content_x := bls.FrStr(&point_to_evaluate_at)
	return content_x
}

func ProposerReply(ks *KZGSettings, all_poly [][]bls.Fr, content_x string, all_commitments []string) ([]string, []string) {
	all_evaluation := make([]string, len(all_commitments))
	all_proof := make([]string, len(all_commitments))

	var point_to_evaluate_at bls.Fr
	bls.SetFr(&point_to_evaluate_at, content_x)

	for j := 0; j < len(all_commitments); j++ {
		var value bls.Fr
		bls.EvalPolyAt(&value, all_poly[j], &point_to_evaluate_at)
		all_evaluation[j] = bls.FrStr(&value)

		proof := ks.ComputeProofSingleFr(all_poly[j], point_to_evaluate_at)
		all_proof[j] = bls.StrG1(proof)
	}
	return all_evaluation, all_proof
}

func ValidatorCheck(ks *KZGSettings, content_x string, all_commitments []string, all_evaluation []string, all_proof []string) bool{
	var point_to_evaluate_at bls.Fr
	bls.SetFr(&point_to_evaluate_at, content_x)

	for j := 0; j < len(all_commitments); j++ {
		var evaluation_j bls.Fr
		bls.SetFr(&evaluation_j, all_evaluation[j])

		var commitment bls.G1Point
		bls.SetG1(&commitment, all_commitments[j])

		var proof bls.G1Point
		bls.SetG1(&proof, all_commitments[j])
		
		if !ks.CheckProofSingle(&commitment, &proof, &point_to_evaluate_at, &evaluation_j){
			return false
		}
	}

	return true
}

func SaveToFile(toSave string, filename ...string){
	thefilename := ""
	if len(filename) == 0 {
		thefilename = "./tmp/data"
	} else {
		thefilename = filename[0]
	}

	f, _ := os.Create(thefilename)
	io.WriteString(f, toSave)
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
	return string(toReturn[:])
}