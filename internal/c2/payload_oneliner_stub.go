package c2

import "fmt"

// OnelinerKind 单行 payload 的语言/形式
type OnelinerKind string

const (
	OnelinerBash       OnelinerKind = "bash"
	OnelinerNc         OnelinerKind = "nc"
	OnelinerNcMkfifo   OnelinerKind = "nc_mkfifo"
	OnelinerPython     OnelinerKind = "python"
	OnelinerPerl       OnelinerKind = "perl"
	OnelinerPowerShell OnelinerKind = "powershell"
	OnelinerCurl       OnelinerKind = "curl_beacon"
)

func AllOnelinerKinds() []OnelinerKind {
	return []OnelinerKind{
		OnelinerBash, OnelinerNc, OnelinerNcMkfifo,
		OnelinerPython, OnelinerPerl,
		OnelinerPowerShell, OnelinerCurl,
	}
}

var tcpOnelinerKinds = map[OnelinerKind]bool{
	OnelinerBash: true, OnelinerNc: true, OnelinerNcMkfifo: true,
	OnelinerPython: true, OnelinerPerl: true, OnelinerPowerShell: true,
}

var httpOnelinerKinds = map[OnelinerKind]bool{
	OnelinerCurl: true,
}

func OnelinerKindsForListener(listenerType string) []OnelinerKind {
	switch ListenerType(listenerType) {
	case ListenerTypeTCPReverse:
		return []OnelinerKind{
			OnelinerBash, OnelinerNc, OnelinerNcMkfifo,
			OnelinerPython, OnelinerPerl, OnelinerPowerShell,
		}
	case ListenerTypeHTTPBeacon, ListenerTypeHTTPSBeacon, ListenerTypeWebSocket:
		return []OnelinerKind{OnelinerCurl}
	default:
		return nil
	}
}

func IsOnelinerCompatible(listenerType string, kind OnelinerKind) bool {
	switch ListenerType(listenerType) {
	case ListenerTypeTCPReverse:
		return tcpOnelinerKinds[kind]
	case ListenerTypeHTTPBeacon, ListenerTypeHTTPSBeacon, ListenerTypeWebSocket:
		return httpOnelinerKinds[kind]
	default:
		return false
	}
}

type OnelinerInput struct {
	Kind         OnelinerKind
	Host         string
	Port         int
	HTTPBaseURL  string
	ImplantToken string
}

func GenerateOneliner(in OnelinerInput) (string, error) {
	return "", fmt.Errorf("oneliner generation unavailable: restore internal/c2/payload_oneliner.go (Windows Defender may block it)")
}

func urlEncodeForShell(s string) string {
	return s
}
