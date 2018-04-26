/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */

package stateless

import (
	"errors"
	"fmt"

	"github.com/ont-project/ontology-framework/common/log"
	"github.com/ont-project/ontology-framework/core/meta"
	"github.com/ont-project/ontology-framework/core/meta/helper"
	"github.com/ont-project/ontology-framework/core/meta/payload"
	"github.com/ont-project/ontology-framework/core/meta/types"
	"github.com/ont-project/ontology-framework/util"
)

// VerifyTransaction verifys received single transaction
func VerifyTransaction(tx *meta.Transaction) error {
	if err := checkTransactionSignatures(tx); err != nil {
		log.Info("transaction verify error:", err)
		return errors.New("[validator error] transaction contract")
	}

	if err := checkTransactionPayload(tx); err != nil {
		log.Warn("[VerifyTransaction],", err)
		return errors.New("[validator error] transaction payload")
	}

	return nil
}

func checkTransactionSignatures(tx *meta.Transaction) error {
	hash := tx.Hash()
	address := make(map[types.Address]bool, len(tx.Sigs))
	for _, sig := range tx.Sigs {
		m := int(sig.M)
		kn := len(sig.PubKeys)
		sn := len(sig.SigData)

		if kn > 24 || sn < m || m > kn {
			return errors.New("wrong tx sig param length")
		}

		if kn == 1 {
			err := helper.Verify(sig.PubKeys[0], hash[:], sig.SigData[0])
			if err != nil {
				return errors.New("signature verification failed")
			}

			address[types.AddressFromPubKey(sig.PubKeys[0])] = true
		} else {
			if err := helper.VerifyMultiSignature(hash[:], sig.PubKeys, m, sig.SigData); err != nil {
				return err
			}

			addr, _ := types.AddressFromMultiPubKeys(sig.PubKeys, m)
			address[addr] = true
		}
	}

	// check all payers in address
	for _, fee := range tx.Fee {
		if address[fee.Payer] == false {
			return errors.New("signature missing for payer: " + util.ToHexString(fee.Payer[:]))
		}
	}

	return nil
}

func checkTransactionPayload(tx *meta.Transaction) error {

	switch pld := tx.Payload.(type) {
	case *payload.DeployCode:
		return nil
	case *payload.InvokeCode:
		return nil
	case *payload.Bookkeeping:
		return nil
	default:
		return errors.New(fmt.Sprint("[txValidator], unimplemented transaction payload type.", pld))
	}
	return nil
}
