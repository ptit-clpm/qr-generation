package utils

import (
	"fmt"
	"testing"
)

func TestCRC16_Standard(t *testing.T) {
	// CRC-16/CCITT-FALSE: poly 0x1021, init 0xFFFF, no final XOR
	// Known test: "123456789" => 0x29B1
	result := crc16("123456789")
	if result != "29B1" {
		t.Errorf("CRC16('123456789') = %s, want 29B1", result)
	}
}

func TestVietQRContent(t *testing.T) {
	content := VietQRContent("MB", "0924517780", "NGUYEN DUY QUANG NHAT", "QRPRO-1-ABCD1234", 99000)
	fmt.Println("VietQR content:", content)
}
