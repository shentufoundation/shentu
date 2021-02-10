package vm

import (
	"fmt"

	"github.com/hyperledger/burrow/acm"
	"github.com/hyperledger/burrow/acm/acmstate"
	"github.com/hyperledger/burrow/execution/engine"
	"github.com/hyperledger/burrow/execution/errors"
	"github.com/hyperledger/burrow/execution/evm"
	"github.com/hyperledger/burrow/execution/exec"
	"github.com/hyperledger/burrow/execution/native"
	"github.com/hyperledger/burrow/logging"
)

type CVM struct {
	options  engine.Options
	sequence uint64
	// Provide any foreign dispatchers to allow calls between VMs
	engine.Externals
	externalDispatcher engine.Dispatcher
	// User dispatcher.CallableProvider to get access to other VMs
	logger *logging.Logger

	// After execution, refund counter is set to the memory refund value
	refund uint64
}

func NewCVM(options engine.Options) *CVM {
	// Set defaults
	if options.MemoryProvider == nil {
		options.MemoryProvider = wrappedDDMP
	}
	if options.Logger == nil {
		options.Logger = logging.NewNoopLogger()
	}
	if options.Natives == nil {
		options.Natives = native.MustDefaultNatives()
	}
	vm := &CVM{
		options: options,
	}
	vm.logger = options.Logger.WithScope("NewVM").With("evm_nonce", options.Nonce)
	vm.externalDispatcher = engine.Dispatchers{&vm.Externals, options.Natives, vm}
	return vm
}

// wrappedDDMP is a wrapper for DefaultDynamicMemoryProvider, to wrap the new MemoryProvider into a gasMemory
func wrappedDDMP(err errors.Sink) engine.Memory {
	memory := gasMemory{
		Memory:      engine.NewDynamicMemory(0, 0x1000000, err),
		lastGasCost: 0,
		refund:      0,
	}
	return memory
}

// Initiate an EVM call against the provided state pushing events to eventSink. code should contain the EVM bytecode,
// input the CallData (readable by CALLDATALOAD), value the amount of native token to transfer with the call
// an quantity metering the number of computational steps available to the execution according to the gas schedule.
func (vm *CVM) Execute(st acmstate.ReaderWriter, blockchain engine.Blockchain, eventSink exec.EventSink,
	params engine.CallParams, code []byte) ([]byte, error) {

	// Make it appear as if natives are stored in state
	st = native.NewState(vm.options.Natives, st)

	state := engine.State{
		CallFrame:  engine.NewCallFrame(st).WithMaxCallStackDepth(vm.options.CallStackMaxDepth),
		Blockchain: blockchain,
		EventSink:  eventSink,
	}

	output, err := vm.Contract(code).Call(state, params)
	if err == nil {
		// Only sync back when there was no exception
		err = state.CallFrame.Sync()
	}
	// Always return output - we may have a reverted exception for which the return is meaningful
	return output, err
}

// SetNonce setss a new nonce and resets the sequence number. Nonces should only be used once!
// A global counter or sufficient randomness will work.
func (vm *CVM) SetNonce(nonce []byte) {
	vm.options.Nonce = nonce
	vm.sequence = 0
}

// SetLogger sets the logger for the CVM instance.
func (vm *CVM) SetLogger(logger *logging.Logger) {
	vm.logger = logger
}

// Dispatch dispatches an account to be used externally from another engine.
func (vm *CVM) Dispatch(acc *acm.Account) engine.Callable {
	// Let the EVM handle code-less (e.g. those created by a call) contracts (so only return nil if there is _other_ non-EVM code)
	if len(acc.EVMCode) == 0 && len(acc.Code()) != 0 {
		return nil
	}
	return vm.Contract(acc.EVMCode)
}

// Contract returns a CVMContract with the provided CVM and code.
func (vm *CVM) Contract(code []byte) *CVMContract {
	return &CVMContract{
		CVM:  vm,
		Code: evm.NewCode(code),
	}
}

func (vm *CVM) debugf(format string, a ...interface{}) {
	if vm.options.DebugOpcodes {
		fmt.Printf(format, a...)
	}
}
