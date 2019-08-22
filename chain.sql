-- Open Source License
-- Copyright (c) 2019 Vincent Bernardoff <vb@luminar.eu.org>
-- Copyright (c) 2019 Nomadic Labs <contact@nomadic-labs.com>
--
-- Permission is hereby granted, free of charge, to any person obtaining a
-- copy of this software and associated documentation files (the "Software"),
-- to deal in the Software without restriction, including without limitation
-- the rights to use, copy, modify, merge, publish, distribute, sublicense,
-- and/or sell copies of the Software, and to permit persons to whom the
-- Software is furnished to do so, subject to the following conditions:
--
-- The above copyright notice and this permission notice shall be included
-- in all copies or substantial portions of the Software.
--
-- THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
-- IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
-- FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
-- THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
-- LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
-- FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
-- DEALINGS IN THE SOFTWARE.

-- Version 0.1 (03 June 2019)

create table chain (
  hash char(15) primary key
  -- chain id
);

-- this table inlines blocks and block headers
-- see lib_base/block_header.mli
create table block (
  hash char(51) unique,
  -- Block hash.
  -- 51 = 32 bytes hashes encoded in b58check + length of prefix "B"
  -- see lib_crypto/base58.ml
  level int,
  -- Height of the block, from the genesis block.
  proto int not null,
  -- Number of protocol changes since genesis modulo 256.
  predecessor char(51) not null,
  -- Hash of the preceding block.
  timestamp timestamp not null,
  -- Timestamp at which the block is claimed to have been created.
  validation_passes smallint not null,
  -- Number of validation passes (also number of lists of operations).
  merkle_root char(53) not null,
  -- see [operations_hash]
  -- Hash of the list of lists (actually root hashes of merkle trees)
  -- of operations included in the block. There is one list of
  -- operations per validation pass.
  -- 53 = 32 bytes hashes encoded in b58 check + "LLo" prefix
  fitness varchar(64) not null,
  -- A sequence of sequences of unsigned bytes, ordered by length and
  -- then lexicographically. It represents the claimed fitness of the
  -- chain ending in this block.
  context_hash char(52) not null,
  -- Hash of the state of the context after application of this block.
  primary key (hash, level)
);

-- Record extracted from
-- Alpha_block_services.operation
create table operation (
  hash char(51) primary key, -- operation hash
  chain char(15) not null, -- chain id
  block_hash char(51) not null, -- block hash
  foreign key (chain) references chain (hash), -- link to [chain] table
  foreign key (block_hash) references block (hash) -- link to [block] table
);
create index operation_chain on operation (chain);
create index operation_block on operation (block_hash);

-- protocol-specific content of an operation
create table operation_alpha (
  hash char(51) not null,
  id smallint not null,
  -- index of op in contents_list
  operation_kind smallint not null,
  -- from mezos/chain_db.ml
  -- see proto_alpha/operation_repr.ml
  -- (this would better be called "kind")
  -- type of operation alpha
  -- 0: Endorsement
  -- 1: Seed_nonce_revelation
  -- 2: double_endorsement_evidence
  -- 3: Double_baking_evidence
  -- 4: Activate_account
  -- 5: Proposals
  -- 6: Ballot
  -- 7: Manager_operation { operation = Reveal _ ; _ }
  -- 8: Manager_operation { operation = Transaction _ ; _ }
  -- 9: Manager_operation { operation = Origination _ ; _ }
  -- 10: Manager_operation { operation = Delegation _ ; _ }
  primary key (hash, id),
  foreign key (hash) references operation (hash)
);
create index operation_alpha_cat on operation_alpha (operation_kind);

-- implicit accounts (including deactivated ones)
create table implicit (
  pkh char(36) primary key,
  -- b58-encoded public key hash: tz1/tz1/tz3...
  activated char(51),
  -- hash of block at which activation was performed
  -- (see mezos/chain_db.ml/upsert_activated)
  revealed char(51),
  -- hash of block at which revelation was performed
  -- (see mezos/chain_db.ml/upsert_activated)
  pk varchar(55),
  -- Full public key (optional)
  foreign key (activated) references block (hash),
  foreign key (revealed)  references block (hash)
);
create index implicit_activated on implicit (activated);
create index implicit_revealed on implicit (revealed);

-- endorsements
-- /!\ It seems like this table is not filled by Mezos
-- currently. I could not find the related code, at least.
create table endorsement (
  block_hash char(51),
  op char(51),
  id smallint,
  level int,
  pkh char(36),
  slot smallint,
  primary key (block_hash, op, id, level, pkh, slot),
  foreign key (block_hash) references block (hash),
  foreign key (op)  references operation (hash),
  foreign key (pkh) references implicit (pkh)
);

/* Block info is obtained through
   lib_shell_services/Make()()/info
   instantiated in
   proto_alpha/lib_client/proto_alpha.ml
} */

-- Block table.
-- from the doc:
-- "level_position = cycle * blocks_per_cycle + cycle_position"
create table block_alpha (
  hash char(51) primary key,
  -- block hash
  baker char(36) not null,
  -- pkh of baker
  level_position int not null,
  /* Verbatim from lib_protocol/level_repr:
     The level of the block relative to the block that
     starts protocol alpha. This is specific to the
     protocol alpha. Other protocols might or might not
     include a similar notion.
  */
  cycle int not null,
  -- cycle
  cycle_position int not null,
  /* Verbatim from lib_protocol/level_repr:
     The current level of the block relative to the first
     block of the current cycle.
  */
  voting_period int not null,
  /* increasing integer.
     from proto_alpha/level_repr:
     voting_period = level_position / blocks_per_voting_period */
  voting_period_position int not null,
  -- voting_period_position = remainder(level_position / blocks_per_voting_period)
  voting_period_kind smallint not null,
  /* Proposal = 0
     Testing_vote = 1
     Testing = 2
     Promotion_vote = 3
     Defined implicitly in mezos/tezos_sql.ml via use of Obj.magic on the
     type proto_alpha/Voting_period.kind */
  consumed_gas varchar(64) not null,
  /* total gas consumed by block. Arbitrary-precision integer, max set by protocol
     represented as hex dump of binary (little-endian) form of unsigned integer.
     Note: in Mezos implem, this value cannot be negative because of Z.of_bits
     (which is reasonable). */
  foreign key (hash) references block (hash),
  foreign key (baker) references implicit (pkh)
);
create index block_alpha_baker on block_alpha (baker);
create index block_alpha_level_position on block_alpha (level_position);
create index block_alpha_cycle on block_alpha (cycle);
create index block_alpha_cycle_position on block_alpha (cycle_position);

-- deactivated accounts
create table deactivated (
  pkh char(36) not null,
  -- pkh of the deactivated account (tz1...)
  block_hash char(51) not null,
  -- block hash at which deactivation occured
  primary key (pkh, block_hash),
  foreign key (pkh) references implicit (pkh),
  foreign key (block_hash) references block (hash)
);

-- contract (implicit:tz1... or originated:KT1...) table
-- two ways of updating this table:
-- - on bootstrap, scanning preexisting contracts
-- - when scanning ops, looking at an origination/revelation
create table contract (
  address char(36) primary key,
  -- contract address, b58check format
  block_hash char(51) not null,
  -- block hash
  mgr char(36),
  -- manager
  delegate char(36),
  -- delegate
  spendable bool not null,
  -- spendable flag, soon obsolete!
  delegatable bool not null,
  -- delegatable flag, soon obsolete?
  credit bigint,
  -- credit
  preorig char(36),
  -- comment from proto_alpha/apply:
  -- The preorigination field is only used to early return
  -- the address of an originated contract in Michelson.
  -- It cannot come from the outside.
  script text,
  -- Json-encoded Micheline script
  foreign key (block_hash) references block (hash),
  foreign key (mgr) references implicit (pkh),
  foreign key (delegate) references implicit (pkh),
  foreign key (preorig) references contract (address)
);
create index contract_block on contract (block_hash);
create index contract_mgr on contract (mgr);
create index contract_delegate on contract (delegate);
create index contract_preorig on contract (preorig);

-- transaction table
create table tx (
  operation_hash char(51) not null,
  -- operation hash (starts with "o", see lib_crypto/base58)
  op_id smallint not null,
  -- index of the operation in the block's list of operations
  source char(36) not null,
  -- source address
  destination char(36) not null,
  -- dest address
  fee bigint not null,
  -- fees
  amount bigint not null,
  -- amount
  parameters text,
  -- optional parameters to contract in json-encoded Micheline
  primary key (operation_hash, op_id),
  foreign key (operation_hash, op_id) references operation_alpha (hash, id),
  foreign key (source) references contract (address),
  foreign key (destination) references contract (address)
);
create index tx_source on tx (source);
create index tx_destination on tx (destination);

-- origination table
create table origination (
  operation_hash char(51) not null,
  -- operation hash
  op_id smallint not null,
  -- index of the operation in the block's list of operations
  source char(36) not null,
  -- source of origination op
  k char(36) not null,
  -- address of originated contract
  primary key (operation_hash, op_id),
  foreign key (operation_hash, op_id) references operation_alpha (hash, id),
  foreign key (source) references contract (address),
  foreign key (k) references contract (address)
);
create index origination_source on origination (source);
create index origination_k on origination (k);

create table delegation (
  operation_hash char(51) not null,
  -- operation hash
  op_id smallint not null,
  -- index of the operation in the block's list of operations
  source char(36) not null,
  -- source of the delegation op
  pkh char(36),
  -- optional delegate
  primary key (operation_hash, op_id),
  foreign key (operation_hash, op_id) references operation_alpha (hash, id),
  foreign key (source) references contract (address),
  foreign key (pkh) references implicit (pkh)
);
create index delegation_source on delegation (source);
create index delegation_pkh on delegation (pkh);

create table balance (
  block_hash char(51) not null,
  -- block hash
  operation_hash char(51),
  -- operation hash
  op_id smallint,
  -- index of the operation in the blocks list of operations
  balance_kind smallint not null,
  -- balance kind:
  -- 0 : Contract
  -- 1 : Rewards
  -- 2 : Fees
  -- 3 : Deposits
  -- see proto_alpha/delegate_storage.ml/balance
  contract_address char(36) not null,
  -- b58check encoded address of contract (either implicit or originated)
  cycle int,
  -- cycle
  diff bigint not null,
  -- balance update
  -- credited if positve
  -- debited if negative
  foreign key (block_hash) references block (hash),
  foreign key (operation_hash, op_id) references operation_alpha (hash, id),
  foreign key (contract_address) references contract (address)
);
create index balance_block   on balance (block_hash);
create index balance_op    on balance (operation_hash, op_id);
create index balance_cat   on balance (balance_kind);
create index balance_k     on balance (contract_address);
create index balance_cycle on balance (cycle);

-- snapshots
-- the snapshot block for a given cycle is obtained as follows
-- at the last block of cycle n, the snapshot block for cycle n+6 is selected
-- Use [Storage.Roll.Snapshot_for_cycle.get ctxt cycle] in proto_alpha to
-- obtain this value.
-- RPC: /chains/main/blocks/${block}/context/raw/json/cycle/${cycle}
-- where:
-- ${block} denotes a block (either by hash or level)
-- ${cycle} denotes a cycle which must be in [cycle_of(level)-5,cycle_of(level)+7]
create table snapshot (
  cycle int,
  level int,
  primary key (cycle, level)
);

-- Could be useful for baking.
-- create table delegate (
--   cycle int not null,
--   level int not null,
--   pkh char(36) not null,
--   balance bigint not null,
--   frozen_balance bigint not null,
--   staking_balance bigint not null,
--   delegated_balance bigint not null,
--   deactivated bool not null,
--   grace smallint not null,
--   primary key (cycle, pkh),
--   foreign key (cycle, level) references snapshot (cycle, level),
--   foreign key (pkh) references implicit (pkh)
-- );

-- Delegated contract table
-- It seems this table is not filled by Mezos yet
create table delegated_contract (
  delegate char(36),
  -- tz1 of the delegate
  delegator char(36),
  -- address of the delegator (for now, KT1 but this could change)
  cycle int,
  level int,
  primary key (delegate, delegator, cycle, level),
  foreign key (delegate) references implicit (pkh),
  foreign key (delegator) references contract (address),
  foreign key (cycle, level) references snapshot (cycle, level)
);
create index delegated_contract_cycle on delegated_contract (cycle);
create index delegated_contract_level on delegated_contract (level);

-- Could be useful for baking.
-- create table stake (
--   delegate char(36) not null,
--   level int not null,
--   k char(36) not null,
--   kind smallint not null,
--   diff bigint not null,
--   primary key (delegate, level, k, kind, diff),
--   foreign key (delegate) references implicit (pkh),
--   foreign key (k) references contract (address)
-- );

-- VIEWS
create view level as
  select
   level,
   hash
  from block
  order by level asc;

--  tx_full view
create view tx_full as
  select
   operation_hash, -- operation hash
   op_id, -- index in list of operations
   b.hash as block_hash, -- block hash
   b.level as level, -- block level
   b.timestamp as timestamp, -- timestamp
   source, -- source
   k1.mgr as source_mgr, -- manager of source
   destination, -- destination
   k2.mgr as destination_mgr, -- manager of destination
   fee, -- fees
   amount, -- amount transfered
   parameters -- parameters to target contract, if any
  from tx
  join operation on tx.operation_hash = operation.hash
  join block b on operation.block_hash = b.hash
  join contract k1 on tx.source = k1.address
  join contract k2 on tx.destination = k2.address;

create view balance_full as
  select
   block.level, -- block level
   block_alpha.cycle, -- cycle
   cycle_position, -- position in cycle
   operation_hash, -- operation hash
   op_id, -- index in list of operations
   contract_address, -- address of contract
   balance_kind, -- balance kind
   diff -- balance update
  from block
  natural join block_alpha
  join balance on block.hash = balance.block_hash;
