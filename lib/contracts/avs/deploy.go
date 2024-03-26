package avs

import (
	"context"
	"math/big"

	"github.com/omni-network/omni/contracts/bindings"
	"github.com/omni-network/omni/lib/chainids"
	"github.com/omni-network/omni/lib/contracts"
	"github.com/omni-network/omni/lib/create3"
	"github.com/omni-network/omni/lib/errors"
	"github.com/omni-network/omni/lib/ethclient/ethbackend"
	"github.com/omni-network/omni/lib/netconf"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

const (
	metadataURI = "https://raw.githubusercontent.com/omni-network/omni/main/static/avs-metadata.json"
)

//nolint:gochecknoglobals // static abi type
var (
	avsABI   = mustGetABI(bindings.OmniAVSMetaData)
	proxyABI = mustGetABI(bindings.TransparentUpgradeableProxyMetaData)
)

type DeploymentConfig struct {
	Create3Factory   common.Address
	Create3Salt      string
	Eigen            EigenDeployments
	Deployer         common.Address
	Owner            common.Address
	ProxyAdmin       common.Address
	Portal           common.Address
	OmniChainID      uint64
	MetadataURI      string
	StrategyParams   []StrategyParam
	EthStakeInbox    common.Address
	MinOperatorStake *big.Int
	MaxOperatorCount uint32
	AllowlistEnabled bool
	ExpectedAddr     common.Address
}

func (cfg DeploymentConfig) Validate() error {
	if (cfg.Create3Factory == common.Address{}) {
		return errors.New("create3 factory not set")
	}
	if cfg.Create3Salt == "" {
		return errors.New("create3 salt not set")
	}
	if err := cfg.Eigen.Validate(); err != nil {
		return errors.Wrap(err, "eigen deployments")
	}
	if (cfg.Deployer == common.Address{}) {
		return errors.New("deployer is zero")
	}
	if (cfg.Owner == common.Address{}) {
		return errors.New("owner is zero")
	}
	if (cfg.ProxyAdmin == common.Address{}) {
		return errors.New("proxy admin is zero")
	}
	if cfg.MetadataURI == "" {
		return errors.New("metadata uri not set")
	}
	if cfg.MinOperatorStake == nil {
		return errors.New("min operator stake not set")
	}
	if cfg.MaxOperatorCount == 0 {
		return errors.New("max operator count not set")
	}

	return nil
}

func getDeployCfg(chainID uint64, network netconf.ID) (DeploymentConfig, error) {
	if !chainids.IsMainnetOrTestnet(chainID) && network == netconf.Devnet {
		return devnetCfg(), nil
	}

	if chainID == chainids.Holesky && network == netconf.Testnet {
		return testnetCfg(), nil
	}

	if !chainids.IsMainnet(chainID) && network == netconf.Staging {
		return stagingCfg(), nil
	}

	return DeploymentConfig{}, errors.New("unsupported chain for network", "chain_id", chainID, "network", network)
}

func testnetCfg() DeploymentConfig {
	return DeploymentConfig{
		Create3Factory:   contracts.TestnetCreate3Factory(),
		Create3Salt:      contracts.AVSSalt(netconf.Testnet),
		Deployer:         contracts.TestnetDeployer(),
		Owner:            contracts.TestnetAVSAdmin(),
		ProxyAdmin:       contracts.TestnetProxyAdmin(),
		Eigen:            holeskyEigenDeployments(),
		StrategyParams:   holeskeyStrategyParams(),
		MetadataURI:      metadataURI,
		OmniChainID:      netconf.Testnet.Static().OmniExecutionChainID,
		MinOperatorStake: big.NewInt(1e18), // 1 ETH
		MaxOperatorCount: 200,
		AllowlistEnabled: false,
		ExpectedAddr:     contracts.TestnetAVS(),
	}
}

func stagingCfg() DeploymentConfig {
	return DeploymentConfig{
		Create3Factory:   contracts.StagingCreate3Factory(),
		Create3Salt:      contracts.AVSSalt(netconf.Staging),
		Deployer:         contracts.StagingDeployer(),
		Owner:            contracts.StagingAVSAdmin(),
		ProxyAdmin:       contracts.StagingProxyAdmin(),
		Eigen:            devnetEigenDeployments,
		StrategyParams:   devnetStrategyParams(),
		MetadataURI:      metadataURI,
		OmniChainID:      netconf.Staging.Static().OmniExecutionChainID,
		MinOperatorStake: big.NewInt(1e18), // 1 ETH
		MaxOperatorCount: 10,
		AllowlistEnabled: true,
		ExpectedAddr:     contracts.StagingAVS(),
	}
}

func devnetCfg() DeploymentConfig {
	return DeploymentConfig{
		Create3Factory:   contracts.DevnetCreate3Factory(),
		Create3Salt:      contracts.AVSSalt(netconf.Devnet),
		Deployer:         contracts.DevnetDeployer(),
		Owner:            contracts.DevnetAVSAdmin(),
		ProxyAdmin:       contracts.DevnetProxyAdmin(),
		Eigen:            devnetEigenDeployments,
		MetadataURI:      metadataURI,
		OmniChainID:      netconf.Devnet.Static().OmniExecutionChainID,
		StrategyParams:   devnetStrategyParams(),
		EthStakeInbox:    common.HexToAddress("0x1234"), // TODO: replace with actual address
		MinOperatorStake: big.NewInt(1e18),              // 1 ETH
		MaxOperatorCount: 10,
		AllowlistEnabled: true,
		ExpectedAddr:     contracts.DevnetAVS(),
	}
}

func AddrForNetwork(network netconf.ID) (common.Address, bool) {
	switch network {
	case netconf.Mainnet:
		return contracts.MainnetAVS(), true
	case netconf.Testnet:
		return contracts.TestnetAVS(), true
	case netconf.Staging:
		return contracts.StagingAVS(), true
	case netconf.Devnet:
		return contracts.DevnetAVS(), true
	default:
		return common.Address{}, false
	}
}

// IsDeployed checks if the OmniAVS contract is deployed to the provided backend
// to its expected network address.
func IsDeployed(ctx context.Context, network netconf.ID, backend *ethbackend.Backend) (bool, contracts.Deployment, error) {
	chainID, err := backend.ChainID(ctx)
	if err != nil {
		return false, contracts.Deployment{}, errors.Wrap(err, "chain id")
	}

	cfg, err := getDeployCfg(chainID.Uint64(), network)
	if err != nil {
		return false, contracts.Deployment{}, errors.Wrap(err, "get deployment config")
	}

	factory, err := bindings.NewCreate3(cfg.Create3Factory, backend)
	if err != nil {
		return false, contracts.Deployment{}, errors.Wrap(err, "new create3")
	}

	salt := create3.HashSalt(cfg.Create3Salt)
	height, err := factory.GetDeployedHeight(nil, cfg.Deployer, salt)
	if err != nil {
		return false, contracts.Deployment{}, errors.Wrap(err, "get deployed height")
	}

	if height.Uint64() == 0 {
		return false, contracts.Deployment{}, nil
	}

	deployment := contracts.Deployment{
		Address:     create3.Address(cfg.Create3Factory, cfg.Create3Salt, cfg.Deployer),
		BlockHeight: height.Uint64(),
	}

	return true, deployment, nil
}

// DeployIfNeeded deploys a new AVS contract if it is not already deployed.
func DeployIfNeeded(ctx context.Context, network netconf.ID, backend *ethbackend.Backend) (contracts.Deployment, error) {
	deployed, deployment, err := IsDeployed(ctx, network, backend)
	if err != nil {
		return contracts.Deployment{}, errors.Wrap(err, "is deployed")
	}
	if deployed {
		return deployment, nil
	}

	return Deploy(ctx, network, backend)
}

// Deploy deploys the AVS contract and returns the address and receipt.
// It only allows deployments to explicitly supported chains.
func Deploy(ctx context.Context, network netconf.ID, backend *ethbackend.Backend) (contracts.Deployment, error) {
	chainID, err := backend.ChainID(ctx)
	if err != nil {
		return contracts.Deployment{}, errors.Wrap(err, "chain id")
	}

	cfg, err := getDeployCfg(chainID.Uint64(), network)
	if err != nil {
		return contracts.Deployment{}, errors.Wrap(err, "get deployment config")
	}

	return deploy(ctx, cfg, backend)
}

func deploy(ctx context.Context, cfg DeploymentConfig, backend *ethbackend.Backend) (contracts.Deployment, error) {
	if err := cfg.Validate(); err != nil {
		return contracts.Deployment{}, errors.Wrap(err, "validate config")
	}

	deployerTxOpts, err := backend.BindOpts(ctx, cfg.Deployer)
	if err != nil {
		return contracts.Deployment{}, errors.Wrap(err, "bind deployer opts")
	}

	ownerTxOpts, err := backend.BindOpts(ctx, cfg.Owner)
	if err != nil {
		return contracts.Deployment{}, errors.Wrap(err, "bind owner opts")
	}

	factory, err := bindings.NewCreate3(cfg.Create3Factory, backend)
	if err != nil {
		return contracts.Deployment{}, errors.Wrap(err, "new create3")
	}

	salt := create3.HashSalt(cfg.Create3Salt)

	addr, err := factory.GetDeployedAddr(nil, deployerTxOpts.From, salt)
	if err != nil {
		return contracts.Deployment{}, errors.Wrap(err, "get deployed")
	} else if (cfg.ExpectedAddr != common.Address{}) && addr != cfg.ExpectedAddr {
		return contracts.Deployment{}, errors.New("unexpected address", "expected", cfg.ExpectedAddr, "actual", addr)
	}

	impl, tx, _, err := bindings.DeployOmniAVS(deployerTxOpts, backend, cfg.Eigen.DelegationManager, cfg.Eigen.AVSDirectory)
	if err != nil {
		return contracts.Deployment{}, errors.Wrap(err, "deploy impl")
	}

	deployReceipt, err := backend.WaitMined(ctx, tx)
	if err != nil {
		return contracts.Deployment{}, errors.Wrap(err, "wait mined portal")
	} else if deployReceipt.Status != ethtypes.ReceiptStatusSuccessful {
		return contracts.Deployment{}, errors.New("deploy impl failed")
	}

	initCode, err := packInitCode(cfg, impl)
	if err != nil {
		return contracts.Deployment{}, errors.Wrap(err, "pack init code")
	}

	tx, err = factory.Deploy(deployerTxOpts, salt, initCode)
	if err != nil {
		return contracts.Deployment{}, errors.Wrap(err, "deploy proxy")
	}

	receipt, err := backend.WaitMined(ctx, tx)
	if err != nil {
		return contracts.Deployment{}, errors.Wrap(err, "wait mined proxy")
	} else if receipt.Status != ethtypes.ReceiptStatusSuccessful {
		return contracts.Deployment{}, errors.New("deploy proxy failed")
	}

	deployment := contracts.Deployment{
		Address:     addr,
		BlockHeight: receipt.BlockNumber.Uint64(),
	}

	// the contract has been deployed. transactions below are admin calls
	// TODO: move to avs initializer, so that deployment is the last transaction

	avs, err := bindings.NewOmniAVS(addr, backend)
	if err != nil {
		return deployment, errors.Wrap(err, "bind avs")
	}

	if !cfg.AllowlistEnabled {
		// only wait mained second admin call below (SetMetadataURI)
		_, err = avs.DisableAllowlist(ownerTxOpts)
		if err != nil {
			return deployment, errors.Wrap(err, "disable allowlist")
		}
	}

	tx, err = avs.SetMetadataURI(ownerTxOpts, cfg.MetadataURI)
	if err != nil {
		return deployment, errors.Wrap(err, "set metadata uri")
	}

	receipt, err = backend.WaitMined(ctx, tx)
	if err != nil {
		return deployment, errors.Wrap(err, "wait mined set metadata uri")
	}
	if receipt.Status != ethtypes.ReceiptStatusSuccessful {
		return deployment, errors.New("set metadata uri failed")
	}

	return deployment, nil
}

func packInitCode(cfg DeploymentConfig, impl common.Address) ([]byte, error) {
	initializer, err := packInitialzer(cfg)
	if err != nil {
		return nil, err
	}

	return contracts.PackInitCode(proxyABI, bindings.TransparentUpgradeableProxyBin, impl, cfg.ProxyAdmin, initializer)
}

// packInitializer encodes the initializer parameters for the AVS contract.
func packInitialzer(cfg DeploymentConfig) ([]byte, error) {
	enc, err := avsABI.Pack("initialize",
		cfg.Owner, cfg.Portal, cfg.OmniChainID, cfg.EthStakeInbox,
		cfg.MinOperatorStake, cfg.MaxOperatorCount, strategyParams(cfg))

	if err != nil {
		return nil, errors.Wrap(err, "pack initializer")
	}

	return enc, nil
}

// strategyParams converts the deployment config's strategy params to the.
func strategyParams(cfg DeploymentConfig) []bindings.IOmniAVSStrategyParam {
	params := make([]bindings.IOmniAVSStrategyParam, len(cfg.StrategyParams))
	for i, sp := range cfg.StrategyParams {
		params[i] = bindings.IOmniAVSStrategyParam{
			Strategy:   sp.Strategy,
			Multiplier: sp.Multiplier,
		}
	}

	return params
}

// mustGetABI returns the metadata's ABI as an abi.ABI type.
// It panics on error.
func mustGetABI(metadata *bind.MetaData) *abi.ABI {
	abi, err := metadata.GetAbi()
	if err != nil {
		panic(err)
	}

	return abi
}
