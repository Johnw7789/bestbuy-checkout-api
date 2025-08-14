package shr

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
)

var (
	akamaiPath = regexp.MustCompile(`<script type="text/javascript"\s+(?:nonce=".*")?\s+src="((?i)[a-z\d/\-_]+)"></script>`)
	scriptPathExpr = regexp.MustCompile(`(?i)src\s*=\s*["'](\/[A-Za-z0-9._\-/]+)\?v=([A-Za-z0-9-]+)&t=([A-Za-z0-9]+)`)
	scriptPathExprV2 = regexp.MustCompile(`(?i)(\/[A-Za-z0-9._\-/]+\?v=[A-Za-z0-9-]+)`)
)

func GetRandStr(length int) string {
	var chars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz123456789")

	randChars := make([]rune, length)
	for i := range randChars {
		randChars[i] = chars[rand.Intn(len(chars))]
	}

	return string(randChars)
}

func Parse(value string, a string, b string) string {
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}

	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}

	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}

	return value[posFirstAdjusted:posLast]
}

func ParseV2(str string, start string, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return
	}

	s += len(start)

	e := strings.Index(str[s:], end)
	if e == -1 {
		return
	}

	e += s + e - 1

	return str[s:e]
}

func ParseV3(str string, start string, end string) string {
	var match []byte
	index := strings.Index(str, start)

	if index == -1 {
		return string(match)
	}

	index += len(start)

	for {
		char := str[index]

		if strings.HasPrefix(str[index:index+len(match)], end) {
			break
		}

		match = append(match, char)
		index++
	}

	return string(match)
}

func Reverse(s string) string {
	rs := []rune(s)

	for i, j := 0, len(rs)-1; i < j; i, j = i+1, j-1 {
		rs[i], rs[j] = rs[j], rs[i]
	}

	return string(rs)
}

// pretty print struct
func PPJson(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}

func ParseScriptPath(str string) string {
	scriptPath := Parse(str, `<script src="`, `"></script>`)

	return scriptPath
}

func GetIPAddr() (string, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}

// * This seems to work for both forms of sbsd we could encounter: (1) sbsd script blocking access to site and (2) sbsd script not blocking access, instead just present in the home/prod/etc page html
func ParseSBSD(reader io.Reader) (string, string, string, error) {
	// scriptPathExpr = regexp.MustCompile(`(?i)(\/[A-Za-z0-9._\-/]+\?v=[A-Za-z0-9-]+)`)
	// scriptPathExpr = regexp.MustCompile(`<script type="text/javascript"\s+(?:nonce=".*")?\s+src="((?i)[a-z\d/\-_]+)"></script>`)

	// * Attempt to find the script path in the format: /MeHEyIiHD_WeC/4vUYQz_tbhFEh/w/cGw3S6zu/WwJcAlNnRw/IyJre/Ew5VEg?v=5b80929a-bd85-5c59-8fa7-5935475b6006&t=170
	src, err := io.ReadAll(reader)
	if err != nil {
		return "", "", "", errors.New("could not find sbsd script path")
	}

	matches := scriptPathExpr.FindSubmatch(src)

	if len(matches) < 4 {
		// * Attempt to find the script path in the format: /MeHEyIiHD_WeC/4vUYQz_tbhFEh/w/cGw3S6zu/WwJcAlNnRw/IyJre/Ew5VEg?v=5b80929a-bd85-5c59-8fa7-5935475b6006
		matches := scriptPathExprV2.FindSubmatch(src)

		if len(matches) < 2 {
			return "", "", "", errors.New("sbsd script path not found (2)")
		}

		path := string(matches[1])
		pathSplit := strings.Split(path, "?v=")
		if len(pathSplit) < 2 {
			return "", "", "", errors.New("sbsd script path not found (3)")
		}

		uuid := pathSplit[1]

		return string(matches[1]), "", uuid, nil
	}

	return string(matches[1]) + "?v=" + string(matches[2]) + "&t=" + string(matches[3]), string(matches[1]), string(matches[2]), nil
}

func ParseAkamaiPath(reader io.Reader) (string, error) {
	src, err := io.ReadAll(reader)
	if err != nil {
		return "", errors.New("could not read akamai script body")
	}

	matches := akamaiPath.FindSubmatch(src)
	if len(matches) < 2 {
		return "", errors.New("akamai script path not found")
	}

	return string(matches[1]), nil
}