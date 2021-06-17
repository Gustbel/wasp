// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package wasmproc

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

type ScMaps struct {
	ScSandboxObject
}

func NewScMaps(vm *wasmProcessor) *ScMaps {
	a := &ScMaps{}
	a.vm = vm
	return a
}

func (a *ScMaps) GetObjectID(keyId int32, typeId int32) int32 {
	return GetArrayObjectID(a, keyId, typeId, func() WaspObject {
		return NewScDict(a.vm)
	})
}
