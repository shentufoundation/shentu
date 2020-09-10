package runner

import (
	"fmt"
	"net/url"
	"strconv"
	"sync"

	"github.com/tendermint/tendermint/libs/kv"
	ctkClient "github.com/tendermint/tendermint/rpc/client/http"
	tendermintTypes "github.com/tendermint/tendermint/types"

	"github.com/certikfoundation/shentu/toolsets/oracle-operator/querier"
	runnerTypes "github.com/certikfoundation/shentu/toolsets/oracle-operator/runner/types"
	oracleTypes "github.com/certikfoundation/shentu/toolsets/oracle-operator/types"
)

// Listen listens for events from CertiK chain.
func Listen(ctx oracleTypes.Context, ctkMsgChan chan<- interface{}, errorChan chan<- error) {
	// load configuration and logger
	logger := ctx.Logger()
	node := ctx.ClientContext().NodeURI
	logger.Info("start to listen to certik-chain", "node", node)

	// initialize client
	client, err := ctkClient.New(ctx.ClientContext().NodeURI, "/websocket")
	if err != nil {
		logger.Error("ctkClient dialing", "error", err.Error())
		errorChan <- err
		return
	}

	// start the listener
	err = client.Start()
	if err != nil {
		logger.Error("ctkClient subscribing", "error", err.Error())
		errorChan <- err
		return
	}
	defer client.Stop()

	// subscribe the TXs according to the query
	query := "tm.event='Tx'"
	txChan, err := client.Subscribe(ctx.Context(), "", query, 1000)
	if err != nil {
		logger.Error("ctkClient subscribing", "error", err.Error())
		errorChan <- err
		return
	}

	for {
		select {
		case <-ctx.Context().Done():
			logger.Info("stop listening...")
			return
		case tx := <-txChan:
			// get tendermint transaction data in struct of ResponseDeliverTx
			txData, ok := tx.Data.(tendermintTypes.EventDataTx)
			if !ok {
				logger.Error("ctkClient monitored messaging", "error")
				return
			}
			go listenHandler(ctx, txData, ctkMsgChan, errorChan)
		}
	}
}

// listenHandler deals with events inside tx data result.
func listenHandler(ctx oracleTypes.Context, txData tendermintTypes.EventDataTx, ctkMsgChan chan<- interface{}, errorChan chan<- error) {
	logger := ctx.Logger()
	for _, event := range txData.Result.Events {
		switch event.Type {
		case "create_task":
			logger.Info("Received event", "eventType", "create_task")
			go handleMsgCreateTask(ctx.WithLoggerLabels("type", "CreateTask"), event.Attributes, ctkMsgChan, errorChan)
		}
	}
}

// getMsgCreateTask parses TX data and pass organized message to endpoint querier.
func handleMsgCreateTask(ctx oracleTypes.Context, kvs []kv.Pair, ctkMsgChan chan<- interface{}, errorChan chan<- error) {
	logger := ctx.Logger()

	logger.Debug("handling create_task event from certik-chain...")
	msgCreateTaskRequest, err := getMsgCreateTask(ctx, kvs)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	var wg sync.WaitGroup
	msgPrimitiveRespChan := make(chan runnerTypes.MsgPrimitiveResponse, 1000)
	wg.Add(len(ctx.Config().Runner.Combination.Primitives))
	for _, primitive := range ctx.Config().Runner.Combination.Primitives {
		go handlPrimitive(ctx.WithLoggerLabels("primitive", primitive.PrimitiveContractAddr),
			msgCreateTaskRequest, primitive, msgPrimitiveRespChan, &wg)
	}

	wg.Wait()
	// close msgPrimtiveRespChan when all handlePrimitive goroutine finished
	close(msgPrimitiveRespChan)
	switch ctx.Config().Runner.Combination.Type {
	case "linear":
		go handlStrategyLinear(ctx.WithLoggerLabels("strategy", "linear"), msgCreateTaskRequest,
			msgPrimitiveRespChan, ctkMsgChan)
	default:
		logger.Error("Combination type not found", "strategyType", ctx.Config().Runner.Combination.Type)
		errorChan <- fmt.Errorf("combination type %v not found", ctx.Config().Runner.Combination.Type)
	}
}

// getMsgCreateTask parses TX data of creating tasks.
func getMsgCreateTask(ctx oracleTypes.Context, kvs []kv.Pair) (oracleTypes.MsgCreateTask, error) {
	// load logger
	logger := ctx.Logger()
	logger.Debug("parsing create_task event from certik-chain...")

	var contract, function string
	for _, v := range kvs {
		switch string(v.GetKey()) {
		case "contract":
			contract = string(v.GetValue())
		case "function":
			function = string(v.GetValue())
		default:
			logger.Info("kv pair in event msg", v.GetKey(), v.GetValue())
		}
	}

	if contract == "" || function == "" {
		return oracleTypes.MsgCreateTask{}, fmt.Errorf("missing required field in event content")
	}

	msgCreateTask := oracleTypes.MsgCreateTask{
		Contract: contract,
		Function: function,
	}

	logger.Debug("parsed create_task event from certik-chain",
		"Contract", msgCreateTask.Contract,
		"Function", msgCreateTask.Function)

	return msgCreateTask, nil
}

// handlePrimitive gets score for each primitive.
func handlPrimitive(ctx oracleTypes.Context, msgCreateTask oracleTypes.MsgCreateTask, primitive runnerTypes.Primitive,
	msgPrimitiveRespChan chan<- runnerTypes.MsgPrimitiveResponse, wg *sync.WaitGroup) {
	logger := ctx.Logger()
	logger.Info("handle primitive")
	args := []string{msgCreateTask.Contract, msgCreateTask.Function}
	retBool, retString, err := BuildQueryInsight(ctx.ClientContext(), primitive.PrimitiveContractAddr, ContractFnName, args)
	if err != nil {
		logger.Error(err.Error())
		wg.Done()
		return
	}
	var score uint8
	if retBool {
		_, err := url.ParseRequestURI(retString)
		if err != nil {
			logger.Error(err.Error())
			wg.Done()
			return
		}

		score, err = querier.HandleRequest(ctx.WithLoggerLabels("submodule", "querier"), msgCreateTask, retString)
		logger.Debug("Got score from endpoint url", "url", retString, "score", score)
		if err != nil {
			logger.Error(err.Error())
			wg.Done()
			return
		}
	} else {
		retScore, err := strconv.ParseUint(retString, 10, 8)
		if err != nil {
			logger.Error(err.Error())
			wg.Done()
			return
		}
		score = uint8(retScore)
		logger.Debug("Got score from Security Primitive on certik chain", "score", score)
	}

	msgPrimitiveRespChan <- runnerTypes.MsgPrimitiveResponse{
		Score:     score,
		Primitive: primitive,
	}
	wg.Done()
}

// handleStrategyLinear handles linear combination strategy.
func handlStrategyLinear(ctx oracleTypes.Context, msgCreateTask oracleTypes.MsgCreateTask,
	msgPrimitiveRespChan <-chan runnerTypes.MsgPrimitiveResponse, ctkMsgChan chan<- interface{}) {
	logger := ctx.Logger()
	logger.Info("handle linear strategy")
	sum := float32(0)
	weightSum := float32(0)
	if len(msgPrimitiveRespChan) == 0 {
		logger.Error("none of primitives returns security score.")
		return
	}

	for resp := range msgPrimitiveRespChan {
		sum += float32(resp.Score) * resp.Primitive.Weight
		weightSum += resp.Primitive.Weight
	}
	score := uint8(sum / weightSum)

	logger.Debug("Got Security Score", "Score", score,
		"Contract", msgCreateTask.Contract,
		"Function", msgCreateTask.Function)
	ctkMsgChan <- oracleTypes.MsgRespondToTask{
		Contract: msgCreateTask.Contract,
		Function: msgCreateTask.Function,
		Score:    score,
		Operator: ctx.CLIContext.FromAddress,
	}
}
