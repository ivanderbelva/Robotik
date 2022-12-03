package consensus

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
)

// ChainHeaderReader menampung koleksi method untuk mengakses verifikasi blockchain lokal
type ChainHeaderReader interface {
	// Menerima konfigurasi blockchain
	Config() *params.ChainConfig

	// Menerima header saat ini dari local chain
	CurrentHeader() *types.Header

	// Menerima header block dari database dengan hash nomer
	GetHeader(hash common.Hash, number uint64) *types.Header

	// Menerima header block dari databse dengan nomer
	GetHeaderByNumber(number uint64) *types.Header

	// Menerima header block dari databse dengan hash
	GetHeaderByHash(hash common.Hash) *types.Header

	// Menerima kesulitan total dari database
	GetTd(hash common.Hash, number uint64) *big.Int
}

// ChainReader menampung koleksi method untuk mengakses lokal blockchain ketika verifikasi
type ChainReader interface {
	ChainHeaderReader

	// Menerima block dari database dengan hash nomer
	GetBlock(hash common.Hash, number uint64) *types.Block
}

// Algoritma mesin konsensus
type Engine interface {
	// Author retrieves the Ethereum address of the account that minted the given
	// block, which may be different from the header's coinbase if a consensus
	// engine is based on signatures.
	// Author menerima alamat akun yang diberi block
	Author(header *types.Header) (common.Address, error)

	// VerifyHeader checks whether a header conforms to the consensus rules of a
	// given engine. Verifying the seal may be done optionally here, or explicitly
	// via the VerifySeal method.
	// Mengecek apakah header berpengaturan consensus di engine
	VerifyHeader(chain ChainHeaderReader, header *types.Header, seal bool) error

	// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
	// concurrently. The method returns a quit channel to abort the operations and
	// a results channel to retrieve the async verifications (the order is that of
	// the input slice).
	// Mirip dengan VerifyHeader namun memverifikasi secara batch
	VerifyHeaders(chain ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error)

	// VerifyUncles verifies that the given block's uncles conform to the consensus
	// rules of a given engine.
	// Memverifikasi block ber uncle consensus
	VerifyUncles(chain ChainReader, block *types.Block) error

	// Prepare initializes the consensus fields of a block header according to the
	// rules of a particular engine. The changes are executed inline.
	// Menyiapkan inisialisasi consensus block header sesua peraturan mesin
	Prepare(chain ChainHeaderReader, header *types.Header) error

	// Finalize runs any post-transaction state modifications (e.g. block rewards)
	// but does not assemble the block.
	// Note: The block header and state database might be updated to reflect any
	// consensus rules that happen at finalization (e.g. block rewards).
	// Menjalankan modifikasi state pasca transaksi
	Finalize(chain ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
		uncles []*types.Header)

	// FinalizeAndAssemble runs any post-transaction state modifications (e.g. block
	// rewards) and assembles the final block.
	// Note: The block header and state database might be updated to reflect any
	// consensus rules that happen at finalization (e.g. block rewards).
	// Menjalankan modifikasi state pasca transaksi
	FinalizeAndAssemble(chain ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
		uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error)

	// Seal generates a new sealing request for the given input block and pushes
	// the result into the given channel.
	// Note, the method returns immediately and will send the result async. More
	// than one result may also be returned depending on the consensus algorithm.
	// Membuat permintaan penyegelan block input yang diberikan
	Seal(chain ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error

	// Mengembalikan block hash yang berkaitan dengan penyegelan
	SealHash(header *types.Header) common.Hash

	// Algoritma untuk pengaturan kesulitan untuk block yang sesuai
  CalcDifficulty(chain ChainHeaderReader, time uint64, parent *types.Header) *big.Int

	// Mengembalikan API RPC yang mesin consensus ini sediakan
	APIs(chain ChainHeaderReader) []rpc.API

	// Menutup thread di background yang diatur oleh mesin konsensus
	Close() error
}

// PoW merupakan mesin consensus berbasis proof-of-work
type PoW interface {
	Engine

	// Mengembalikan hashrate mining saat ini dari proof-of-work mesin consensus
	Hashrate() float64
}
