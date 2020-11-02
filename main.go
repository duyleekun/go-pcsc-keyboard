package main

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/sf1/go-card/smartcard"
)

var (
	user32        = syscall.NewLazyDLL("user32.dll")
	sendInputProc = user32.NewProc("SendInput")
)

func sendKey(s uint8) {
	type keyboardInput struct {
		wVk         uint16
		wScan       uint16
		dwFlags     uint32
		time        uint32
		dwExtraInfo uint64
	}

	type input struct {
		inputType uint32
		ki        keyboardInput
		padding   uint64
	}

	var i input
	i.inputType = 1 //INPUT_KEYBOARD

	i.ki.wVk = uint16(s)
	i.ki.dwFlags = 0
	sendInputProc.Call(
		uintptr(1),
		uintptr(unsafe.Pointer(&i)),
		uintptr(unsafe.Sizeof(i)),
	)
	i.ki.dwFlags = 2
	sendInputProc.Call(
		uintptr(1),
		uintptr(unsafe.Pointer(&i)),
		uintptr(unsafe.Sizeof(i)),
	)
	// log.Printf("ret: %v error: %v", ret, err)
}

func main() {
	done := false
	ctx, err := smartcard.EstablishContext()
	defer println(err)
	defer ctx.Release()
	for !done {
		fmt.Printf("Please insert card")
		reader, _ := ctx.WaitForCardPresent()
		// handle error, if any

		card, err := reader.Connect()
		fmt.Printf("%%", err)

		command := smartcard.Command2(0xff, 0xca, 0x00, 0x00, 0x00)
		response, _ := card.TransmitAPDU(command)
		println(err)

		responseData := response.Data()

		//Reverse byte if needed
		// for i, j := 0, len(responseData)-1; i < j; i, j = i+1, j-1 {
		// 	responseData[i], responseData[j] = responseData[j], responseData[i]
		// }

		// handle error, if any
		fmt.Printf("Response: %X %d\n", responseData, len(responseData))
		str := fmt.Sprintf("%X\n", responseData)

		for i := 0; i < len(str); i++ {
			sendKey(str[i])
		}
		sendKey(0x0D)

		//sendKey("A")
		card.Disconnect()
		fmt.Printf("Please remove card\n")
		reader.WaitUntilCardRemoved()
	}
	return
}
