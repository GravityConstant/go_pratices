package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"mime/quotedprintable"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	prefixEncoding = `=?GB2312?B?`
	suffixEncoding = `?=`
)

var (
	subject   = `=?GB2312?B?vLzK9bK/s8LOsMbm16rV/cnqx+s=?=`
	textPlain = `1/C+tLXEwey1vKO6DQogDQogICAgxPq6wyENCiANCqGhztLKx7y8yvWyv7XEs8LOsMbmo6zT2jIwMTjE6jjUwjE2yNXI69awuavLvqOstb298czs0tG+rdPQwb249tTCsOsu1NrV4rbOyrG85MDvo6zO0rmk1/fFrMGmo6zE3Lm7yqTIzrjazru5pNf3o6zP1s/yuavLvszhs/bXqtX9yerH66GjDQogDQogICAgzai5/b+qzbfBvdbctcS5pNf3o6zO0tbY0MLK7M+kwcvI7bz+v6q3orXE1fu49sH3s8yho9PaOdTCt93O0tX9yr2yztPruavLvrjVuNXG9LavtcTSu7j20MK1xM/uxL+ho9TauaTX99bQo6zO0tK71rHRz7jx0qrH89fUvLqjrMjP1ea8sMqx1/a6w8HstbyyvNbDtcTDv9K7z+7Izs7xo6y8sMqx0+u/zbf+sr+8sNTL06qyv82sysLM1sLbs8zQ8r+qt6K1xL7fzOXQ6Mfzo6zC+tfjz+C52LK/w8W21NDCz7XNs7XE0qrH86Gj1NqzzNDyv6q3orn9s8zW0LK7ts/M4bjfs+TKtdfUvLqjrL/Lt/7U2s+1zbPR0Leiuf2zzNbQtcS8vMr10NTJz7XE0qrH86Gj1NoxMNTCtde7+bG+zeqzycjtvP61xLXa0rvG2r+qt6LIzs7xoaMNCiANCiAgICC5q8u+v+3Lycjax6K1xLmk1/e31c6no6zNxb3hz/LJz7XExvPStc7Eu6+jrMq5ztLU2r3Ptsy1xMqxvOTE2srK06bBy9XiwO+1xLmk1/e7t76zo6zNrMqxyMPO0rrcv+zT682sysLDx7PJzqrBy7rcusO1xLmk1/e777DpoaO+rbn91eLBvbj21MKw67XEuaTX96OsztLP1tTa0tG+rcTcubu2wMGitKbA7bG+1rC5pNf3o6y1sci7ztK7udPQuty24LK71+O1xLXYt72jrLSmwO3OyszitcS+rdHpt73D5tPQtP3M4bjfo6zNxbbT0K3X98TcwabSstDo0qq9+NK7sr3U9se/o6zQ6NKqsru2z7zM0PjRp8+w0tTM4bjf19S8urXExNzBpqGjDQogDQogICAgztK63M+yu7bV4rfduaTX96Os1eLBvbj21MKw68C0ztLRp7W9wcu63Lbgo6y40M7ywcu63LbgoaPX986q0rvD+8rU08PG2tSxuaSjrM7Sz6PN+8Tc1Nq5q8u+09C4/LrDtcS3otW5o6y4/LzTxsjH0LXEz6PN+9LU0rvD+9X9yr3UsbmktcTJ7bfd1NrV4sDvuaTX96OsyrXP1tfUvLq1xLfctrfEv7Hqo6zM5c/W19S8urXEyMvJ+rzb1rWjrLrNuavLvtK7xvCzybOkoaPSss+jzfvO0tTaseCzzLXE1eLM9bXAwrfJz9S919/UvdS2o6yzyc6q0rvD+7P2yau1xLPM0PLUsaGjDQogDQogICAg1Nq0y87SzOGz9tX9yr3J6sfro6zPo837xNzU2jEx1MK33bmry764+NPo16rV/cX617yjrL/Sx+u49867wey1vLj4ztK8zND4ts3BttfUvLqhosq1z9bA7c/rtcS7+rvhoaPO0rvh08PHq9DptcTMrLbIus2xpcL6tcTIyMfp1/a6w87StcSxvtawuaTX96Oszqq5q8u+tLTU7Lzb1rWjrM2suavLvtK7xvDVuc37w8C6w7XEzrTAtCENCiANCiAgICC0y9bCDQogICAgvrTA8SENCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIMnqx+vIy6O6s8LOsMbmDQogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDIwMTjE6jEw1MIzMMjVDQoNCg0KcGhwMDJAMjAwOTQwMC5jbg0K`
)

func main() {
	subjStr := decodeHeader(subject)
	resb, err := DecodeGB2312([]byte(subjStr))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(resb))

	subjStr = decodeBase64Body(textPlain)
	resb, err = DecodeGB2312([]byte(subjStr))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(resb))
}

func base64Decode(s string) (string, error) {
	res, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func decodeHeader(s string) string {
	var res string

	s = filterLF(s)
	ss := strings.Split(s, " ")
	for _, val := range ss {
		tmp := val
		if strings.Index(val, prefixEncoding) == 0 {
			tmp, _ = base64Decode(strings.TrimSuffix(strings.TrimPrefix(val, prefixEncoding), suffixEncoding))
		}
		res += tmp
	}

	return res
}

func decodeBase64Body(s string) string {
	tmp, err := base64Decode(s)
	if err != nil {
		return ""
	}
	return tmp
}

func filterLF(input string) string {
	var res []string

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		tmp := scanner.Text()
		res = append(res, strings.TrimSpace(tmp))
	}

	if err := scanner.Err(); err != nil {
		return ""
	}

	return strings.Join(res, " ")
}

func decodeQuotedPrintable(filename string) string {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return ""
	}
	b, err := ioutil.ReadAll(quotedprintable.NewReader(bytes.NewReader(bs)))
	if err != nil {
		return ""
	}
	return string(b)
}

func DecodeGB2312(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
