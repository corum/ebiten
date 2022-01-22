// Copyright 2022 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gamepad

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")

	getModuleHandleWProc = kernel32.NewProc("GetModuleHandleW")
)

type nativeGamepads struct {
	dinput8    windows.Handle
	dinput8API uintptr // IDirectInput8W*
	xinput     windows.Handle

	directInput8Create uintptr
}

type nativeGamepad struct {
}

var iidIDirectInput8W = windows.GUID{
	Data1: 0xbf798031,
	Data2: 0x483a,
	Data3: 0x4da2,
	Data4: [...]byte{0xaa, 0x99, 0x5d, 0x64, 0xed, 0x36, 0x97, 0x00},
}

func (g *nativeGamepads) init() error {
	if h, err := windows.LoadLibrary("dinput8.dll"); err == nil {
		g.dinput8 = h

		p, err := windows.GetProcAddress(h, "DirectInput8Create")
		if err != nil {
			return err
		}
		g.directInput8Create = p
	}

	for _, dll := range []string{
		"xinput1_4.dll",
		"xinput1_3.dll",
		"xinput9_1_0.dll",
		"xinput1_2.dll",
		"xinput1_1.dll",
	} {
		if h, err := windows.LoadLibrary(dll); err == nil {
			g.xinput = h
			break
		}
	}

	if g.dinput8 != 0 {
		m, _, err := getModuleHandleWProc.Call(0)
		if err != nil && err != windows.ERROR_SUCCESS {
			return fmt.Errorf("gamepad: GetModuleHandleW failed: %w", err)
		}
		if m == 0 {
			return fmt.Errorf("gamepad: GetModuleHandleW returned 0")
		}

		const directInputVersion = 0x0800

		r, _, err := syscall.Syscall6(g.directInput8Create, 5,
			m, directInputVersion, uintptr(unsafe.Pointer(&iidIDirectInput8W)), uintptr(unsafe.Pointer(&g.dinput8API)), 0,
			0)
		if err != nil && err != windows.ERROR_SUCCESS {
			// This returns windows.ERROR_INSUFFICIENT_BUFFER. Why?
			return fmt.Errorf("gamepad: DirectInput8Create failed: %w", err)
		}
		if r != 0 {
			return fmt.Errorf("gamepad: DirectInput8Create returned %d", r)
		}
		println(g.dinput8API)
	}

	return nil
}

func (*nativeGamepad) present() bool {
	return false
}

func (*nativeGamepad) update() {
}

func (*nativeGamepad) axisNum() int {
	return 0
}

func (*nativeGamepad) buttonNum() int {
	return 0
}

func (*nativeGamepad) hatNum() int {
	return 0
}

func (*nativeGamepad) axisValue(axis int) float64 {
	return 0
}

func (*nativeGamepad) isButtonPressed(button int) bool {
	return false
}

func (*nativeGamepad) hatState(hat int) int {
	return 0
}
