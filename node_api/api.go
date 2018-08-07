package node_api

import (
	"github.com/filecoin-project/go-filecoin/api"
	"github.com/filecoin-project/go-filecoin/node"
)

// API is an actual implementation of the filecoin core api interface.
type API struct {
	node *node.Node

	actor     *ActorAPI
	address   *AddressAPI
	block     *BlockAPI
	bootstrap *BootstrapAPI
	chain     *ChainAPI
	config    *ConfigAPI
	client    *ClientAPI
	daemon    *DaemonAPI
	dag       *DagAPI
	id        *IdAPI
	init      *InitAPI
	log       *LogAPI
	message   *MessageAPI
	miner     *MinerAPI
	mining    *MiningAPI
	mpool     *MpoolAPI
	orderbook *OrderbookAPI
	paych     *PaychAPI
	ping      *PingAPI
	show      *ShowAPI
	swarm     *SwarmAPI
	version   *VersionAPI
	wallet    *WalletAPI
}

// Assert that API fullfills the api.API interface.
var _ api.API = (*API)(nil)

// NewAPI constructs a new instance of the API.
func NewAPI(node *node.Node) api.API {
	api := &API{
		node: node,
	}
	api.actor = NewActorAPI(api)
	api.address = NewAddressAPI(api)
	api.block = NewBlockAPI(api)
	api.bootstrap = NewBootstrapAPI(api)
	api.chain = NewChainAPI(api)
	api.config = NewConfigAPI(api)
	api.client = NewClientAPI(api)
	api.daemon = NewDaemonAPI(api)
	api.dag = NewDagAPI(api)
	api.id = NewIdAPI(api)
	api.init = NewInitAPI(api)
	api.log = NewLogAPI(api)
	api.message = NewMessageAPI(api)
	api.miner = NewMinerAPI(api)
	api.mining = NewMiningAPI(api)
	api.mpool = NewMpoolAPI(api)
	api.orderbook = NewOrderbookAPI(api)
	api.paych = NewPaychAPI(api)
	api.ping = NewPingAPI(api)
	api.show = NewShowAPI(api)
	api.swarm = NewSwarmAPI(api)
	api.version = NewVersionAPI(api)
	api.wallet = NewWalletAPI(api)

	return api
}

func (api *API) Actor() api.Actor {
	return api.actor
}

func (api *API) Address() api.Address {
	return api.address
}

func (api *API) Block() api.Block {
	return api.block
}

func (api *API) Bootstrap() api.Bootstrap {
	return api.bootstrap
}

func (api *API) Chain() api.Chain {
	return api.chain
}

func (api *API) Config() api.Config {
	return api.config
}

func (api *API) Client() api.Client {
	return api.client
}

func (api *API) Daemon() api.Daemon {
	return api.daemon
}

func (api *API) Dag() api.Dag {
	return api.dag
}

func (api *API) Id() api.Id {
	return api.id
}

func (api *API) Init() api.Init {
	return api.init
}

func (api *API) Log() api.Log {
	return api.log
}

func (api *API) Message() api.Message {
	return api.message
}

func (api *API) Miner() api.Miner {
	return api.miner
}

func (api *API) Mining() api.Mining {
	return api.mining
}

func (api *API) Mpool() api.Mpool {
	return api.mpool
}

func (api *API) Orderbook() api.Orderbook {
	return api.orderbook
}

func (api *API) Paych() api.Paych {
	return api.paych
}

func (api *API) Ping() api.Ping {
	return api.ping
}

func (api *API) Show() api.Show {
	return api.show
}

func (api *API) Swarm() api.Swarm {
	return api.swarm
}

func (api *API) Version() api.Version {
	return api.version
}

func (api *API) Wallet() api.Wallet {
	return api.wallet
}