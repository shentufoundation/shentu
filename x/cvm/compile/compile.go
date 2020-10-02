package compile

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hyperledger/burrow/deploy/compile"
	"github.com/hyperledger/burrow/logging"
)

const NoABI = "\000"

func BytecodeEVM(basename, workDir, abiFile string, logger *logging.Logger) (*compile.Response, error) {
	bytecode, err := ioutil.ReadFile(workDir + "/" + basename)
	if err != nil {
		return nil, err
	}
	if basenameSplit := strings.Split(basename, "."); basenameSplit[len(basenameSplit)-1] == "wasm" {
		bytecode = []byte(hex.EncodeToString(bytecode))
		logger.TraceMsg("Command Output", "ewasm", bytecode)
	} else {
		logger.TraceMsg("Command Output", "bytecode", bytecode)
	}

	abi := json.RawMessage(NoABI)
	if abiFile != "" {
		abiBasename, abiWorkDir, err := ResolveFilename(abiFile)
		if err != nil {
			return nil, err
		}

		abiBasenameSplit := strings.Split(abiBasename, ".")
		if abiBasenameSplit[len(abiBasenameSplit)-1] != "abi" {
			return nil, errors.New("ABI file extension must be .abi")
		}

		abi, err = ioutil.ReadFile(abiWorkDir + "/" + abiBasename)
		if err != nil {
			return nil, err
		}
		logger.TraceMsg("Command Output", "abi", abi)
	}

	return newResponse(basename, bytecode, abi)
}

func DeepseaEVM(basename, workDir string, logger *logging.Logger) (*compile.Response, error) {
	bytecode, err := runDeepseaCompile(basename, workDir, "bytecode")
	if err != nil {
		return nil, err
	}
	logger.TraceMsg("Command Output", "bytecode", bytecode)

	abi, err := runDeepseaCompile(basename, workDir, "abi")
	if err != nil {
		return nil, err
	}
	logger.TraceMsg("Command Output", "abi", abi)

	return newResponse(basename, bytecode, abi)
}

func runDeepseaCompile(basename, workDir, callParam string) ([]byte, error) {
	shellCmd := exec.Command("dsc", basename, callParam)
	if workDir != "" {
		shellCmd.Dir = workDir
	}
	output, err := shellCmd.CombinedOutput()
	return output, err
}

func newResponse(basename string, bytecode, abi json.RawMessage) (*compile.Response, error) {
	contract := compile.SolidityContract{
		Abi: abi,
		Evm: struct {
			Bytecode         compile.ContractCode
			DeployedBytecode compile.ContractCode
		}{Bytecode: compile.ContractCode{
			Object: strings.TrimSpace(string(bytecode)),
		}},
	}

	metamap := make([]compile.MetadataMap, 1)
	metamap = append(metamap, compile.MetadataMap{
		Metadata: compile.Metadata{
			Abi: abi,
		},
	})

	contract.MetadataMap = metamap

	respItem := compile.ResponseItem{
		Filename:   basename,
		Objectname: strings.Split(basename, ".")[0],
		Contract:   contract,
	}

	respItemArray := []compile.ResponseItem{respItem}

	resp := compile.Response{
		Objects: respItemArray,
	}

	return &resp, nil
}

func CheckFileExists(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return err
	} else {
		return nil
	}
}

func ResolveFilename(filename string) (string, string, error) {
	basename := filepath.Base(filename)
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return "", "", err
	}

	err = CheckFileExists(absPath)
	if err != nil {
		return "", "", err
	}

	workDir := filepath.Dir(absPath)

	return basename, workDir, nil
}
