package types

import (
	"context"
	"time"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/log"

	clientcontext "github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// Context define the type of data for service node.
type Context struct {
	*clientcontext.CLIContext
	ctx       context.Context
	txBuilder *authtypes.TxBuilder
	config    *Config
	logger    log.Logger
	service   Service
}

// Context returns the internal context object.
func (c Context) Context() context.Context { return c.ctx }

// ClientContext returns a copy of the certik chain client context object value.
func (c Context) ClientContext() clientcontext.CLIContext { return *c.CLIContext }

// Codec returns a reference to the client codec.
func (c Context) Codec() *amino.Codec { return c.CLIContext.Codec }

// TxBuilder returns a copy of the certik chain transaction builder object.
func (c Context) TxBuilder() authtypes.TxBuilder { return *c.txBuilder }

// Config returns a copy of the service node global configuration.
func (c Context) Config() Config { return *c.config }

// Logger returns the logger for service node internal use.
func (c Context) Logger() log.Logger { return c.logger }

// Service returns the service registered on certik chain and run by the service node.
func (c Context) Service() Service { return c.service }

// NewContextWithDefaultConfigAndLogger returns a new context with global configuration set from a config file.
func NewContextWithDefaultConfigAndLogger() (Context, error) {
	if err := initConfig(); err != nil {
		return Context{}, err
	}
	if err := initLogger(); err != nil {
		return Context{}, err
	}
	return NewContext(&config, logger), nil
}

// NewContext creates a new context.
func NewContext(config *Config, logger log.Logger) Context {
	return Context{
		CLIContext: &clientcontext.CLIContext{},
		ctx:        context.Background(),
		txBuilder:  &authtypes.TxBuilder{},
		config:     config,
		logger:     logger,
	}
}

// WithContext returns a copy of the context with an updated internal context.
func (c Context) WithContext(ctx context.Context) Context {
	c.ctx = ctx
	return c
}

// WithClientContext returns a copy of the context with an updated CosmoSDK client context.
func (c Context) WithClientContext(ctx *clientcontext.CLIContext) Context {
	c.CLIContext = ctx
	return c
}

// WithTxBuilder returns a copy of the context with an updated tx builder.
func (c Context) WithTxBuilder(txBuilder *authtypes.TxBuilder) Context {
	c.txBuilder = txBuilder
	return c
}

// WithConfig returns a copy of the context with an updated configuration setting.
func (c Context) WithConfig(config *Config) Context {
	c.config = config
	return c
}

// WithLogger returns a copy of the context with an updated logger.
func (c Context) WithLogger(logger log.Logger) Context {
	c.logger = logger
	return c
}

// WithLoggerLabels returns a copy of the context with updated logger labels.
func (c Context) WithLoggerLabels(keyvals ...interface{}) Context {
	c.logger = c.logger.With(keyvals...)
	return c
}

// // WithService returns a copy of the context with an updated service info.
func (c Context) WithService(s Service) Context {
	c.service = s
	return c
}

// WithValue returns a copy of the context with an extra key-value information.
func (c Context) WithValue(key, value interface{}) Context {
	c.ctx = context.WithValue(c.ctx, key, value)
	return c
}

type ServiceID string

// Service defines the data structure of a service.
type Service struct {
	ID              ServiceID      `json:"service_id"`
	Name            string         `json:"name"`
	Protocol        string         `json:"protocol"`
	ContractAddress string         `json:"contract_address"`
	RequestDataABI  string         `json:"request_data_abi"`
	ResponseDataABI string         `json:"response_data_abi"`
	Owner           sdk.AccAddress `json:"service_owner"`
	ProviderID      ServiceID      `json:"service_provider_id"`
	Tags            []string       `json:"tags"`
	Description     string         `json:"descriptions"`
	CreateTime      time.Time      `json:"create_time"`
	RevokeTime      time.Time      `json:"revoke_time"`
}
