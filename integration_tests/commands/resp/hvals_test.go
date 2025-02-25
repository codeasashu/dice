// This file is part of DiceDB.
// Copyright (C) 2024 DiceDB (dicedb.io).
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package resp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHVals(t *testing.T) {
	conn := getLocalConnection()
	defer conn.Close()

	testCases := []TestCase{
		{
			name:     "RESP HVALS with multiple fields",
			commands: []string{"HSET hvalsKey field value", "HSET hvalsKey field2 value1", "HVALS hvalsKey"},
			expected: []interface{}{int64(1), int64(1), []interface{}{"value", "value1"}},
		},
		{
			name:     "RESP HVALS with non-existing key",
			commands: []string{"HVALS hvalsKey01"},
			expected: []interface{}{[]any{}},
		},
		{
			name:     "HVALS on wrong key type",
			commands: []string{"SET hvalsKey02 field", "HVALS hvalsKey02"},
			expected: []interface{}{"OK", "ERR -WRONGTYPE Operation against a key holding the wrong kind of value"},
		},
		{
			name:     "HVALS with wrong number of arguments",
			commands: []string{"HVALS hvalsKey03 x", "HVALS"},
			expected: []interface{}{"ERR wrong number of arguments for 'hvals' command", "ERR wrong number of arguments for 'hvals' command"},
		},
		{
			name:     "RESP One or more vals exist",
			commands: []string{"HSET key field value", "HSET key field1 value1", "HVALS key"},
			expected: []interface{}{int64(1), int64(1), []interface{}{"value", "value1"}},
		},
		{
			name:     "RESP No values exist",
			commands: []string{"HVALS key"},
			expected: []interface{}{[]interface{}{}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			FireCommand(conn, "HDEL key field")
			FireCommand(conn, "HDEL key field1")

			for i, cmd := range tc.commands {
				result := FireCommand(conn, cmd)
				switch e := tc.expected[i].(type) {
				case []interface{}:
					assert.ElementsMatch(t, e, tc.expected[i])
				default:
					assert.Equal(t, tc.expected[i], result)
				}
			}
		})
	}

	FireCommand(conn, "HDEL key field")
	FireCommand(conn, "HDEL key field1")
}
