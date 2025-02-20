// Code generated by ent, DO NOT EDIT.

package runtime

import (
	"time"

	"github.com/google/uuid"
	"github.com/omni-network/omni/explorer/db/ent/block"
	"github.com/omni-network/omni/explorer/db/ent/chain"
	"github.com/omni-network/omni/explorer/db/ent/msg"
	"github.com/omni-network/omni/explorer/db/ent/receipt"
	"github.com/omni-network/omni/explorer/db/ent/schema"
	"github.com/omni-network/omni/explorer/db/ent/xprovidercursor"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	blockFields := schema.Block{}.Fields()
	_ = blockFields
	// blockDescBlockHash is the schema descriptor for BlockHash field.
	blockDescBlockHash := blockFields[2].Descriptor()
	// block.BlockHashValidator is a validator for the "BlockHash" field. It is called by the builders before save.
	block.BlockHashValidator = blockDescBlockHash.Validators[0].(func([]byte) error)
	// blockDescTimestamp is the schema descriptor for Timestamp field.
	blockDescTimestamp := blockFields[3].Descriptor()
	// block.DefaultTimestamp holds the default value on creation for the Timestamp field.
	block.DefaultTimestamp = blockDescTimestamp.Default.(time.Time)
	// blockDescCreatedAt is the schema descriptor for CreatedAt field.
	blockDescCreatedAt := blockFields[4].Descriptor()
	// block.DefaultCreatedAt holds the default value on creation for the CreatedAt field.
	block.DefaultCreatedAt = blockDescCreatedAt.Default.(time.Time)
	chainFields := schema.Chain{}.Fields()
	_ = chainFields
	// chainDescUUID is the schema descriptor for UUID field.
	chainDescUUID := chainFields[0].Descriptor()
	// chain.DefaultUUID holds the default value on creation for the UUID field.
	chain.DefaultUUID = chainDescUUID.Default.(func() uuid.UUID)
	// chainDescCreatedAt is the schema descriptor for CreatedAt field.
	chainDescCreatedAt := chainFields[1].Descriptor()
	// chain.DefaultCreatedAt holds the default value on creation for the CreatedAt field.
	chain.DefaultCreatedAt = chainDescCreatedAt.Default.(time.Time)
	msgHooks := schema.Msg{}.Hooks()
	msg.Hooks[0] = msgHooks[0]
	msgFields := schema.Msg{}.Fields()
	_ = msgFields
	// msgDescUUID is the schema descriptor for UUID field.
	msgDescUUID := msgFields[0].Descriptor()
	// msg.DefaultUUID holds the default value on creation for the UUID field.
	msg.DefaultUUID = msgDescUUID.Default.(func() uuid.UUID)
	// msgDescSourceMsgSender is the schema descriptor for SourceMsgSender field.
	msgDescSourceMsgSender := msgFields[2].Descriptor()
	// msg.SourceMsgSenderValidator is a validator for the "SourceMsgSender" field. It is called by the builders before save.
	msg.SourceMsgSenderValidator = msgDescSourceMsgSender.Validators[0].(func([]byte) error)
	// msgDescDestAddress is the schema descriptor for DestAddress field.
	msgDescDestAddress := msgFields[3].Descriptor()
	// msg.DestAddressValidator is a validator for the "DestAddress" field. It is called by the builders before save.
	msg.DestAddressValidator = msgDescDestAddress.Validators[0].(func([]byte) error)
	// msgDescTxHash is the schema descriptor for TxHash field.
	msgDescTxHash := msgFields[9].Descriptor()
	// msg.TxHashValidator is a validator for the "TxHash" field. It is called by the builders before save.
	msg.TxHashValidator = msgDescTxHash.Validators[0].(func([]byte) error)
	// msgDescCreatedAt is the schema descriptor for CreatedAt field.
	msgDescCreatedAt := msgFields[10].Descriptor()
	// msg.DefaultCreatedAt holds the default value on creation for the CreatedAt field.
	msg.DefaultCreatedAt = msgDescCreatedAt.Default.(time.Time)
	receiptHooks := schema.Receipt{}.Hooks()
	receipt.Hooks[0] = receiptHooks[0]
	receiptFields := schema.Receipt{}.Fields()
	_ = receiptFields
	// receiptDescUUID is the schema descriptor for UUID field.
	receiptDescUUID := receiptFields[0].Descriptor()
	// receipt.DefaultUUID holds the default value on creation for the UUID field.
	receipt.DefaultUUID = receiptDescUUID.Default.(func() uuid.UUID)
	// receiptDescRelayerAddress is the schema descriptor for RelayerAddress field.
	receiptDescRelayerAddress := receiptFields[4].Descriptor()
	// receipt.RelayerAddressValidator is a validator for the "RelayerAddress" field. It is called by the builders before save.
	receipt.RelayerAddressValidator = receiptDescRelayerAddress.Validators[0].(func([]byte) error)
	// receiptDescTxHash is the schema descriptor for TxHash field.
	receiptDescTxHash := receiptFields[8].Descriptor()
	// receipt.TxHashValidator is a validator for the "TxHash" field. It is called by the builders before save.
	receipt.TxHashValidator = receiptDescTxHash.Validators[0].(func([]byte) error)
	// receiptDescCreatedAt is the schema descriptor for CreatedAt field.
	receiptDescCreatedAt := receiptFields[9].Descriptor()
	// receipt.DefaultCreatedAt holds the default value on creation for the CreatedAt field.
	receipt.DefaultCreatedAt = receiptDescCreatedAt.Default.(time.Time)
	xprovidercursorFields := schema.XProviderCursor{}.Fields()
	_ = xprovidercursorFields
	// xprovidercursorDescUUID is the schema descriptor for UUID field.
	xprovidercursorDescUUID := xprovidercursorFields[0].Descriptor()
	// xprovidercursor.DefaultUUID holds the default value on creation for the UUID field.
	xprovidercursor.DefaultUUID = xprovidercursorDescUUID.Default.(func() uuid.UUID)
	// xprovidercursorDescCreatedAt is the schema descriptor for CreatedAt field.
	xprovidercursorDescCreatedAt := xprovidercursorFields[3].Descriptor()
	// xprovidercursor.DefaultCreatedAt holds the default value on creation for the CreatedAt field.
	xprovidercursor.DefaultCreatedAt = xprovidercursorDescCreatedAt.Default.(time.Time)
	// xprovidercursorDescUpdatedAt is the schema descriptor for UpdatedAt field.
	xprovidercursorDescUpdatedAt := xprovidercursorFields[4].Descriptor()
	// xprovidercursor.DefaultUpdatedAt holds the default value on creation for the UpdatedAt field.
	xprovidercursor.DefaultUpdatedAt = xprovidercursorDescUpdatedAt.Default.(time.Time)
}

const (
	Version = "v0.13.1"                                         // Version of ent codegen.
	Sum     = "h1:uD8QwN1h6SNphdCCzmkMN3feSUzNnVvV/WIkHKMbzOE=" // Sum of ent codegen.
)
