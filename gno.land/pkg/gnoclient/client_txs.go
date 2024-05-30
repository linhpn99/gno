package gnoclient

import (
	"github.com/gnolang/gno/gno.land/pkg/sdk/vm"
	"github.com/gnolang/gno/tm2/pkg/amino"
	ctypes "github.com/gnolang/gno/tm2/pkg/bft/rpc/core/types"
	"github.com/gnolang/gno/tm2/pkg/crypto"
	"github.com/gnolang/gno/tm2/pkg/errors"
	"github.com/gnolang/gno/tm2/pkg/sdk/bank"
	"github.com/gnolang/gno/tm2/pkg/std"
)

// Define various error messages for different validation failures
var (
	ErrEmptyPackage      = errors.New("empty package to run")
	ErrEmptyPkgPath      = errors.New("empty pkg path")
	ErrEmptyFuncName     = errors.New("empty function name")
	ErrInvalidGasWanted  = errors.New("invalid gas wanted")
	ErrInvalidGasFee     = errors.New("invalid gas fee")
	ErrMissingSigner     = errors.New("missing Signer")
	ErrMissingRPCClient  = errors.New("missing RPCClient")
	ErrInvalidToAddress  = errors.New("invalid send to address")
	ErrInvalidSendAmount = errors.New("invalid send amount")
	ErrInvalidMsgType    = errors.New("invalid msg type")
)

// BaseTxCfg defines the base transaction configuration, shared by all message types
type BaseTxCfg struct {
	GasFee         string // Gas fee
	GasWanted      int64  // Gas wanted
	AccountNumber  uint64 // Account number
	SequenceNumber uint64 // Sequence number
	Memo           string // Memo
}

// MsgCall - syntax sugar for vm.MsgCall
type MsgCall struct {
	PkgPath  string   // Package path
	FuncName string   // Function name
	Args     []string // Function arguments
	Send     string   // Send amount
}

// MsgSend - syntax sugar for bank.MsgSend
type MsgSend struct {
	ToAddress crypto.Address // Send to address
	Send      string         // Send amount
}

// MsgRun - syntax sugar for vm.MsgRun
type MsgRun struct {
	Package *std.MemPackage // Package to run
	Send    string          // Send amount
}

// MsgAddPackage - syntax sugar for vm.MsgAddPackage
type MsgAddPackage struct {
	Package *std.MemPackage // Package to add
	Deposit string          // Coin deposit
}

// MsgNoop - syntax sugar for vm.NoopMsg
type MsgNoop struct {
	Caller crypto.Address
}

type Msg interface {
	validateMsg() error
	getCoins() (std.Coins, error)
}

// Call executes one or more MsgCall calls on the blockchain
func (c *Client) Call(cfg BaseTxCfg, msgs ...MsgCall) (*ctypes.ResultBroadcastTxCommit, error) {
	// Validate required client fields.
	if err := c.validateClient(); err != nil {
		return nil, err
	}

	// Validate base transaction config
	if err := cfg.validateBaseTxConfig(); err != nil {
		return nil, err
	}

	// Parse MsgCall slice
	vmMsgs := make([]std.Msg, 0, len(msgs))
	for _, msg := range msgs {
		// Validate MsgCall fields
		if err := msg.validateMsg(); err != nil {
			return nil, err
		}

		// Parse send coins
		send, err := msg.getCoins()
		if err != nil {
			return nil, err
		}

		caller := c.Signer.Info().GetAddress()

		// Unwrap syntax sugar to vm.MsgCall slice
		vmMsgs = append(vmMsgs, vm.MsgCall{
			Caller:  caller,
			PkgPath: msg.PkgPath,
			Func:    msg.FuncName,
			Args:    msg.Args,
			Send:    send,
		})

	}

	return c.sendTransaction(cfg, vmMsgs...)
}

// Run executes one or more MsgRun calls on the blockchain
func (c *Client) Run(cfg BaseTxCfg, msgs ...MsgRun) (*ctypes.ResultBroadcastTxCommit, error) {
	// Validate required client fields.
	if err := c.validateClient(); err != nil {
		return nil, err
	}

	// Validate base transaction config
	if err := cfg.validateBaseTxConfig(); err != nil {
		return nil, err
	}

	// Parse MsgRun slice
	vmMsgs := make([]std.Msg, 0, len(msgs))
	for _, msg := range msgs {
		// Validate MsgCall fields
		if err := msg.validateMsg(); err != nil {
			return nil, err
		}

		// Parse send coins
		send, err := msg.getCoins()
		if err != nil {
			return nil, err
		}

		caller := c.Signer.Info().GetAddress()

		msg.Package.Name = "main"
		msg.Package.Path = ""

		// Unwrap syntax sugar to vm.MsgCall slice
		vmMsgs = append(vmMsgs, vm.MsgRun{
			Caller:  caller,
			Package: msg.Package,
			Send:    send,
		})

	}

	return c.sendTransaction(cfg, vmMsgs...)
}

// Send executes one or more MsgSend calls on the blockchain
func (c *Client) Send(cfg BaseTxCfg, msgs ...MsgSend) (*ctypes.ResultBroadcastTxCommit, error) {
	// Validate required client fields.
	if err := c.validateClient(); err != nil {
		return nil, err
	}

	// Validate base transaction config
	if err := cfg.validateBaseTxConfig(); err != nil {
		return nil, err
	}

	// Parse MsgSend slice
	vmMsgs := make([]std.Msg, 0, len(msgs))
	for _, msg := range msgs {
		// Validate MsgSend fields
		if err := msg.validateMsg(); err != nil {
			return nil, err
		}

		// Parse send coins
		send, err := std.ParseCoins(msg.Send)
		if err != nil {
			return nil, err
		}

		caller := c.Signer.Info().GetAddress()

		// Unwrap syntax sugar to vm.MsgSend slice
		vmMsgs = append(vmMsgs, bank.MsgSend{
			FromAddress: caller,
			ToAddress:   msg.ToAddress,
			Amount:      send,
		})
	}

	return c.sendTransaction(cfg, vmMsgs...)
}

// AddPackage executes one or more AddPackage calls on the blockchain
func (c *Client) AddPackage(cfg BaseTxCfg, msgs ...MsgAddPackage) (*ctypes.ResultBroadcastTxCommit, error) {
	// Validate required client fields.
	if err := c.validateClient(); err != nil {
		return nil, err
	}

	// Validate base transaction config
	if err := cfg.validateBaseTxConfig(); err != nil {
		return nil, err
	}

	// Parse MsgRun slice
	vmMsgs := make([]std.Msg, 0, len(msgs))
	for _, msg := range msgs {
		// Validate MsgCall fields
		if err := msg.validateMsg(); err != nil {
			return nil, err
		}

		// Parse deposit coins
		deposit, err := msg.getCoins()
		if err != nil {
			return nil, err
		}

		caller := c.Signer.Info().GetAddress()

		// Unwrap syntax sugar to vm.MsgAddPackage slice
		vmMsgs = append(vmMsgs, vm.MsgAddPackage{
			Creator: caller,
			Package: msg.Package,
			Deposit: deposit,
		})

	}

	return c.sendTransaction(cfg, vmMsgs...)
}

// Sponsor allows sending one or more transactions (represented by `msgs`) using the signer's account to pay transaction fees.
// The `sponsoree` account benefits from these transactions without incurring any gas costs.
func (c *Client) Sponsor(cfg BaseTxCfg, sponsoree crypto.Address, msgs ...Msg) (*ctypes.ResultBroadcastTxCommit, error) {
	// Validate required client fields.
	if err := c.validateClient(); err != nil {
		return nil, err
	}

	// Validate base transaction config
	if err := cfg.validateBaseTxConfig(); err != nil {
		return nil, err
	}

	// Parse Msg slice
	vmMsgs := make([]std.Msg, 0, len(msgs)+1)

	// the first msg in list must be MsgNoop
	vmMsgs = append(vmMsgs, vm.MsgNoop{
		Caller: c.Signer.Info().GetAddress(),
	})

	for _, msg := range msgs {
		// Validate MsgCall fields
		if err := msg.validateMsg(); err != nil {
			return nil, err
		}

		// Parse send/deposit coins
		coins, err := msg.getCoins()
		if err != nil {
			return nil, err
		}

		switch m := msg.(type) {
		case MsgNoop:
		// Ignore duplicate MsgNoop (optional)
		// return nil, errors.New("MsgNoop already present")  // uncomment to enforce single MsgNoop

		case MsgCall:
			// Unwrap syntax sugar to vm.MsgCall slice
			vmMsgs = append(vmMsgs, vm.MsgCall{
				Caller:  sponsoree,
				PkgPath: m.PkgPath,
				Func:    m.FuncName,
				Args:    m.Args,
				Send:    coins,
			})

		case MsgSend:
			// Unwrap syntax sugar to vm.MsgSend slice
			vmMsgs = append(vmMsgs, bank.MsgSend{
				FromAddress: sponsoree,
				ToAddress:   m.ToAddress,
				Amount:      coins,
			})

		case MsgRun:
			m.Package.Name = "main"
			m.Package.Path = ""

			// Unwrap syntax sugar to vm.MsgRun slice
			vmMsgs = append(vmMsgs, vm.MsgRun{
				Caller:  sponsoree,
				Package: m.Package,
				Send:    coins,
			})

		case MsgAddPackage:
			// Unwrap syntax sugar to vm.MsgAddPackage slice
			vmMsgs = append(vmMsgs, vm.MsgAddPackage{
				Creator: sponsoree,
				Package: m.Package,
				Deposit: coins,
			})

		default:
			return nil, ErrInvalidMsgType
		}
	}

	return c.sendTransaction(cfg, vmMsgs...)
}

// signAndBroadcastTxCommit signs a transaction and broadcasts it, returning the result
func (c *Client) signAndBroadcastTxCommit(tx std.Tx, accountNumber, sequenceNumber uint64) (*ctypes.ResultBroadcastTxCommit, error) {
	caller := c.Signer.Info().GetAddress()

	if sequenceNumber == 0 || accountNumber == 0 {
		account, _, err := c.QueryAccount(caller)
		if err != nil {
			return nil, errors.Wrap(err, "query account")
		}
		accountNumber = account.AccountNumber
		sequenceNumber = account.Sequence
	}

	signCfg := SignCfg{
		UnsignedTX:     tx,
		SequenceNumber: sequenceNumber,
		AccountNumber:  accountNumber,
	}
	signedTx, err := c.Signer.Sign(signCfg)
	if err != nil {
		return nil, errors.Wrap(err, "sign")
	}

	bz, err := amino.Marshal(signedTx)
	if err != nil {
		return nil, errors.Wrap(err, "marshaling tx binary bytes")
	}

	bres, err := c.RPCClient.BroadcastTxCommit(bz)
	if err != nil {
		return nil, errors.Wrap(err, "broadcasting bytes")
	}

	if bres.CheckTx.IsErr() {
		return bres, errors.Wrap(bres.CheckTx.Error, "check transaction failed: log:%s", bres.CheckTx.Log)
	}
	if bres.DeliverTx.IsErr() {
		return bres, errors.Wrap(bres.DeliverTx.Error, "deliver transaction failed: log:%s", bres.DeliverTx.Log)
	}

	return bres, nil
}

// sendTransaction creates and sends a transaction containing the provided messages.
// It uses the given transaction configuration for gas fee and other parameters.
func (c *Client) sendTransaction(cfg BaseTxCfg, msgs ...std.Msg) (*ctypes.ResultBroadcastTxCommit, error) {
	// Parse gas fee
	gasFeeCoins, err := std.ParseCoin(cfg.GasFee)
	if err != nil {
		return nil, err
	}

	// Pack transaction
	tx := std.Tx{
		Msgs:       msgs,
		Fee:        std.NewFee(cfg.GasWanted, gasFeeCoins),
		Signatures: nil,
		Memo:       cfg.Memo,
	}

	// Sign and broadcast the transaction, then return the result.
	return c.signAndBroadcastTxCommit(tx, cfg.AccountNumber, cfg.SequenceNumber)
}

// TODO: Add more functionality, examples, and unit tests.
