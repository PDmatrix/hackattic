package help_me_unpack

import (
	b64 "encoding/base64"
	"encoding/binary"
	"encoding/json"
	"math"
)

type HelpMeUnpack struct{}

type Data struct {
	Bytes string `json:"bytes"`
}

type Output struct {
	Int             int32   `json:"int"`
	Uint            uint32  `json:"uint"`
	Short           int16   `json:"short"`
	Float           float32 `json:"float"`
	Double          float64 `json:"double"`
	BigEndianDouble float64 `json:"big_endian_double"`
}

func (d HelpMeUnpack) Solve(input string) (interface{}, error) {
	data := new(Data)
	output := new(Output)
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return nil, err
	}
	sDec, _ := b64.StdEncoding.DecodeString(data.Bytes)

	output.Int = int32(binary.LittleEndian.Uint32(sDec[0:4]))
	output.Uint = binary.LittleEndian.Uint32(sDec[4:8])
	output.Short = int16(binary.LittleEndian.Uint16(sDec[8:10]))
	output.Float = math.Float32frombits(binary.LittleEndian.Uint32(sDec[12:16]))
	output.Double = math.Float64frombits(binary.LittleEndian.Uint64(sDec[16:24]))
	output.BigEndianDouble = math.Float64frombits(binary.BigEndian.Uint64(sDec[24:32]))

	return output, nil
}
