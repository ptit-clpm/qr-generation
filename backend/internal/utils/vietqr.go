package utils

import (
	"fmt"
)

func VietQRContent(bankCode, accountNo, accountName, content string, amount float64) string {
	p := newPayload()

	p.Append("00", "01")
	p.Append("01", "12")

	cai := newPayload()
	cai.Append("00", bankCode)
	cai.Append("01", accountNo)

	ma := newPayload()
	ma.Append("00", "A000000775")
	ma.Append("01", cai.String())
	p.Append("38", ma.String())

	p.Append("53", "704")

	amt := fmt.Sprintf("%.0f", amount)
	p.Append("54", amt)

	p.Append("58", "VN")

	p.Append("59", accountName)

	p.Append("60", "HANOI")

	ad := newPayload()
	ad.Append("01", content)
	p.Append("62", ad.String())

	raw := p.String()
	crc := crc16(raw)
	raw += "6304" + crc

	return raw
}

type payloadBuilder struct{ s string }

func newPayload() *payloadBuilder { return &payloadBuilder{} }

func (p *payloadBuilder) Append(id, value string) {
	length := fmt.Sprintf("%02d", len(value))
	p.s += id + length + value
}

func (p *payloadBuilder) String() string { return p.s }

func crc16(data string) string {
	crc := uint16(0xFFFF)
	for i := 0; i < len(data); i++ {
		crc ^= uint16(data[i]) << 8
		for j := 0; j < 8; j++ {
			if crc&0x8000 != 0 {
				crc = (crc << 1) ^ 0x1021
			} else {
				crc <<= 1
			}
		}
	}
	return fmt.Sprintf("%04X", crc)
}
