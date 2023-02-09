package types

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
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
				FindingDesc: []byte(finding.Title),
			}
			encAny, _ = codectypes.NewAnyWithValue(&encKeyMsg)
			finding.FindingDesc = encAny

			testByte := finding.GetFindingDesc().GetFindingDesc()

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
				FindingPoc: []byte(finding.Title),
			}
			encAny, _ = codectypes.NewAnyWithValue(&encKeyMsg)
			finding.FindingPoc = encAny

			testByte := finding.GetFindingPoc().GetFindingPoc()

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

func TestGetEncryptedComment(t *testing.T) {
	testCases := []struct {
		name    string
		args    string
		expPass bool
	}{
		{"empty Finding", "", true},
		{"no empty Finding", "string", true},
	}

	for _, testCase := range testCases {
		eciesEncryptedComment := EciesEncryptedComment{
			FindingComment: []byte(testCase.args),
		}
		encAny, err := codectypes.NewAnyWithValue(&eciesEncryptedComment)
		require.NoError(t, err)

		finding := Finding{}
		finding.FindingComment = encAny

		testByte := finding.GetFindingComment().GetFindingComment()

		if testCase.expPass {
			require.Equal(t, testByte, []byte(testCase.args))
		} else {
			require.NotEqual(t, testByte, []byte(testCase.args))
		}
	}

}
