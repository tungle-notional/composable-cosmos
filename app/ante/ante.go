package ante

import (
	ibcante "github.com/cosmos/ibc-go/v7/modules/core/ante"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"

	"github.com/cosmos/cosmos-sdk/codec"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"

	tfmwKeeper "github.com/notional-labs/centauri/v6/x/transfermiddleware/keeper"
	txBoundaryAnte "github.com/notional-labs/centauri/v6/x/tx-boundary/ante"
	txBoundaryKeeper "github.com/notional-labs/centauri/v6/x/tx-boundary/keeper"

	democracyante "github.com/cosmos/interchain-security/v3/app/consumer-democracy/ante"
	consumerante "github.com/cosmos/interchain-security/v3/app/consumer/ante"
	ccvconsumerkeeper "github.com/cosmos/interchain-security/v3/x/ccv/consumer/keeper"
)

// Link to default ante handler used by cosmos sdk:
// https://github.com/cosmos/cosmos-sdk/blob/v0.43.0/x/auth/ante/ante.go#L41
func NewAnteHandler(
	_ servertypes.AppOptions,
	ak ante.AccountKeeper,
	sigGasConsumer ante.SignatureVerificationGasConsumer,
	signModeHandler signing.SignModeHandler,
	channelKeeper *ibckeeper.Keeper,
	tfmwKeeper tfmwKeeper.Keeper,
	txBoundaryKeeper txBoundaryKeeper.Keeper,
	consumerKeeper ccvconsumerkeeper.Keeper,
	codec codec.BinaryCodec,
) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(), //  // outermost AnteDecorator. SetUpContext must be called first
		consumerante.NewDisabledModulesDecorator("/cosmos.evidence", "/cosmos.slashing"),
		democracyante.NewForbiddenProposalsDecorator(IsProposalWhitelisted, IsModuleWhiteList),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(ak),
		ante.NewConsumeGasForTxSizeDecorator(ak),
		NewIBCPermissionDecorator(codec, tfmwKeeper),
		txBoundaryAnte.NewStakingPermissionDecorator(codec, txBoundaryKeeper),
		ante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(ak),
		ante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
		ante.NewSigVerificationDecorator(ak, signModeHandler),
		ante.NewIncrementSequenceDecorator(ak),
		ibcante.NewRedundantRelayDecorator(channelKeeper),
	)
}
