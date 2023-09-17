package domain

import (
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ShouldParsePredicateFalse(t *testing.T) {
	// GIVEN a template string
	b := `{{false}}`
	// WHEN parsing template
	out, err := ParseTemplate(b, &PrincipalExt{Delegate: &types.Principal{}}, &types.Resource{}, &services.AuthRequest{})
	// THEN it should not fail
	require.NoError(t, err)
	require.Equal(t, "false", string(out))
}

func Test_ShouldParsePredicateTrue(t *testing.T) {
	// GIVEN a template string
	b := `{{false}}`
	// WHEN parsing template
	out, err := ParseTemplate(b, &PrincipalExt{Delegate: &types.Principal{}}, &types.Resource{}, &services.AuthRequest{})
	// THEN it should not fail
	require.NoError(t, err)
	require.Equal(t, "false", string(out))
}

func Test_ShouldCompareDateFormat(t *testing.T) {
	// GIVEN a template string
	b := `{{$year := TimeNow "2006"}}{{le "2023" $year}}  `
	// WHEN parsing name
	out, err := ParseTemplate(b, &PrincipalExt{Delegate: &types.Principal{}}, &types.Resource{}, &services.AuthRequest{})
	// THEN it should succeed
	require.NoError(t, err)
	require.Equal(t, "true", string(out))
}

func Test_ShouldParseDateFormat(t *testing.T) {
	// GIVEN a template string
	b := `{{TimeNow "2006"}}`
	// WHEN parsing name
	_, err := ParseTemplate(b, &PrincipalExt{Delegate: &types.Principal{}}, &types.Resource{}, &services.AuthRequest{})
	// THEN it should succeed
	require.NoError(t, err)
}

func Test_ShouldParseTimeNow(t *testing.T) {
	// GIVEN a template string
	b := `{{TimeNow "mm"}}`
	// WHEN parsing name
	_, err := ParseTemplate(b, &PrincipalExt{Delegate: &types.Principal{}}, &types.Resource{}, &services.AuthRequest{})
	// THEN it should succeed
	require.NoError(t, err)
}

func Test_ShouldParseInt(t *testing.T) {
	// GIVEN a template string
	b := `{{Int 3}}`
	// WHEN parsing int
	_, err := ParseTemplate(b, &PrincipalExt{Delegate: &types.Principal{}}, &types.Resource{}, &services.AuthRequest{})
	// THEN it should succeed
	require.NoError(t, err)
}

func Test_ShouldParseFloat(t *testing.T) {
	// GIVEN a template string
	b := `{{Float 3}}`
	// WHEN parsing int
	out, err := ParseTemplate(b, &PrincipalExt{Delegate: &types.Principal{}}, &types.Resource{}, &services.AuthRequest{})
	// THEN it should succeed
	require.NoError(t, err)
	require.Equal(t, "3", string(out))
}

func Test_ShouldParseLT(t *testing.T) {
	// GIVEN a template string
	b := `{{LT 3 5}}`
	// WHEN parsing int
	out, err := ParseTemplate(b, &PrincipalExt{Delegate: &types.Principal{}}, &types.Resource{}, &services.AuthRequest{})
	// THEN it should succeed
	require.NoError(t, err)
	require.Equal(t, "true", string(out))
}

func Test_ShouldParseLE(t *testing.T) {
	// GIVEN a template string
	b := `{{LE 3 5}}`
	// WHEN parsing int
	out, err := ParseTemplate(b, &PrincipalExt{Delegate: &types.Principal{}}, &types.Resource{}, &services.AuthRequest{})
	// THEN it should succeed
	require.NoError(t, err)
	require.Equal(t, "true", string(out))
}

func Test_ShouldParseEQ(t *testing.T) {
	// GIVEN a template string
	b := `{{EQ 3 5}}`
	// WHEN parsing int
	out, err := ParseTemplate(b, &PrincipalExt{Delegate: &types.Principal{}}, &types.Resource{}, &services.AuthRequest{})
	// THEN it should succeed
	require.NoError(t, err)
	require.Equal(t, "false", string(out))
}

func Test_ShouldParseGT(t *testing.T) {
	// GIVEN a template string
	b := `{{GT 3 5}}`
	// WHEN parsing int
	out, err := ParseTemplate(b, &PrincipalExt{Delegate: &types.Principal{}}, &types.Resource{}, &services.AuthRequest{})
	// THEN it should succeed
	require.NoError(t, err)
	require.Equal(t, "false", string(out))
}

func Test_ShouldParseGE(t *testing.T) {
	// GIVEN a template string
	b := `{{GE 3 5}}`
	// WHEN parsing int
	out, err := ParseTemplate(b, &PrincipalExt{Delegate: &types.Principal{}}, &types.Resource{}, &services.AuthRequest{})
	// THEN it should succeed
	require.NoError(t, err)
	require.Equal(t, "false", string(out))
}

func Test_ShouldFailOnNoVariables(t *testing.T) {
	// GIVEN a template string
	b := `
{{with .Account -}}
Account: {{.}}
{{- end}}
Money: {{.Money}}
{{if .Note -}}
Note: {{.Note}}
{{- end}}
`

	// WHEN parsing template
	_, err := ParseTemplate(b, &PrincipalExt{Delegate: &types.Principal{}}, &types.Resource{}, &services.AuthRequest{})

	// THEN it should not fail without params
	require.NoError(t, err)

	// AND it should not fail with params
	_, err = ParseTemplate(b, &PrincipalExt{Delegate: &types.Principal{}}, &types.Resource{}, &services.AuthRequest{
		Context: map[string]string{"Account": "x123", "Money": "12", "Note": "ty"},
	})
	require.NoError(t, err)
}

func Test_ShouldParseBeginsWith(t *testing.T) {
	// GIVEN a template string
	b := `{{BeginsWith "test1" "test"}}`
	// WHEN parsing int
	out, err := ParseTemplate(b, &PrincipalExt{Delegate: &types.Principal{}}, &types.Resource{}, &services.AuthRequest{})
	// THEN it should succeed
	require.NoError(t, err)
	require.Equal(t, "true", string(out))
}

func Test_ShouldParseEndsWith(t *testing.T) {
	// GIVEN a template string
	b := `{{EndsWith "test1" "est1"}}`
	// WHEN parsing int
	out, err := ParseTemplate(b, &PrincipalExt{Delegate: &types.Principal{}}, &types.Resource{}, &services.AuthRequest{})
	// THEN it should succeed
	require.NoError(t, err)
	require.Equal(t, "true", string(out))
}

func Test_ShouldParseContains(t *testing.T) {
	// GIVEN a template string
	b := `{{Contains "test1" "est1"}}`
	// WHEN parsing int
	out, err := ParseTemplate(b, &PrincipalExt{Delegate: &types.Principal{}}, &types.Resource{}, &services.AuthRequest{})
	// THEN it should succeed
	require.NoError(t, err)
	require.Equal(t, "true", string(out))
}

func Test_ShouldParseIncludes(t *testing.T) {
	// GIVEN a template string
	b := `{{Includes "test1 test2" "test1"}}`
	// WHEN parsing int
	out, err := ParseTemplate(b, &PrincipalExt{Delegate: &types.Principal{}}, &types.Resource{}, &services.AuthRequest{})
	// THEN it should succeed
	require.NoError(t, err)
	require.Equal(t, "true", string(out))
}

func Test_ShouldParseStringToFloatArray(t *testing.T) {
	// GIVEN a template string
	b := `{{StringToFloatArray "1 3.5, 12.3: 34.123"}}`
	// WHEN parsing int
	out, err := ParseTemplate(b, &PrincipalExt{Delegate: &types.Principal{}}, &types.Resource{}, &services.AuthRequest{})
	// THEN it should succeed
	require.NoError(t, err)
	require.Equal(t, "[1 3.5 12.3 34.123]", string(out))
}

func Test_ShouldParseDistanceWithinKM(t *testing.T) {
	// GIVEN a template string
	b := `{{DistanceWithinKM "47.620422,-122.349358" "46.879967,-121.726906" 100}}`
	// WHEN parsing int
	out, err := ParseTemplate(b, &PrincipalExt{Delegate: &types.Principal{}}, &types.Resource{}, &services.AuthRequest{})
	// THEN it should succeed
	require.NoError(t, err)
	require.Equal(t, "true", string(out))
}

func Test_ShouldParseTimeInRange(t *testing.T) {
	// GIVEN a template string
	b := `{{TimeInRange "11:00am" "10:00am" "8:00pm"}}`
	// WHEN parsing int
	out, err := ParseTemplate(b, &PrincipalExt{Delegate: &types.Principal{}}, &types.Resource{}, &services.AuthRequest{})
	// THEN it should succeed
	require.NoError(t, err)
	require.Equal(t, "true", string(out))
}

func Test_ShouldParseIPAddrInRange(t *testing.T) {
	// GIVEN a template string
	b := `{{IPInRange "211.211.211.5" "211.211.211.0/24"}}`
	// WHEN parsing int
	out, err := ParseTemplate(b, &PrincipalExt{Delegate: &types.Principal{}}, &types.Resource{}, &services.AuthRequest{})
	// THEN it should succeed
	require.NoError(t, err)
	require.Equal(t, "true", string(out))
}

func Test_ShouldParseIPAddrLoopback(t *testing.T) {
	// GIVEN a template string
	b := `{{IsLoopback "127.0.0.1"}}`
	// WHEN parsing int
	out, err := ParseTemplate(b, &PrincipalExt{Delegate: &types.Principal{}}, &types.Resource{}, &services.AuthRequest{})
	// THEN it should succeed
	require.NoError(t, err)
	require.Equal(t, "true", string(out))
}

func Test_ShouldParseIPAddrIsMulticast(t *testing.T) {
	// GIVEN a template string
	b := `{{IsMulticast "224.0.0.1"}}`
	// WHEN parsing int
	out, err := ParseTemplate(b, &PrincipalExt{Delegate: &types.Principal{}}, &types.Resource{}, &services.AuthRequest{})
	// THEN it should succeed
	require.NoError(t, err)
	require.Equal(t, "true", string(out))
}

func Test_ShouldParseTrue(t *testing.T) {
	// GIVEN a template string
	b := `{{True true}}`
	// WHEN parsing int
	out, err := ParseTemplate(b, &PrincipalExt{Delegate: &types.Principal{}}, &types.Resource{}, &services.AuthRequest{})
	// THEN it should succeed
	require.NoError(t, err)
	require.Equal(t, "true", string(out))
}

func Test_ShouldParseNot(t *testing.T) {
	// GIVEN a template string
	b := `{{Not true}}`
	// WHEN parsing int
	out, err := ParseTemplate(b, &PrincipalExt{Delegate: &types.Principal{}}, &types.Resource{}, &services.AuthRequest{})
	// THEN it should succeed
	require.NoError(t, err)
	require.Equal(t, "false", string(out))
}

func Test_ShouldCalculateDistanceBetweenLatLng(t *testing.T) {
	lat1 := 47.620422
	lon1 := -122.349358
	lat2 := 46.879967
	lon2 := -121.726906
	dist := kmDistance([]float64{lat1, lon1}, []float64{lat2, lon2})
	require.True(t, dist < 100)
	dist = kmDistance(nil, []float64{lat2, lon2})
	require.True(t, dist > 100)
	dist = kmDistance([]float64{lat1, lon1}, nil)
	require.True(t, dist > 100)
}

func Test_ShouldCheckTimeInRange(t *testing.T) {
	require.True(t, IsTimeInRange("11:00am", "10:00am", "8:00pm"))
}

func Test_ShouldCheckIPAddressInRange(t *testing.T) {
	ip := "211.211.211.5"
	mask := "211.211.211.0/24"

	inRange, err := IPInRange(ip, mask)
	require.NoError(t, err)
	require.True(t, inRange)
}

func Test_ShouldCheckIPAddressLoopback(t *testing.T) {
	isLoop, err := IsLoopback("127.0.0.1")
	require.NoError(t, err)
	require.True(t, isLoop)
	isLoop, err = IsLoopback("::1")
	require.NoError(t, err)
	require.True(t, isLoop)
	isLoop, err = IsLoopback("192.168.1.1")
	require.NoError(t, err)
	require.False(t, isLoop)
}

func Test_ShouldCheckIPAddressMulticast(t *testing.T) {
	isMulticast, err := IsMulticast("224.0.0.1") // multicast for IPv4
	require.NoError(t, err)
	require.True(t, isMulticast)
	isMulticast, err = IsMulticast("FF02::1") // multicast for IPv6
	require.NoError(t, err)
	require.True(t, isMulticast)
	isMulticast, err = IsMulticast("192.168.1.1")
	require.NoError(t, err)
	require.False(t, isMulticast)
}

func Test_ShouldCheckIncludesStringOrArray(t *testing.T) {
	require.True(t, includesStringOrArray([]string{"one", "two"}, "one"))
	require.True(t, includesStringOrArray("one ; two", "one"))
	require.True(t, includesStringOrArray("one, two", "one"))
	require.True(t, includesStringOrArray("one, two", ""))
}
