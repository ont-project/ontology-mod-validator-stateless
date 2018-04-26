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

	ec "github.com/ont-project/ontology-framework/common/config"
	"github.com/ont-project/ontology-framework/common/di"
	"github.com/ont-project/ontology-framework/core"
	"github.com/ont-project/ontology-framework/core/channel"
	"github.com/ont-project/ontology-framework/core/meta"
	"github.com/ont-project/ontology/runtime"
)

type StatelessValidatorModule struct {
	runtime.RuntimeModule
}

type StatelessValidatorConfig struct {
	ec.Config
	channel.ChannelSetConfig
}

// 0x8000 for stateless validator
var validatorExtra ec.Config = ec.Config{Magic: core.MagicValidator, Version: 0x8001}

func init() {
	di.ConfigMap.Register(ec.GetExtraId(validatorExtra), &StatelessValidatorConfig{})
	di.EngineMap.Register(ec.GetExtraId(validatorExtra), &StatelessValidatorModule{})
}

func (self *StatelessValidatorModule) NewStub() (channel.Stub, error) {
	return &core.ValidatorStub{}, nil
}

func (self *StatelessValidatorModule) ExecuteGo(stub channel.Stub, cmd channel.Command, param ...interface{}) (interface{}, error) {
	switch cmd {
	case core.Verify:
		return nil, self.Verify(param[0].(core.VerifyType), param[1:]...)
	default:
		return nil, errors.New("unknown command")
	}
}

func (self *StatelessValidatorModule) Verify(verifyType core.VerifyType, data ...interface{}) error {
	switch verifyType {
	case core.Transaction:
		if err := VerifyTransaction(data[0].(*meta.Transaction)); err != nil {
			return nil
		}
		return errors.New("verify error, transaction")
	case core.Block:
		fallthrough
	case core.Signature:
		fallthrough
	case core.Program:
		fallthrough
	default:
		return errors.New("verify error, verify type mismatch")
	}
}
