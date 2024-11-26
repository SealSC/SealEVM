package evmInt256

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

func (i *Int) MarshalJSON() ([]byte, error) {
	hexStr := "0x" + i.Text(16)
	return json.Marshal(hexStr)
}

func (i *Int) UnmarshalJSON(data []byte) error {
	if i.Int == nil {
		i.Int = new(big.Int)
	}

	var raw json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	str := string(raw)
	if len(str) > 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	if strings.HasPrefix(str, "0x") {
		str = strings.TrimPrefix(str, "0x")
		if _, ok := i.SetString(str, 16); !ok {
			return fmt.Errorf("invalid hex string: %s", str)
		}
	} else if val, err := strconv.ParseInt(str, 10, 64); err == nil {
		i.SetInt64(val)
	} else if _, ok := i.SetString(str, 10); !ok {
		return fmt.Errorf("invalid decimal string: %s", str)
	}

	return nil
}
