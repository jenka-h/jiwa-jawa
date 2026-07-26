package err

import (
	"errors"
	"fmt"
)

// PrintError prints an error message with [Error] prefix
func PrintError(msg string) {
	fmt.Printf("[Error] %s\n", msg)
}

// PrintUsageError prints a usage error message
func PrintUsageError(msg string) {
	fmt.Printf("[Error] %s\n", msg)
}

// PrintErrorf prints a formatted error message
func PrintErrorf(format string, args ...interface{}) {
	fmt.Printf("[Error] "+format+"\n", args...)
}

/*
pakai method errors.Is untuk pengecekan spesifik nantinya
*/
var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrPacketTooShort     = errors.New("packet too short")
	ErrInvalidMessageType = errors.New("invalid message type")
	ErrInvalidACK         = errors.New("invalid ACK packet")
)
