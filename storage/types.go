package storage

import (
	"context"
	"time"
)

/*
// Types based on Postgres schema
// TODO: revisit

type Block struct {
	// Block hash.
	// 51 = 32 bytes hashes encoded in b58check + length of prefix "B"
	// see lib_crypto/base58.ml
	Hash             string
	Level            int       // Height of the block, from the genesis block.
	Proto            int       // Number of protocol changes since genesis modulo 256.
	Predecessor      string    // Hash of the preceding block.
	Timestamp        time.Time // Timestamp at which the block is claimed to have been created.
	ValidationPasses int       // Number of validation passes (also number of lists of operations).
	// see [operations_hash]
	// Hash of the list of lists (actually root hashes of merkle trees)
	// of operations included in the block. There is one list of
	// operations per validation pass.
	// 53 = 32 bytes hashes encoded in b58 check + "LLo" prefix
	MerkleRoot string
	// A sequence of sequences of unsigned bytes, ordered by length and
	// then lexicographically. It represents the claimed fitness of the
	// chain ending in this block.
	Fitness     string
	ContextHash string // Hash of the state of the context after application of this block.
}

// From the doc:
// "level_position = cycle * blocks_per_cycle + cycle_position"
type BlockAlpha struct {
	Block

	Hash  string // Block hash
	Baker string // PKH of baker
	// Verbatim from lib_protocol/level_repr:
	// The level of the block relative to the block that
	// starts protocol alpha. This is specific to the
	// protocol alpha. Other protocols might or might not
	// include a similar notion.
	LevelPosition int
	Cycle         int // Cycle
	// Verbatim from lib_protocol/level_repr:
	// The current level of the block relative to the first
	// block of the current cycle.
	CyclePosition int
	// Increasing integer.
	// From proto_alpha/level_repr:
	// voting_period = level_position / blocks_per_voting_period
	VotingPeriod         int // voting_period_position = remainder(level_position / blocks_per_voting_period)
	VotingPeriodPosition int
	// Proposal = 0
	// Testing_vote = 1
	// Testing = 2
	// Promotion_vote = 3
	// Defined implicitly in mezos/tezos_sql.ml via use of Obj.magic on the
	// type proto_alpha/Voting_period.kind
	VotingPeriodKind int
	// Total gas consumed by block. Arbitrary-precision integer, max set by protocol
	// represented as hex dump of binary (little-endian) form of unsigned integer.
	// Note: in Mezos implem, this value cannot be negative because of Z.of_bits
	// (which is reasonable).
	ConsumedGas string
}

type Operation struct {
	Hash      string // Operation hash
	Chain     string // Chain ID
	BlockHash string // Block hash
}

type OperationAlpha struct {
	Hash string
	ID   int // Index of op in contents_list
	// from mezos/chain_db.ml
	// see proto_alpha/operation_repr.ml
	// (this would better be called "kind")
	// type of operation alpha
	// 0: Endorsement
	// 1: Seed_nonce_revelation
	// 2: double_endorsement_evidence
	// 3: Double_baking_evidence
	// 4: Activate_account
	// 5: Proposals
	// 6: Ballot
	// 7: Manager_operation { operation = Reveal _ ; _ }
	// 8: Manager_operation { operation = Transaction _ ; _ }
	// 9: Manager_operation { operation = Origination _ ; _ }
	// 10: Manager_operation { operation = Delegation _ ; _ }
	OperationKind int
}

const (
	OperationEndorsement = iota
	OperationSeedNonceRevelation
	OperationDoubleEndorsementEvidence
	OperationDoubleBakingEvidence
	OperationActivateAccount
	OperationProposals
	OperationBallot
	OperationReveal
	OperationTransaction
	OperationOrigination
	OperationDelegation
)

// Implicit accounts (including deactivated ones)
type Implicit struct {
	PKH string // b58-encoded public key hash: tz1/tz1/tz3...
	// hash of block at which activation was performed
	// (see mezos/chain_db.ml/upsert_activated)
	Activated string
	// hash of block at which revelation was performed
	// (see mezos/chain_db.ml/upsert_activated)
	Revealed string
	PK       string // Full public key (optional)
}

// Endorsements
type Endorsement struct {
	BlockHash string
	Op        string
	ID        int
	Level     int
	PKH       string
	Slot      int
}

const (
	VotingPeriodProposal = iota
	VotingPeriodTestingVote
	VotingPeriodTesting
	VotingPeriodPromotionVote
)

// Deactivated accounts
type Deactivated struct {
	PKH       string // PKH of the deactivated account (tz1...)
	BlockHash string // Block hash at which deactivation occured
}

// contract (implicit:tz1... or originated:KT1...) table
// two ways of updating this table:
// - on bootstrap, scanning preexisting contracts
// - when scanning ops, looking at an origination/revelation
type Contract struct {
	Address     string // Contract address, b58check format
	BlockHash   string // Block hash
	Mgr         string // Manager
	Delegate    string // Delegate
	Spendable   bool   // Spendable flag, soon obsolete!
	Delegatable bool   // Delegatable flag, soon obsolete?
	Credit      int64  // Credit
	// Comment from proto_alpha/apply:
	// The preorigination field is only used to early return
	// the address of an originated contract in Michelson.
	// It cannot come from the outside.
	Preorig string
	Script  string // Json-encoded Micheline script
}

// Transaction table
type Tx struct {
	OperationHash string // Operation hash (starts with "o", see lib_crypto/base58)
	OpID          int    // Index of the operation in the block's list of operations
	Source        string // Source address
	Destination   string // Dest address
	Fee           int64  // Fees
	Amount        int64  // Amount
	Parameters    string // Optional parameters to contract in json-encoded Micheline
}

// Origination table
type Origination struct {
	OperationHash string // Operation hash
	OpID          int    // Index of the operation in the block's list of operations
	Source        string // Source of origination op
	K             string // Address of originated contract
}

type Delegation struct {
	OperationHash string // Operation hash
	OpID          int    // Index of the operation in the block's list of operations
	Source        string // Source of origination op
	PKH           string // Optional delegate
}

type Balance struct {
	BlockHash     string // Block hash
	OperationHash string // Operation hash
	OpID          int    // Index of the operation in the block's list of operations
	// Balance kind:
	// 0 : Contract
	// 1 : Rewards
	// 2 : Fees
	// 3 : Deposits
	// see proto_alpha/delegate_storage.ml/balance
	BalanceKind     int
	ContractAddress string // b58check encoded address of contract (either implicit or originated)
	Cycle           int    // Cycle
	// Balance update
	// credited if positve
	// debited if negative
	Diff int64
}

const (
	BalanceContract = iota
	BalanceRewards
	BalanceFees
	BalanceDeposits
)

// Snapshots
// the snapshot block for a given cycle is obtained as follows
// at the last block of cycle n, the snapshot block for cycle n+6 is selected
// Use [Storage.Roll.Snapshot_for_cycle.get ctxt cycle] in proto_alpha to
// obtain this value.
// RPC: /chains/main/blocks/${block}/context/raw/json/cycle/${cycle}
// where:
// ${block} denotes a block (either by hash or level)
// ${cycle} denotes a cycle which must be in [cycle_of(level)-5,cycle_of(level)+7]
type Snapshot struct {
	Cycle int
	Level int
}

type DelegatedContract struct {
	Delegate  string // tz1 of the delegate
	Delegator string // Address of the delegator (for now, KT1 but this could change)
	Cycle     int
	Level     int
}

*/

// API specific types
type BalanceUpdate struct {
	BlockLevel     int64     `json:"level"`
	BlockTimestamp time.Time `json:"timestamp"`
	Diff           int64     `json:"diff"`
	Value          int64     `json:"value"`
}

type BalanceStorage interface {
	GetBalanceUpdate(ctx context.Context, address string, start, end time.Time, limit int) ([]*BalanceUpdate, error)
}
