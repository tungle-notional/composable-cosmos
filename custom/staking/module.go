package bank

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	stakingmodule "github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"

	// custombankkeeper "github.com/notional-labs/composable/v6/custom/bank/keeper"
	customstakingkeeper "github.com/notional-labs/composable/v6/custom/staking/keeper"
)

// AppModule wraps around the bank module and the bank keeper to return the right total supply
type AppModule struct {
	stakingmodule.AppModule
	keeper   customstakingkeeper.Keeper
	subspace exported.Subspace
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper customstakingkeeper.Keeper, accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper, ss exported.Subspace) AppModule {
	stakingModule := stakingmodule.NewAppModule(cdc, &keeper.Keeper, accountKeeper, bankKeeper, ss)
	return AppModule{
		AppModule: stakingModule,
		keeper:    keeper,
		subspace:  ss,
	}
}

// RegisterServices registers module services.
// NOTE: Overriding this method as not doing so will cause a panic
// when trying to force this custom keeper into a bankkeeper.BaseKeeper
func (am AppModule) RegisterServices(cfg module.Configurator) {
	// types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(&am.keeper))
	types.RegisterMsgServer(cfg.MsgServer(), customstakingkeeper.NewMsgServerImpl(am.keeper.Keeper, am.keeper))
	querier := keeper.Querier{Keeper: &am.keeper.Keeper}
	types.RegisterQueryServer(cfg.QueryServer(), querier)

	m := stakingkeeper.NewMigrator(&am.keeper.Keeper, am.subspace)
	if err := cfg.RegisterMigration(types.ModuleName, 1, m.Migrate1to2); err != nil {
		panic(fmt.Sprintf("failed to migrate x/staking from version 1 to 2: %v", err))
	}

	if err := cfg.RegisterMigration(types.ModuleName, 2, m.Migrate2to3); err != nil {
		panic(fmt.Sprintf("failed to migrate x/staking from version 2 to 3: %v", err))
	}

	if err := cfg.RegisterMigration(types.ModuleName, 3, m.Migrate3to4); err != nil {
		panic(fmt.Sprintf("failed to migrate x/staking from version 3 to 4: %v", err))
	}
}