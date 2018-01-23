package goics

import (
	"strings"
)

const (
	vParamSep = ";"
)

// IcsNode is a basic token.., with, key, val, and extra params
// to define the type of val.
type IcsNode struct {
	Key    string
	Val    string
	Params map[string]string
}

// ParamsLen returns how many params has a token
func (n *IcsNode) ParamsLen() int {
	if n.Params == nil {
		return 0
	}
	return len(n.Params)
}

// GetOneParam resturns the first param found
// usefull when you know that there is a only
// one param token
func (n *IcsNode) GetOneParam() (string, string) {
	if n.ParamsLen() == 0 {
		return "", ""
	}
	var key, val string
	for k, v := range n.Params {
		key, val = k, v
		break
	}
	return key, val
}

// DecodeLine extracts key, val and extra params from a line
func DecodeLine(line string) *IcsNode {
	if strings.Contains(line, keySep) == false {
		return &IcsNode{}
	}
	key, val := getKeyVal(line)
	//@todo test if val containes , multipleparams
	if strings.Contains(key, vParamSep) == false {
		return &IcsNode{
			Key: key,
			Val: val,
		}
	} 
	// Extract key
	firstParam := strings.Index(key, vParamSep)
	realkey := key[0:firstParam]
	n := &IcsNode{
		Key: realkey,
		Val: val,
	}
	// Extract params
	params := key[firstParam+1:]
	n.Params = decodeParams(params)
	return n
}

// decode extra params linked in key val in the form
// key;param1=val1:val
func decodeParams(arr string) map[string]string {

	p := make(map[string]string)
	var isQuoted = false
	var isParam = true
	var curParam string
	var curVal string
	for _, c := range arr {
		switch {
		// if string is quoted, wait till next quote
		// and capture content
		case c == '"':
			if isQuoted == false {
				isQuoted = true
			} else {
				p[curParam] = curVal
				isQuoted = false
			}
		case c == '=' && isQuoted == false:
			isParam = false
		case c == ';' && isQuoted == false:
			isParam = true
			p[curParam] = curVal
			curParam = ""
			curVal = ""
		default:
			if isParam {
				curParam = curParam + string(c)
			} else {
				curVal = curVal + string(c)
			}
		}
	}
	p[curParam] = curVal
	return p

}

// Returns a key, val... for a line..
func getKeyVal(s string) (key, value string) {
	p := strings.SplitN(s, keySep, 2)
	return p[0], p[1]
}
