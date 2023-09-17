package domain

import (
	"bytes"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/utils"
	"html/template"
	"math"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// ParseTemplate parses GO template with dynamic parameters
func ParseTemplate(
	templateStr string,
	principal *PrincipalExt,
	resource *types.Resource,
	req *services.AuthRequest,
) ([]byte, error) {
	if templateStr == "" {
		return nil, NewInternalError("template is not defined", TemplateCode)
	}
	if !strings.Contains(templateStr, "{{") {
		templateStr = "{{" + templateStr + "}}"
	}
	emptyLineRegex, err := regexp.Compile(`(?m)^\s*$[\r\n]*|[\r\n]+\s+\z`)
	if err != nil {
		return nil, err
	}
	t, err := template.New("").Funcs(TemplateFuncs(principal, req)).Parse(templateStr)
	if err != nil {
		return nil, NewMarshalError(
			fmt.Sprintf("failed to parse '%s' template due to %s",
				templateStr, err))
	}
	var out bytes.Buffer
	data := principal.ToMap(req, resource)
	err = t.Execute(&out, data)
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
			"Body":  templateStr,
			"Data":  data,
		}).Warnf("failed to execute template")
		return nil, NewInternalError(fmt.Sprintf("failed to execute template due to '%s', data=%v",
			err, data), TemplateCode)
	}
	strResponse := strings.TrimSpace(emptyLineRegex.ReplaceAllString(out.String(), ""))
	//if unescape {
	//	strResponse = strings.ReplaceAll(strResponse, "&lt;", "<")
	//}
	return []byte(strResponse), nil
}

// TemplateFuncs returns template functions
func TemplateFuncs(
	principal *PrincipalExt,
	req *services.AuthRequest,
) template.FuncMap {
	return template.FuncMap{
		"Dict": func(values ...any) (map[string]any, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("invalid dict call")
			}
			dict := make(map[string]any, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"Iterate": func(input any) []int {
			count := utils.ToInt(input)
			var i int
			var Items []int
			for i = 0; i < count; i++ {
				Items = append(Items, i)
			}
			return Items
		},
		"Unescape": func(s string) template.HTML {
			return template.HTML(s)
		},
		"Int": func(num any) int64 {
			return utils.ToInt64(num)
		},
		"Float": func(num any) float64 {
			return utils.ToFloat64(num)
		},
		"Not": func(s any) bool {
			return !utils.ToBoolean(s)
		},
		"True": func(s any) bool {
			return utils.ToBoolean(s)
		},
		"LT": func(a any, b any) bool {
			return utils.ToFloat64(a) < utils.ToFloat64(b)
		},
		"LE": func(a any, b any) bool {
			return utils.ToFloat64(a) <= utils.ToFloat64(b)
		},
		"EQ": func(a any, b any) bool {
			return utils.ToFloat64(a) == utils.ToFloat64(b)
		},
		"GT": func(a any, b any) bool {
			return utils.ToFloat64(a) > utils.ToFloat64(b)
		},
		"GE": func(a any, b any) bool {
			return utils.ToFloat64(a) >= utils.ToFloat64(b)
		},
		"Nth": func(a any, b any) bool {
			return utils.ToInt(a)%utils.ToInt(b) == 0
		},
		"TimeNow": func(format string) string {
			return time.Now().Format(format)
		},
		"TimeInRange": func(current string, start string, end string) bool {
			return IsTimeInRange(current, start, end)
		},
		"Includes": func(arr any, s string) bool {
			return includesStringOrArray(arr, s)
		},
		"HasRole": func(role string) bool {
			return utils.Includes(principal.RoleNames(), role)
		},
		"ActionIncludes": func(actions ...string) bool {
			return utils.Includes(actions, req.Action)
		},
		"HasGroup": func(group string) bool {
			return utils.Includes(principal.GroupNames(), group)
		},
		"HasRelation": func(relation string) bool {
			return utils.Includes(principal.RelationNamesByResourceName(req.Resource), relation)
		},
		"BeginsWith": func(full string, partial string) bool {
			return strings.HasPrefix(full, partial)
		},
		"EndsWith": func(full string, partial string) bool {
			return strings.HasSuffix(full, partial)
		},
		"Contains": func(full string, partial string) bool {
			return strings.Contains(full, partial)
		},
		"StringToFloatArray": func(s string) (arr []float64) {
			return parseLatLng(s)
		},
		"IPInRange": func(ipAddr string, cidr string) bool {
			b, _ := IPInRange(ipAddr, cidr)
			return b
		},
		"IsLoopback": func(ipAddr string) bool {
			b, _ := IsLoopback(ipAddr)
			return b
		},
		"IsMulticast": func(ipAddr string) bool {
			b, _ := IsMulticast(ipAddr)
			return b
		},
		"DistanceWithinKM": func(s1 string, s2 string, withinS any) bool {
			latlng1 := parseLatLng(s1)
			latlng2 := parseLatLng(s2)
			within := utils.ToFloat64(withinS)
			if within <= 0 {
				return false
			}
			return kmDistance(latlng1, latlng2) <= within
		},
	}
}

func includesStringOrArray(input any, target string) bool {
	pattern := `\s+|[!@#$%^&*(),.?":{}|<>]`
	var items []string

	switch v := input.(type) {
	case string:
		// Split string by whitespace or special characters
		reg := regexp.MustCompile(pattern)
		items = reg.Split(v, -1)
	case []string:
		items = v
	default:
		// Split string by whitespace or special characters
		reg := regexp.MustCompile(pattern)
		vv := fmt.Sprintf("%v", input)
		items = reg.Split(vv, -1)
	}
	return utils.Includes(items, target)
}

func parseLatLng(s string) (arr []float64) {
	expr := regexp.MustCompile(`[\s,;:]`)
	strArr := expr.Split(s, -1)
	for _, next := range strArr {
		if n, err := strconv.ParseFloat(next, 64); err == nil {
			arr = append(arr, n)
		}
	}
	return arr
}

// IsTimeInRange checks if time is in range of start and end time.
func IsTimeInRange(
	currentTimeStr string,
	startTimeStr string,
	endTimeStr string,
) bool {
	// Convert 12-hour format strings to 24-hour format strings
	layout := "3:04pm"
	currentTime, _ := time.Parse(layout, currentTimeStr)
	startTime, _ := time.Parse(layout, startTimeStr)
	endTime, _ := time.Parse(layout, endTimeStr)

	// Convert parsed time to today's date
	startTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), startTime.Hour(), startTime.Minute(), 0, 0, currentTime.Location())
	endTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), endTime.Hour(), endTime.Minute(), 0, 0, currentTime.Location())

	return currentTime.After(startTime) && currentTime.Before(endTime)
}

// kmDistance defines helper method to calculate distance between two latitude/longitude
func kmDistance(latlng1 []float64, latlng2 []float64) float64 {
	if len(latlng1) != 2 || len(latlng2) != 2 {
		return math.Inf(1)
	}
	earthRadiusKm := 6371.0
	deglat := toRadians(latlng2[0] - latlng1[0])
	deglon := toRadians(latlng2[1] - latlng1[1])
	lat1 := toRadians(latlng1[0])
	lat2 := toRadians(latlng2[0])

	a := math.Sin(deglat/2.0)*math.Sin(deglat/2.0) +
		math.Sin(deglon/2.0)*math.Sin(deglon/2.0)*
			math.Cos(lat1)*math.Cos(lat2)
	c := 2.0 * math.Atan2(math.Sqrt(a), math.Sqrt(1.0-a))
	return earthRadiusKm * c
}

func toRadians(degrees float64) float64 {
	return math.Pi * degrees / 180.0
}

// IPInRange checks if ipaddress is in CIDR range
func IPInRange(ipAddress string, cidr string) (bool, error) {
	// Parse the CIDR range.
	_, subnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false, err
	}

	// Parse the IP address.
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return false, fmt.Errorf("invalid IP address: %s", ipAddress)
	}

	// Check if the IP is in the CIDR range.
	return subnet.Contains(ip), nil
}

// IsLoopback checks if the provided IP address is a loopback address.
func IsLoopback(ipAddress string) (bool, error) {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return false, fmt.Errorf("invalid IP address: %s", ipAddress)
	}
	return ip.IsLoopback(), nil
}

// IsMulticast checks if the provided IP address is a multicast address.
func IsMulticast(ipAddress string) (bool, error) {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return false, fmt.Errorf("invalid IP address: %s", ipAddress)
	}
	return ip.IsMulticast(), nil
}
