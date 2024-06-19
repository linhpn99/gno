package client

import (
	"flag"
	"fmt"

	"github.com/gnolang/gno/tm2/pkg/amino"
	"github.com/gnolang/gno/tm2/pkg/commands"
	"github.com/gnolang/gno/tm2/pkg/crypto/keys"
	"github.com/gnolang/gno/tm2/pkg/errors"
	"github.com/gnolang/gno/tm2/pkg/std"
)

type MakeTxCfg struct {
	RootCfg *BaseCfg

	GasWanted int64
	GasFee    string
	Memo      string

	Broadcast bool
	// Valid options are SimulateTest, SimulateSkip or SimulateOnly.
	Simulate string
	ChainID  string

	// Optional
	Sponsor string
}

// These are the valid options for MakeTxConfig.Simulate.
const (
	SimulateTest = "test"
	SimulateSkip = "skip"
	SimulateOnly = "only"
)

func (c *MakeTxCfg) Validate() error {
	switch c.Simulate {
	case SimulateTest, SimulateSkip, SimulateOnly:
	default:
		return fmt.Errorf("invalid simulate option: %q", c.Simulate)
	}
	return nil
}

func NewMakeTxCmd(rootCfg *BaseCfg, io commands.IO) *commands.Command {
	cfg := &MakeTxCfg{
		RootCfg: rootCfg,
	}

	cmd := commands.NewCommand(
		commands.Metadata{
			Name:       "maketx",
			ShortUsage: "<subcommand> [flags] [<arg>...]",
			ShortHelp:  "composes a tx document to sign",
		},
		cfg,
		commands.HelpExec,
	)

	cmd.AddSubCommands(
		NewMakeSendCmd(cfg, io),
	)

	return cmd
}

func (c *MakeTxCfg) RegisterFlags(fs *flag.FlagSet) {
	fs.Int64Var(
		&c.GasWanted,
		"gas-wanted",
		0,
		"gas requested for tx",
	)

	fs.StringVar(
		&c.GasFee,
		"gas-fee",
		"",
		"gas payment fee",
	)

	fs.StringVar(
		&c.Memo,
		"memo",
		"",
		"any descriptive text",
	)

	fs.BoolVar(
		&c.Broadcast,
		"broadcast",
		false,
		"simulate and broadcast",
	)

	fs.StringVar(
		&c.Simulate,
		"simulate",
		"test",
		`select how to simulate the transaction (only useful with --broadcast); valid options are
		- test: attempts simulating the transaction, and if successful performs broadcasting (default)
		- skip: avoids performing transaction simulation
		- only: avoids broadcasting transaction (ie. dry run)`,
	)

	fs.StringVar(
		&c.ChainID,
		"chainid",
		"dev",
		"chainid to sign for (only useful with --broadcast)",
	)

	fs.StringVar(
		&c.Sponsor,
		"sponsor",
		"",
		"onchain address of the sponsor",
	)
}

func SignHandler(
	cfg *MakeTxCfg,
	nameOrBech32 string,
	tx *std.Tx,
	pass string,
) error {
	baseopts := cfg.RootCfg
	txopts := cfg

	kb, err := keys.NewKeyBaseFromDir(cfg.RootCfg.Home)
	if err != nil {
		return err
	}

	info, err := kb.GetByNameOrAddress(nameOrBech32)
	if err != nil {
		return err
	}
	accountAddr := info.GetAddress()

	var accountNumber uint64 = 0
	var sequence uint64 = 0

	qopts := &QueryCfg{
		RootCfg: baseopts,
		Path:    fmt.Sprintf("auth/accounts/%s", accountAddr),
	}

	qres, err := QueryHandler(qopts)
	if err != nil {
		if !tx.IsSponsorTx() {
			return errors.Wrap(err, "query account")
		}
	} else {
		var qret struct {
			BaseAccount std.BaseAccount
		}

		err = amino.UnmarshalJSON(qres.Response.Data, &qret)
		if err != nil {
			return err
		}
		accountNumber = qret.BaseAccount.AccountNumber
		sequence = qret.BaseAccount.Sequence
	}

	sOpts := signOpts{
		chainID:         txopts.ChainID,
		accountSequence: sequence,
		accountNumber:   accountNumber,
	}

	// sign tx
	kOpts := keyOpts{
		keyName:     nameOrBech32,
		decryptPass: pass,
	}

	if err := signTx(tx, kb, sOpts, kOpts); err != nil {
		return fmt.Errorf("unable to sign transaction, %w", err)
	}

	return nil
}

func ExecSign(cfg *MakeTxCfg, args []string, tx *std.Tx, io commands.IO) error {
	if err := cfg.Validate(); err != nil {
		return err
	}

	baseopts := cfg.RootCfg

	// query account
	nameOrBech32 := args[0]

	var err error
	var pass string
	if baseopts.Quiet {
		pass, err = io.GetPassword("", baseopts.InsecurePasswordStdin)
	} else {
		pass, err = io.GetPassword("Enter password.", baseopts.InsecurePasswordStdin)
	}

	if err != nil {
		return err
	}

	return SignHandler(cfg, nameOrBech32, tx, pass)
}

func ExecBroadcast(cfg *MakeTxCfg, tx *std.Tx, io commands.IO) error {
	// broadcast signed tx
	bopts := &BroadcastCfg{
		RootCfg: cfg.RootCfg,
		tx:      tx,

		DryRun:       cfg.Simulate == SimulateOnly,
		testSimulate: cfg.Simulate == SimulateTest,
	}

	bres, err := BroadcastHandler(bopts)
	if err != nil {
		return errors.Wrap(err, "broadcast tx")
	}
	if bres.CheckTx.IsErr() {
		return errors.Wrap(bres.CheckTx.Error, "check transaction failed: log:%s", bres.CheckTx.Log)
	}
	if bres.DeliverTx.IsErr() {
		return errors.Wrap(bres.DeliverTx.Error, "deliver transaction failed: log:%s", bres.DeliverTx.Log)
	}

	io.Println(string(bres.DeliverTx.Data))
	io.Println("OK!")
	io.Println("GAS WANTED:", bres.DeliverTx.GasWanted)
	io.Println("GAS USED:  ", bres.DeliverTx.GasUsed)
	io.Println("HEIGHT:    ", bres.Height)
	io.Println("EVENTS:    ", string(bres.DeliverTx.EncodeEvents()))

	return nil
}

func ExecSignAndBroadcast(
	cfg *MakeTxCfg,
	args []string,
	tx *std.Tx,
	io commands.IO,
) error {
	err := ExecSign(cfg, args, tx, io)
	if err != nil {
		return err
	}

	return ExecBroadcast(cfg, tx, io)
}
