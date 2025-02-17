package nvimboat

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/neovim/go-client/nvim"
)

func extracUrls(content string) (links []string) {
	re := regexp.MustCompile(`\b((?:https?|ftp|file):\/\/[-a-zA-Z0-9+&@#\/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#\/%=~_|])`)
	matches := re.FindAll([]byte(content), -1)
	for _, l := range matches {
		links = append(links, string(l))
	}
	return links
}

func renderHTML(content string) ([]string, error) {
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(content)
	if err != nil {
		return nil, err
	}
	return strings.Split(markdown, "\n"), nil
}

func unixToDate(unixTime int64) (string, error) {

	tz, err := time.LoadLocation("Local")
	if err != nil {
		return "", err
	}
	ut := time.Unix(unixTime, 0)
	dateString := ut.In(tz).Format("02 Jan 06")

	return dateString, nil
}

func trimTrail(nv *nvim.Nvim, buffer nvim.Buffer) (err error) {
	currentLines, err := nv.BufferLines(buffer, 0, -1, false)
	if err != nil {
		return
	}
	var lines []string
	for _, l := range currentLines {
		lines = append(lines, strings.TrimRight(string(l), " "))
	}
	err = setLines(nv, buffer, lines)
	return
}

func addColumn(nv *nvim.Nvim, buf nvim.Buffer, col []string) (err error) {
	currentLines, err := nv.BufferLines(buf, 0, -1, false)
	if err != nil {
		return
	}
	var (
		diff  int
		lines = []string{}
	)
	diff = (len(col) - len(currentLines))
	for i := 0; i < diff; i++ {
		currentLines = append(currentLines, []byte{})
	}
	for i, c := range col {
		lines = append(lines, string(currentLines[i])+" │ "+c)
	}
	err = setLines(nv, buf, lines)
	if err != nil {
		err = fmt.Errorf("addColumn: %w", err)
		return
	}
	vcl, err := virtColLens(nv)
	if err != nil {
		err = fmt.Errorf("addColumn: %w", err)
		return
	}
	maxLineLen := slices.Max(vcl)

	for i, l := range lines {
		diff = maxLineLen - vcl[i]
		lines[i] = l + strings.Repeat(" ", diff)
	}
	err = setLines(nv, buf, lines)
	if err != nil {
		err = fmt.Errorf("addColumn: %w", err)
		return
	}
	return err
}

func setLines(nv *nvim.Nvim, buffer nvim.Buffer, lines []string) (err error) {
	err = nv.SetBufferLines(buffer, 0, -1, false, strings2bytes(lines))
	if err != nil {
		err = fmt.Errorf("setLines: %w", err)
		return
	}
	return
}

func virtColLens(nv *nvim.Nvim) (evalResult []int, err error) {
	virtCols := "map(range(1, line('$')), \"virtcol([v:val, '$'])\")"
	err = nv.Eval(virtCols, &evalResult)
	if err != nil {
		err = fmt.Errorf("virtCols: %w", err)
		return
	}
	return
}

func makeUnreadRatio(unreadCount, articleCount int) (prefix string) {
	if unreadCount > 0 {
		prefix = "N (" + strconv.Itoa(unreadCount) + "/" + strconv.Itoa(articleCount) + ")"
		return
	}
	prefix = "(" + strconv.Itoa(unreadCount) + "/" + strconv.Itoa(articleCount) + ")"
	return
}

func strings2bytes(stringSlice []string) (byteSlices [][]byte) {
	for _, s := range stringSlice {
		byteSlices = append(byteSlices, []byte(s))
	}
	return
}

func sortMapKeys(m interface{}) (keyList []string) {
	keys := reflect.ValueOf(m).MapKeys()
	for _, key := range keys {
		keyList = append(keyList, key.Interface().(string))
	}
	sort.Strings(keyList)
	return
}

func prettyStruct(s any) (pretty string) {
	marshal, _ := json.MarshalIndent(s, "", "	")
	return string(marshal)
}
