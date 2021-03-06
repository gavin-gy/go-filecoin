package fastesting

import (
	"context"
	"io/ioutil"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ipfs/go-ipfs-files"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/go-filecoin/testhelpers"
	"github.com/filecoin-project/go-filecoin/tools/fast"
	"github.com/filecoin-project/go-filecoin/tools/fast/series"
	localplugin "github.com/filecoin-project/go-filecoin/tools/iptb-plugins/filecoin/local"
	"github.com/filecoin-project/go-filecoin/types"
)

// TestEnvironment provides common setup for writing tests using FAST
type TestEnvironment struct {
	fast.Environment

	t   *testing.T
	ctx context.Context

	pluginName string
	pluginOpts map[string]string

	fastenvOpts fast.EnvironmentOpts

	GenesisMiner *fast.Filecoin
}

// NewTestEnvironment creates a TestEnvironment with a basic setup for writing tests using the FAST library.
func NewTestEnvironment(ctx context.Context, t *testing.T, fastenvOpts fast.EnvironmentOpts) (context.Context, *TestEnvironment) {
	require := require.New(t)

	// Create a directory for the test using the test name (mostly for FAST)
	// Replace the forward slash as tempdir can't handle them
	dir, err := ioutil.TempDir("", strings.Replace(t.Name(), "/", ".", -1))
	require.NoError(err)

	// Create an environment that includes a genesis block with 1MM FIL
	env, err := fast.NewEnvironmentMemoryGenesis(big.NewInt(1000000), dir, types.TestProofsMode)
	require.NoError(err)

	// Setup options for nodes.
	options := make(map[string]string)
	options[localplugin.AttrLogJSON] = "1"                                        // Enable JSON logs
	options[localplugin.AttrLogLevel] = "5"                                       // Set log level to Debug
	options[localplugin.AttrFilecoinBinary] = testhelpers.MustGetFilecoinBinary() // Get the filecoin binary

	genesisURI := env.GenesisCar()
	genesisMiner, err := env.GenesisMiner()
	require.NoError(err)

	fastenvOpts = fast.EnvironmentOpts{
		InitOpts:   append([]fast.ProcessInitOption{fast.POGenesisFile(genesisURI)}, fastenvOpts.InitOpts...),
		DaemonOpts: append([]fast.ProcessDaemonOption{fast.POBlockTime(time.Millisecond)}, fastenvOpts.DaemonOpts...),
	}

	// Setup the first node which is used to help coordinate the other nodes by providing
	// funds, mining for the network, etc
	genesis, err := env.NewProcess(ctx, localplugin.PluginName, options, fastenvOpts)
	require.NoError(err)

	err = series.SetupGenesisNode(ctx, genesis, genesisMiner.Address, files.NewReaderFile(genesisMiner.Owner))
	require.NoError(err)

	// Define a MiningOnce function which will bet set on the context to provide
	// a way to mine blocks in the series used during testing
	var MiningOnce series.MiningOnceFunc = func() {
		_, err := genesis.MiningOnce(ctx)
		require.NoError(err)
	}

	ctx = series.SetCtxMiningOnce(ctx, MiningOnce)
	ctx = series.SetCtxSleepDelay(ctx, time.Second)

	return ctx, &TestEnvironment{
		Environment:  env,
		t:            t,
		ctx:          ctx,
		pluginName:   localplugin.PluginName,
		pluginOpts:   options,
		fastenvOpts:  fastenvOpts,
		GenesisMiner: genesis,
	}
}

// RequireNewNode builds a new node for the environment
func (env *TestEnvironment) RequireNewNode() *fast.Filecoin {
	require := require.New(env.t)

	p, err := env.NewProcess(env.ctx, env.pluginName, env.pluginOpts, env.fastenvOpts)
	require.NoError(err)

	return p
}

// RequireNewNodeStarted builds a new node using RequireNewNode, then initializes
// and starts it
func (env *TestEnvironment) RequireNewNodeStarted() *fast.Filecoin {
	require := require.New(env.t)

	p := env.RequireNewNode()

	err := series.InitAndStart(env.ctx, p)
	require.NoError(err)

	return p
}

// RequireNewNodeConnected builds a new node using RequireNewNodeStarted, then
// connect it to the environment GenesisMiner node
func (env *TestEnvironment) RequireNewNodeConnected() *fast.Filecoin {
	require := require.New(env.t)

	p := env.RequireNewNodeStarted()

	err := series.Connect(env.ctx, env.GenesisMiner, p)
	require.NoError(err)

	return p
}

// RequireNodeNodeWithFunds builds a new node using RequireNewNodeStarted, then
// sends it funds from the environment GenesisMiner node
func (env *TestEnvironment) RequireNewNodeWithFunds(funds int) *fast.Filecoin {
	require := require.New(env.t)

	p := env.RequireNewNodeConnected()

	err := series.SendFilecoinDefaults(env.ctx, env.GenesisMiner, p, funds)
	require.NoError(err)

	return p
}
