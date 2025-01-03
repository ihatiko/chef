package query

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/ihatiko/chef/components/tech/toml"
)

type queryTestNode struct {
	value    interface{}
	position toml.Position
}

func valueString(root interface{}) string {
	result := "" //fmt.Sprintf("%T:", root)
	switch node := root.(type) {
	case *Result:
		items := []string{}
		for i, v := range node.Values() {
			items = append(items, fmt.Sprintf("%s:%s",
				node.Positions()[i].String(), valueString(v)))
		}
		sort.Strings(items)
		result = "[" + strings.Join(items, ", ") + "]"
	case queryTestNode:
		result = fmt.Sprintf("%s:%s",
			node.position.String(), valueString(node.value))
	case []interface{}:
		items := []string{}
		for _, v := range node {
			items = append(items, valueString(v))
		}
		sort.Strings(items)
		result = "[" + strings.Join(items, ", ") + "]"
	case *toml.Tree:
		// workaround for unreliable map key ordering
		items := []string{}
		for _, k := range node.Keys() {
			v := node.GetPath([]string{k})
			items = append(items, k+":"+valueString(v))
		}
		sort.Strings(items)
		result = "{" + strings.Join(items, ", ") + "}"
	case map[string]interface{}:
		// workaround for unreliable map key ordering
		items := []string{}
		for k, v := range node {
			items = append(items, k+":"+valueString(v))
		}
		sort.Strings(items)
		result = "{" + strings.Join(items, ", ") + "}"
	case int64:
		result += fmt.Sprintf("%d", node)
	case string:
		result += "'" + node + "'"
	case float64:
		result += fmt.Sprintf("%f", node)
	case bool:
		result += fmt.Sprintf("%t", node)
	case time.Time:
		result += fmt.Sprintf("'%v'", node)
	}
	return result
}

func assertValue(t *testing.T, result, ref interface{}) {
	pathStr := valueString(result)
	refStr := valueString(ref)
	if pathStr != refStr {
		t.Errorf("values do not match")
		t.Log("test:", pathStr)
		t.Log("ref: ", refStr)
	}
}

func assertParseError(t *testing.T, query string, errString string) {
	_, err := Compile(query)
	if err == nil {
		t.Error("error should be non-nil")
		return
	}
	if err.Error() != errString {
		t.Errorf("error does not match")
		t.Log("test:", err.Error())
		t.Log("ref: ", errString)
	}
}

func assertQueryPositions(t *testing.T, tomlDoc string, query string, ref []interface{}) {
	tree, err := toml.Load(tomlDoc)
	if err != nil {
		t.Errorf("Non-nil toml parse error: %v", err)
		return
	}
	q, err := Compile(query)
	if err != nil {
		t.Error(err)
		return
	}
	results := q.Execute(tree)
	assertValue(t, results, ref)
}

func TestQueryRoot(t *testing.T) {
	assertQueryPositions(t,
		"a = 42",
		"$",
		[]interface{}{
			queryTestNode{
				map[string]interface{}{
					"a": int64(42),
				}, toml.Position{1, 1},
			},
		})
}

func TestQueryKey(t *testing.T) {
	assertQueryPositions(t,
		"[foo]\na = 42",
		"$.foo.a",
		[]interface{}{
			queryTestNode{
				int64(42), toml.Position{2, 1},
			},
		})
}

func TestQueryKeyString(t *testing.T) {
	assertQueryPositions(t,
		"[foo]\na = 42",
		"$.foo['a']",
		[]interface{}{
			queryTestNode{
				int64(42), toml.Position{2, 1},
			},
		})
}

func TestQueryKeyUnicodeString(t *testing.T) {
	assertQueryPositions(t,
		"['f𝟘.o']\na = 42",
		"$['f𝟘.o']['a']",
		[]interface{}{
			queryTestNode{
				int64(42), toml.Position{2, 1},
			},
		})
}

func TestQueryIndexError1(t *testing.T) {
	assertParseError(t, "$.foo.a[5", "(1, 10): expected ',' or ']', not ''")
}

func TestQueryIndexError2(t *testing.T) {
	assertParseError(t, "$.foo.a[]", "(1, 9): expected union sub expression, not ']', 0")
}

func TestQueryIndex(t *testing.T) {
	assertQueryPositions(t,
		"[foo]\na = [0,1,2,3,4,5,6,7,8,9]",
		"$.foo.a[5]",
		[]interface{}{
			queryTestNode{int64(5), toml.Position{2, 1}},
		})
}

func TestQueryIndexNegative(t *testing.T) {
	assertQueryPositions(t,
		"[foo]\na = [0,1,2,3,4,5,6,7,8,9]",
		"$.foo.a[-2]",
		[]interface{}{
			queryTestNode{int64(8), toml.Position{2, 1}},
		})
}

func TestQueryIndexWrong(t *testing.T) {
	assertQueryPositions(t,
		"[foo]\na = [0,1,2,3,4,5,6,7,8,9]",
		"$.foo.a[99]",
		[]interface{}{})
}

func TestQueryIndexEmpty(t *testing.T) {
	assertQueryPositions(t,
		"[foo]\na = []",
		"$.foo.a[5]",
		[]interface{}{})
}

func TestQueryIndexTree(t *testing.T) {
	assertQueryPositions(t,
		"[[foo]]\na = [0,1,2,3,4,5,6,7,8,9]\n[[foo]]\nb = 3",
		"$.foo[1].b",
		[]interface{}{
			queryTestNode{int64(3), toml.Position{4, 1}},
		})
}

func TestQuerySliceError1(t *testing.T) {
	assertParseError(t, "$.foo.a[3:?]", "(1, 11): expected ']' or ':'")
}

func TestQuerySliceError2(t *testing.T) {
	assertParseError(t, "$.foo.a[:::]", "(1, 11): expected ']'")
}

func TestQuerySliceError3(t *testing.T) {
	assertParseError(t, "$.foo.a[::0]", "(1, 11): step cannot be zero")
}

func TestQuerySliceRange(t *testing.T) {
	assertQueryPositions(t,
		"[foo]\na = [0,1,2,3,4,5,6,7,8,9]",
		"$.foo.a[:5]",
		[]interface{}{
			queryTestNode{int64(0), toml.Position{2, 1}},
			queryTestNode{int64(1), toml.Position{2, 1}},
			queryTestNode{int64(2), toml.Position{2, 1}},
			queryTestNode{int64(3), toml.Position{2, 1}},
			queryTestNode{int64(4), toml.Position{2, 1}},
		})
}

func TestQuerySliceStep(t *testing.T) {
	assertQueryPositions(t,
		"[foo]\na = [0,1,2,3,4,5,6,7,8,9]",
		"$.foo.a[0:5:2]",
		[]interface{}{
			queryTestNode{int64(0), toml.Position{2, 1}},
			queryTestNode{int64(2), toml.Position{2, 1}},
			queryTestNode{int64(4), toml.Position{2, 1}},
		})
}

func TestQuerySliceStartNegative(t *testing.T) {
	assertQueryPositions(t,
		"[foo]\na = [0,1,2,3,4,5,6,7,8,9]",
		"$.foo.a[-3:]",
		[]interface{}{
			queryTestNode{int64(7), toml.Position{2, 1}},
			queryTestNode{int64(8), toml.Position{2, 1}},
			queryTestNode{int64(9), toml.Position{2, 1}},
		})
}

func TestQuerySliceEndNegative(t *testing.T) {
	assertQueryPositions(t,
		"[foo]\na = [0,1,2,3,4,5,6,7,8,9]",
		"$.foo.a[:-6]",
		[]interface{}{
			queryTestNode{int64(0), toml.Position{2, 1}},
			queryTestNode{int64(1), toml.Position{2, 1}},
			queryTestNode{int64(2), toml.Position{2, 1}},
			queryTestNode{int64(3), toml.Position{2, 1}},
		})
}

func TestQuerySliceStepNegative(t *testing.T) {
	assertQueryPositions(t,
		"[foo]\na = [0,1,2,3,4,5,6,7,8,9]",
		"$.foo.a[::-2]",
		[]interface{}{
			queryTestNode{int64(9), toml.Position{2, 1}},
			queryTestNode{int64(7), toml.Position{2, 1}},
			queryTestNode{int64(5), toml.Position{2, 1}},
			queryTestNode{int64(3), toml.Position{2, 1}},
			queryTestNode{int64(1), toml.Position{2, 1}},
		})
}

func TestQuerySliceStartOverRange(t *testing.T) {
	assertQueryPositions(t,
		"[foo]\na = [0,1,2,3,4,5,6,7,8,9]",
		"$.foo.a[-99:3]",
		[]interface{}{
			queryTestNode{int64(0), toml.Position{2, 1}},
			queryTestNode{int64(1), toml.Position{2, 1}},
			queryTestNode{int64(2), toml.Position{2, 1}},
		})
}

func TestQuerySliceStartOverRangeNegative(t *testing.T) {
	assertQueryPositions(t,
		"[foo]\na = [0,1,2,3,4,5,6,7,8,9]",
		"$.foo.a[99:7:-1]",
		[]interface{}{
			queryTestNode{int64(9), toml.Position{2, 1}},
			queryTestNode{int64(8), toml.Position{2, 1}},
		})
}

func TestQuerySliceEndOverRange(t *testing.T) {
	assertQueryPositions(t,
		"[foo]\na = [0,1,2,3,4,5,6,7,8,9]",
		"$.foo.a[7:99]",
		[]interface{}{
			queryTestNode{int64(7), toml.Position{2, 1}},
			queryTestNode{int64(8), toml.Position{2, 1}},
			queryTestNode{int64(9), toml.Position{2, 1}},
		})
}

func TestQuerySliceEndOverRangeNegative(t *testing.T) {
	assertQueryPositions(t,
		"[foo]\na = [0,1,2,3,4,5,6,7,8,9]",
		"$.foo.a[2:-99:-1]",
		[]interface{}{
			queryTestNode{int64(2), toml.Position{2, 1}},
			queryTestNode{int64(1), toml.Position{2, 1}},
			queryTestNode{int64(0), toml.Position{2, 1}},
		})
}

func TestQuerySliceWrongRange(t *testing.T) {
	assertQueryPositions(t,
		"[foo]\na = [0,1,2,3,4,5,6,7,8,9]",
		"$.foo.a[5:3]",
		[]interface{}{})
}

func TestQuerySliceWrongRangeNegative(t *testing.T) {
	assertQueryPositions(t,
		"[foo]\na = [0,1,2,3,4,5,6,7,8,9]",
		"$.foo.a[3:5:-1]",
		[]interface{}{})
}

func TestQuerySliceEmpty(t *testing.T) {
	assertQueryPositions(t,
		"[foo]\na = []",
		"$.foo.a[5:]",
		[]interface{}{})
}

func TestQuerySliceTree(t *testing.T) {
	assertQueryPositions(t,
		"[[foo]]\na='nok'\n[[foo]]\na = [0,1,2,3,4,5,6,7,8,9]\n[[foo]]\na='ok'\nb = 3",
		"$.foo[1:].a",
		[]interface{}{
			queryTestNode{
				[]interface{}{
					int64(0), int64(1), int64(2), int64(3), int64(4),
					int64(5), int64(6), int64(7), int64(8), int64(9)},
				toml.Position{4, 1}},
			queryTestNode{"ok", toml.Position{6, 1}},
		})
}

func TestQueryAny(t *testing.T) {
	assertQueryPositions(t,
		"[foo.bar]\na=1\nb=2\n[foo.baz]\na=3\nb=4",
		"$.foo.*",
		[]interface{}{
			queryTestNode{
				map[string]interface{}{
					"a": int64(1),
					"b": int64(2),
				}, toml.Position{1, 1},
			},
			queryTestNode{
				map[string]interface{}{
					"a": int64(3),
					"b": int64(4),
				}, toml.Position{4, 1},
			},
		})
}
func TestQueryUnionSimple(t *testing.T) {
	assertQueryPositions(t,
		"[foo.bar]\na=1\nb=2\n[baz.foo]\na=3\nb=4\n[gorf.foo]\na=5\nb=6",
		"$.*[bar,foo]",
		[]interface{}{
			queryTestNode{
				map[string]interface{}{
					"a": int64(1),
					"b": int64(2),
				}, toml.Position{1, 1},
			},
			queryTestNode{
				map[string]interface{}{
					"a": int64(3),
					"b": int64(4),
				}, toml.Position{4, 1},
			},
			queryTestNode{
				map[string]interface{}{
					"a": int64(5),
					"b": int64(6),
				}, toml.Position{7, 1},
			},
		})
}

func TestQueryRecursionAll(t *testing.T) {
	assertQueryPositions(t,
		"[foo.bar]\na=1\nb=2\n[baz.foo]\na=3\nb=4\n[gorf.foo]\na=5\nb=6",
		"$..*",
		[]interface{}{
			queryTestNode{
				map[string]interface{}{
					"foo": map[string]interface{}{
						"bar": map[string]interface{}{
							"a": int64(1),
							"b": int64(2),
						},
					},
					"baz": map[string]interface{}{
						"foo": map[string]interface{}{
							"a": int64(3),
							"b": int64(4),
						},
					},
					"gorf": map[string]interface{}{
						"foo": map[string]interface{}{
							"a": int64(5),
							"b": int64(6),
						},
					},
				}, toml.Position{1, 1},
			},
			queryTestNode{
				map[string]interface{}{
					"bar": map[string]interface{}{
						"a": int64(1),
						"b": int64(2),
					},
				}, toml.Position{1, 1},
			},
			queryTestNode{
				map[string]interface{}{
					"a": int64(1),
					"b": int64(2),
				}, toml.Position{1, 1},
			},
			queryTestNode{int64(1), toml.Position{2, 1}},
			queryTestNode{int64(2), toml.Position{3, 1}},
			queryTestNode{
				map[string]interface{}{
					"foo": map[string]interface{}{
						"a": int64(3),
						"b": int64(4),
					},
				}, toml.Position{4, 1},
			},
			queryTestNode{
				map[string]interface{}{
					"a": int64(3),
					"b": int64(4),
				}, toml.Position{4, 1},
			},
			queryTestNode{int64(3), toml.Position{5, 1}},
			queryTestNode{int64(4), toml.Position{6, 1}},
			queryTestNode{
				map[string]interface{}{
					"foo": map[string]interface{}{
						"a": int64(5),
						"b": int64(6),
					},
				}, toml.Position{7, 1},
			},
			queryTestNode{
				map[string]interface{}{
					"a": int64(5),
					"b": int64(6),
				}, toml.Position{7, 1},
			},
			queryTestNode{int64(5), toml.Position{8, 1}},
			queryTestNode{int64(6), toml.Position{9, 1}},
		})
}

func TestQueryRecursionUnionSimple(t *testing.T) {
	assertQueryPositions(t,
		"[foo.bar]\na=1\nb=2\n[baz.foo]\na=3\nb=4\n[gorf.foo]\na=5\nb=6",
		"$..['foo','bar']",
		[]interface{}{
			queryTestNode{
				map[string]interface{}{
					"bar": map[string]interface{}{
						"a": int64(1),
						"b": int64(2),
					},
				}, toml.Position{1, 1},
			},
			queryTestNode{
				map[string]interface{}{
					"a": int64(3),
					"b": int64(4),
				}, toml.Position{4, 1},
			},
			queryTestNode{
				map[string]interface{}{
					"a": int64(1),
					"b": int64(2),
				}, toml.Position{1, 1},
			},
			queryTestNode{
				map[string]interface{}{
					"a": int64(5),
					"b": int64(6),
				}, toml.Position{7, 1},
			},
		})
}

func TestQueryFilterFn(t *testing.T) {
	buff, err := ioutil.ReadFile("../example.toml")
	if err != nil {
		t.Error(err)
		return
	}

	assertQueryPositions(t, string(buff),
		"$..[?(int)]",
		[]interface{}{
			queryTestNode{int64(8001), toml.Position{13, 1}},
			queryTestNode{int64(8001), toml.Position{13, 1}},
			queryTestNode{int64(8002), toml.Position{13, 1}},
			queryTestNode{int64(5000), toml.Position{14, 1}},
		})

	assertQueryPositions(t, string(buff),
		"$..[?(string)]",
		[]interface{}{
			queryTestNode{"TOML Example", toml.Position{3, 1}},
			queryTestNode{"Tom Preston-Werner", toml.Position{6, 1}},
			queryTestNode{"GitHub", toml.Position{7, 1}},
			queryTestNode{"GitHub Cofounder & CEO\nLikes tater tots and beer.", toml.Position{8, 1}},
			queryTestNode{"192.168.1.1", toml.Position{12, 1}},
			queryTestNode{"10.0.0.1", toml.Position{21, 3}},
			queryTestNode{"eqdc10", toml.Position{22, 3}},
			queryTestNode{"10.0.0.2", toml.Position{25, 3}},
			queryTestNode{"eqdc10", toml.Position{26, 3}},
		})

	assertQueryPositions(t, string(buff),
		"$..[?(float)]",
		[]interface{}{
			queryTestNode{4e-08, toml.Position{30, 1}},
		})

	tv, _ := time.Parse(time.RFC3339, "1979-05-27T07:32:00Z")
	assertQueryPositions(t, string(buff),
		"$..[?(tree)]",
		[]interface{}{
			queryTestNode{
				map[string]interface{}{
					"name":         "Tom Preston-Werner",
					"organization": "GitHub",
					"bio":          "GitHub Cofounder & CEO\nLikes tater tots and beer.",
					"dob":          tv,
				}, toml.Position{5, 1},
			},
			queryTestNode{
				map[string]interface{}{
					"server":         "192.168.1.1",
					"ports":          []interface{}{int64(8001), int64(8001), int64(8002)},
					"connection_max": int64(5000),
					"enabled":        true,
				}, toml.Position{11, 1},
			},
			queryTestNode{
				map[string]interface{}{
					"alpha": map[string]interface{}{
						"ip": "10.0.0.1",
						"dc": "eqdc10",
					},
					"beta": map[string]interface{}{
						"ip": "10.0.0.2",
						"dc": "eqdc10",
					},
				}, toml.Position{17, 1},
			},
			queryTestNode{
				map[string]interface{}{
					"ip": "10.0.0.1",
					"dc": "eqdc10",
				}, toml.Position{20, 3},
			},
			queryTestNode{
				map[string]interface{}{
					"ip": "10.0.0.2",
					"dc": "eqdc10",
				}, toml.Position{24, 3},
			},
			queryTestNode{
				map[string]interface{}{
					"data": []interface{}{
						[]interface{}{"gamma", "delta"},
						[]interface{}{int64(1), int64(2)},
					},
					"score": 4e-08,
				}, toml.Position{28, 1},
			},
		})

	assertQueryPositions(t, string(buff),
		"$..[?(time)]",
		[]interface{}{
			queryTestNode{tv, toml.Position{9, 1}},
		})

	assertQueryPositions(t, string(buff),
		"$..[?(bool)]",
		[]interface{}{
			queryTestNode{true, toml.Position{15, 1}},
		})
}
