package utils

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

type testItem struct {
	Name string  `json:"name"`
	Age  int     `json:"age"`
	Pay  float64 `json:"pay"`
}

func Test_ShouldMatchPredicate(t *testing.T) {
	item := &testItem{Name: "john", Age: 21, Pay: 21}
	b, err := json.Marshal(item)
	require.NoError(t, err)
	require.True(t, MatchPredicate(nil, nil))
	require.False(t, MatchPredicate(nil, map[string]string{"name": "joe"}))
	require.False(t, MatchPredicate(b, map[string]string{"name": "joe"}))
	require.True(t, MatchPredicate(b, map[string]string{"name": "john"}))
	require.True(t, MatchPredicate(b, map[string]string{"name:==": "john"}))
	require.True(t, MatchPredicate(b, map[string]string{"name:!=": "joe"}))
	require.True(t, MatchPredicate(b, map[string]string{"name:>=": "alice"}))
	require.True(t, MatchPredicate(b, map[string]string{"name:<=": "smith"}))
	require.True(t, MatchPredicate(b, map[string]string{"name:>": "alice"}))
	require.True(t, MatchPredicate(b, map[string]string{"name:<": "smith"}))

	require.True(t, MatchPredicate(b, map[string]string{"age:==": "21"}))
	require.True(t, MatchPredicate(b, map[string]string{"age:!=": "18"}))
	require.True(t, MatchPredicate(b, map[string]string{"age:>=": "18"}))
	require.True(t, MatchPredicate(b, map[string]string{"age:<=": "55"}))
	require.True(t, MatchPredicate(b, map[string]string{"age:>": "10"}))
	require.True(t, MatchPredicate(b, map[string]string{"age:<": "100"}))

	require.True(t, MatchPredicate(b, map[string]string{"pay:==": "21"}))
	require.True(t, MatchPredicate(b, map[string]string{"pay:!=": "18"}))
	require.True(t, MatchPredicate(b, map[string]string{"pay:>=": "18"}))
	require.True(t, MatchPredicate(b, map[string]string{"pay:<=": "55"}))
	require.True(t, MatchPredicate(b, map[string]string{"pay:>": "10"}))
	require.True(t, MatchPredicate(b, map[string]string{"pay:<": "100"}))
}
