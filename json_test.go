package variable_json

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	valResponse = `{}`
)

func TestParse(t *testing.T) {
	result, err := Parse([]byte(valResponse))
	require.Nil(t, err)

	data, err := json.Marshal(result)
	if err != nil {
		log.Printf("err %+v", err)
	}
	require.Nil(t, err)

	log.Printf("data = %s", string(data))
}
