//go:build !bignum_pure && !bignum_hol256
// +build !bignum_pure,!bignum_hol256

package kzg

import (
	"github.com/protolambda/go-kzg/bls"
	"testing"
	"os"
	"io"
	"io/ioutil"
)

func TestKZGforBitcoinRollups(t *testing.T){
	// testing_scale := 12
	fs := NewFFTSettings(12)
	s1, s2 := GenerateTestingSetup("1927409816240961209460912649124", 1 << 12)
	ks := NewKZGSettings(fs, s1, s2)
	// for i := 0; i < len(ks.SecretG1); i++ {
	// 	t.Logf("secret g1 %d: %s", i, bls.StrG1(&ks.SecretG1[i]))
	// }

	poly_to_commit := make([]bls.Fr, uint64(1)<<12)
	for i := 0; i < len(poly_to_commit); i++ {
		poly_to_commit[i] = *bls.RandomFr()
	}

	commitment := ks.CommitToPoly(poly_to_commit)
	
	point_to_evaluate_at := *bls.RandomFr()
	var value bls.Fr
	bls.EvalPolyAt(&value, poly_to_commit, &point_to_evaluate_at)

	proof := ks.ComputeProofSingleFr(poly_to_commit, point_to_evaluate_at)

	if !ks.CheckProofSingle(commitment, proof, &point_to_evaluate_at, &value) {
		t.Fatal("could not verify proof")
	}

	commitment_string := bls.StrG1(commitment)
	proof_string := bls.StrG1(proof)
	var commitment_recover bls.G1Point
	var proof_recover bls.G1Point
	bls.SetG1(&commitment_recover, commitment_string)
	bls.SetG1(&proof_recover, proof_string)
	if !ks.CheckProofSingle(&commitment_recover, &proof_recover, &point_to_evaluate_at, &value) {
		t.Fatal("could not verify proof2")
	}

	content_x := bls.FrStr(&point_to_evaluate_at)
	content_y := bls.FrStr(&value)

	var x_recover bls.Fr
	var y_recover bls.Fr
	bls.SetFr(&x_recover, content_x)
	bls.SetFr(&y_recover, content_y)

	if x_recover != point_to_evaluate_at {
		t.Fatal("x not match")
	}
	if y_recover != value {
		t.Fatal("y not match")
	}

	f_x, _ := os.Create("./tmp/data1_x")
	f_y, _ := os.Create("./tmp/data1_y")
	io.WriteString(f_x, content_x) 
	io.WriteString(f_y, content_y) 
	f_x.Close()
	f_y.Close()

	read_x, _:= ioutil.ReadFile("./tmp/data1_x")
	read_y, _:= ioutil.ReadFile("./tmp/data1_y")
	str_x := string(read_x)
	str_y := string(read_y)
	if str_x != content_x {
		t.Fatal("x read not match")
	}
	if str_y != content_y {
		t.Fatal("y read not match")
	}
}

func TestKZGSettings_CommitToEvalPoly(t *testing.T) {
	fs := NewFFTSettings(4)
	s1, s2 := GenerateTestingSetup("1927409816240961209460912649124", 16+1)
	ks := NewKZGSettings(fs, s1, s2)
	polynomial := testPoly(1, 2, 3, 4, 7, 7, 7, 7, 13, 13, 13, 13, 13, 13, 13, 13)
	evalPoly, err := fs.FFT(polynomial, false)
	if err != nil {
		t.Fatal(err)
	}
	secretG1IFFT, err := fs.FFTG1(ks.SecretG1[:16], true)
	if err != nil {
		t.Fatal(err)
	}

	commitmentByCoeffs := ks.CommitToPoly(polynomial)
	commitmentByEval := CommitToEvalPoly(secretG1IFFT, evalPoly)
	if !bls.EqualG1(commitmentByEval, commitmentByCoeffs) {
		t.Fatalf("expected commitments to be equal, but got:\nby eval: %s\nby coeffs: %s",
			commitmentByEval, commitmentByCoeffs)
	}
}

func TestKZGSettings_CheckProofSingle(t *testing.T) {
	fs := NewFFTSettings(4)
	s1, s2 := GenerateTestingSetup("1927409816240961209460912649124", 16+1)
	ks := NewKZGSettings(fs, s1, s2)
	for i := 0; i < len(ks.SecretG1); i++ {
		t.Logf("secret g1 %d: %s", i, bls.StrG1(&ks.SecretG1[i]))
	}

	polynomial := testPoly(1, 2, 3, 4, 7, 7, 7, 7, 13, 13, 13, 13, 13, 13, 13, 13)
	for i := 0; i < len(polynomial); i++ {
		t.Logf("poly coeff %d: %s", i, bls.FrStr(&polynomial[i]))
	}

	commitment := ks.CommitToPoly(polynomial)
	t.Log("commitment\n", bls.StrG1(commitment))

	proof := ks.ComputeProofSingle(polynomial, 17)
	t.Log("proof\n", bls.StrG1(proof))

	var x bls.Fr
	bls.AsFr(&x, 17)
	var value bls.Fr
	bls.EvalPolyAt(&value, polynomial, &x)
	t.Log("value\n", bls.FrStr(&value))

	if !ks.CheckProofSingle(commitment, proof, &x, &value) {
		t.Fatal("could not verify proof")
	}
}

func testPoly(polynomial ...uint64) []bls.Fr {
	n := len(polynomial)
	polynomialFr := make([]bls.Fr, n, n)
	for i := 0; i < n; i++ {
		bls.AsFr(&polynomialFr[i], polynomial[i])
	}
	return polynomialFr
}
