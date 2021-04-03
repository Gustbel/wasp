// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package chain

import (
	"github.com/iotaledger/goshimmer/packages/ledgerstate"
	"golang.org/x/xerrors"
	"io"

	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/util"
)

func readOutputID(r io.Reader, oid *ledgerstate.OutputID) error {
	n, err := r.Read(oid[:])
	if err != nil {
		return xerrors.Errorf("failed parse OutputID: %w", err)
	}
	if n != ledgerstate.OutputIDLength {
		return xerrors.Errorf("failed parse OutputID: wrong data length")
	}
	return nil
}

func (msg *NotifyReqMsg) Write(w io.Writer) error {
	if _, err := w.Write(msg.StateOutputID.Bytes()); err != nil {
		return err
	}
	if err := util.WriteUint16(w, uint16(len(msg.RequestIDs))); err != nil {
		return err
	}
	for _, reqid := range msg.RequestIDs {
		if _, err := w.Write(reqid[:]); err != nil {
			return err
		}
	}
	return nil
}

func (msg *NotifyReqMsg) Read(r io.Reader) error {
	if err := readOutputID(r, &msg.StateOutputID); err != nil {
		return err
	}
	var arrLen uint16
	err := util.ReadUint16(r, &arrLen)
	if err != nil {
		return err
	}
	if arrLen == 0 {
		return nil
	}
	msg.RequestIDs = make([]coretypes.RequestID, arrLen)
	for i := range msg.RequestIDs {
		_, err = r.Read(msg.RequestIDs[i][:])
		if err != nil {
			return err
		}
	}
	return nil
}

func (msg *NotifyFinalResultPostedMsg) Write(w io.Writer) error {
	if _, err := w.Write(msg.StateOutputID.Bytes()); err != nil {
		return err
	}
	if _, err := w.Write(msg.TxId.Bytes()); err != nil {
		return err
	}
	return nil
}

func (msg *NotifyFinalResultPostedMsg) Read(r io.Reader) error {
	err := readOutputID(r, &msg.StateOutputID)
	if err != nil {
		return err
	}
	if n, err := r.Read(msg.TxId[:]); err != nil || n != ledgerstate.TransactionIDLength {
		return xerrors.Errorf("failed to read transaction id: err=%v", err)
	}
	return nil
}

func (msg *StartProcessingBatchMsg) Write(w io.Writer) error {
	if _, err := w.Write(msg.StateOutputID.Bytes()); err != nil {
		return err
	}
	if err := util.WriteUint16(w, uint16(len(msg.RequestIDs))); err != nil {
		return err
	}
	for i := range msg.RequestIDs {
		if _, err := w.Write(msg.RequestIDs[i].Bytes()); err != nil {
			return err
		}
	}
	if err := msg.FeeDestination.Write(w); err != nil {
		return err
	}
	return nil
}

func (msg *StartProcessingBatchMsg) Read(r io.Reader) error {
	if err := readOutputID(r, &msg.StateOutputID); err != nil {
		return err
	}
	var size uint16
	if err := util.ReadUint16(r, &size); err != nil {
		return err
	}
	msg.RequestIDs = make([]coretypes.RequestID, size)
	for i := range msg.RequestIDs {
		if n, err := r.Read(msg.RequestIDs[i][:]); err != nil || n != ledgerstate.OutputIDLength {
			return err
		}
	}
	if err := msg.FeeDestination.Read(r); err != nil {
		return err
	}
	return nil
}

func (msg *SignedHashMsg) Write(w io.Writer) error {
	if _, err := w.Write(msg.StateOutputID.Bytes()); err != nil {
		return err
	}
	if err := util.WriteUint64(w, uint64(msg.OrigTimestamp)); err != nil {
		return err
	}
	if _, err := w.Write(msg.BatchHash[:]); err != nil {
		return err
	}
	if _, err := w.Write(msg.EssenceHash[:]); err != nil {
		return err
	}
	if err := util.WriteBytes16(w, msg.SigShare); err != nil {
		return err
	}
	return nil
}

func (msg *SignedHashMsg) Read(r io.Reader) error {
	if err := readOutputID(r, &msg.StateOutputID); err != nil {
		return err
	}
	var ts uint64
	if err := util.ReadUint64(r, &ts); err != nil {
		return err
	}
	msg.OrigTimestamp = int64(ts)

	if err := util.ReadHashValue(r, &msg.BatchHash); err != nil {
		return err
	}
	if err := util.ReadHashValue(r, &msg.EssenceHash); err != nil {
		return err
	}
	var err error
	if msg.SigShare, err = util.ReadBytes16(r); err != nil {
		return err
	}
	return nil
}

func (msg *GetBlockMsg) Write(w io.Writer) error {
	if err := util.WriteUint32(w, msg.BlockIndex); err != nil {
		return err
	}
	return nil
}

func (msg *GetBlockMsg) Read(r io.Reader) error {
	if err := util.ReadUint32(r, &msg.BlockIndex); err != nil {
		return err
	}
	return nil
}

func (msg *BlockHeaderMsg) Write(w io.Writer) error {
	if err := util.WriteUint32(w, msg.BlockIndex); err != nil {
		return err
	}
	if err := util.WriteUint16(w, msg.Size); err != nil {
		return err
	}
	if _, err := w.Write(msg.ApprovingOutputID.Bytes()); err != nil {
		return err
	}
	return nil
}

func (msg *BlockHeaderMsg) Read(r io.Reader) error {
	if err := util.ReadUint32(r, &msg.BlockIndex); err != nil {
		return err
	}
	if err := util.ReadUint16(r, &msg.Size); err != nil {
		return err
	}
	if err := readOutputID(r, &msg.ApprovingOutputID); err != nil {
		return err
	}
	return nil
}

func (msg *StateUpdateMsg) Write(w io.Writer) error {
	if err := util.WriteUint32(w, msg.BlockIndex); err != nil {
		return err
	}
	if err := msg.StateUpdate.Write(w); err != nil {
		return err
	}
	if err := util.WriteUint16(w, msg.IndexInTheBlock); err != nil {
		return err
	}
	return nil
}

func (msg *StateUpdateMsg) Read(r io.Reader) error {
	if err := util.ReadUint32(r, &msg.BlockIndex); err != nil {
		return err
	}
	msg.StateUpdate = state.NewStateUpdate(coretypes.RequestID{})
	if err := msg.StateUpdate.Read(r); err != nil {
		return err
	}
	if err := util.ReadUint16(r, &msg.IndexInTheBlock); err != nil {
		return err
	}
	return nil
}

func (msg *BlockIndexPingPongMsg) Write(w io.Writer) error {
	if err := util.WriteUint32(w, msg.BlockIndex); err != nil {
		return err
	}
	return util.WriteBoolByte(w, msg.RSVP)
}

func (msg *BlockIndexPingPongMsg) Read(r io.Reader) error {
	if err := util.ReadUint32(r, &msg.BlockIndex); err != nil {
		return err
	}
	return util.ReadBoolByte(r, &msg.RSVP)
}
