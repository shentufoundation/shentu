package types

import (
	"bytes"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"testing"
)

func TestGetEncryptionKey(t *testing.T) {
	type args struct {
		program []Program
	}

	type errArgs struct {
		shouldPass bool
	}

	tests := []struct {
		name    string
		args    args
		errArgs errArgs
	}{
		{"Program(1)  -> success",
			args{
				program: []Program{
					{
						Description: "for test1",
					},
				},
			},
			errArgs{
				shouldPass: true,
			},
		},
		{"Program(2)  -> empty",
			args{
				program: []Program{
					{
						Description: "",
					},
				},
			},
			errArgs{
				shouldPass: true,
			},
		},
	}

	var encAny *codectypes.Any
	for _, tc := range tests {
		for _, program := range tc.args.program {
			encKeyMsg := EciesPubKey{
				EncryptionKey: []byte(program.Description),
			}
			encAny, _ = codectypes.NewAnyWithValue(&encKeyMsg)
			program.EncryptionKey = encAny

			pubByte := program.GetEncryptionKey().GetEncryptionKey()

			if tc.errArgs.shouldPass {
				if !bytes.Equal(pubByte, []byte(program.Description)) {
					t.Fatal("not equal")
				}
			} else {
				if bytes.Equal(pubByte, []byte(program.Description)) {
					t.Fatal("error equal")
				}
			}

		}
	}
}

func TestGetEncryptedDesc(t *testing.T) {
	type args struct {
		finding []Finding
	}

	type errArgs struct {
		shouldPass bool
	}

	tests := []struct {
		name    string
		args    args
		errArgs errArgs
	}{
		{"Finding(1)  -> success",
			args{
				finding: []Finding{
					{
						Title: "for test1",
					},
				},
			},
			errArgs{
				shouldPass: true,
			},
		},
		{"Finding(2)  -> empty",
			args{
				finding: []Finding{
					{},
				},
			},
			errArgs{
				shouldPass: true,
			},
		},
	}

	var encAny *codectypes.Any
	for _, tc := range tests {
		for _, finding := range tc.args.finding {
			encKeyMsg := EciesEncryptedDesc{
				EncryptedDesc: []byte(finding.Title),
			}
			encAny, _ = codectypes.NewAnyWithValue(&encKeyMsg)
			finding.EncryptedDesc = encAny

			testByte := finding.GetEncryptedDesc().GetEncryptedDesc()

			if tc.errArgs.shouldPass {
				if !bytes.Equal(testByte, []byte(finding.Title)) {
					t.Fatal("not equal")
				}
			} else {
				if bytes.Equal(testByte, []byte(finding.Title)) {
					t.Fatal("error equal")
				}
			}

		}
	}
}
func TestGetEncryptedPoc(t *testing.T) {
	type args struct {
		finding []Finding
	}

	type errArgs struct {
		shouldPass bool
	}

	tests := []struct {
		name    string
		args    args
		errArgs errArgs
	}{
		{"Finding(1)  -> success",
			args{
				finding: []Finding{
					{
						Title: "for test1",
					},
				},
			},
			errArgs{
				shouldPass: true,
			},
		},
		{"Finding(2)  -> empty",
			args{
				finding: []Finding{
					{},
				},
			},
			errArgs{
				shouldPass: true,
			},
		},
	}

	var encAny *codectypes.Any
	for _, tc := range tests {
		for _, finding := range tc.args.finding {
			encKeyMsg := EciesEncryptedPoc{
				EncryptedPoc: []byte(finding.Title),
			}
			encAny, _ = codectypes.NewAnyWithValue(&encKeyMsg)
			finding.EncryptedPoc = encAny

			testByte := finding.GetEncryptedPoc().GetEncryptedPoc()

			if tc.errArgs.shouldPass {
				if !bytes.Equal(testByte, []byte(finding.Title)) {
					t.Fatal("not equal")
				}
			} else {
				if bytes.Equal(testByte, []byte(finding.Title)) {
					t.Fatal("error equal")
				}
			}

		}
	}
}
