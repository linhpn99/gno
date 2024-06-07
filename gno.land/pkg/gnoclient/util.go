package gnoclient

import "github.com/gnolang/gno/tm2/pkg/std"

func (cfg BaseTxCfg) validateBaseTxConfig() error {
	if cfg.GasWanted <= 0 {
		return ErrInvalidGasWanted
	}

	if cfg.GasFee == "" {
		return ErrInvalidGasFee
	}

	return nil
}

func (msg MsgCall) validateMsg() error {
	if msg.PkgPath == "" {
		return ErrEmptyPkgPath
	}

	if msg.FuncName == "" {
		return ErrEmptyFuncName
	}

	return nil
}

func (msg MsgCall) getCoins() (std.Coins, error) {
	coins, err := std.ParseCoins(msg.Send)
	if err != nil {
		return nil, err
	}

	return coins, nil
}

func (msg MsgSend) validateMsg() error {
	if msg.ToAddress.IsZero() {
		return ErrInvalidToAddress
	}

	_, err := std.ParseCoins(msg.Send)
	if err != nil {
		return ErrInvalidSendAmount
	}

	return nil
}

func (msg MsgSend) getCoins() (std.Coins, error) {
	coins, err := std.ParseCoins(msg.Send)
	if err != nil {
		return nil, err
	}

	return coins, nil
}

func (msg MsgRun) validateMsg() error {
	if msg.Package == nil || len(msg.Package.Files) == 0 {
		return ErrEmptyPackage
	}

	return nil
}

func (msg MsgRun) getCoins() (std.Coins, error) {
	coins, err := std.ParseCoins(msg.Send)
	if err != nil {
		return nil, err
	}

	return coins, nil
}

func (msg MsgAddPackage) validateMsg() error {
	if msg.Package == nil || len(msg.Package.Files) == 0 {
		return ErrEmptyPackage
	}

	return nil
}

func (msg MsgAddPackage) getCoins() (std.Coins, error) {
	coins, err := std.ParseCoins(msg.Deposit)
	if err != nil {
		return nil, err
	}

	return coins, nil
}
