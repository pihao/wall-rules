package main

import (
	"reflect"
	"testing"
)

func Test_map2arr(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]struct{}
		want [][]byte
	}{
		{"", map[string]struct{}{"a": {}, "b": {}, "c": {}}, [][]byte{[]byte("a"), []byte("b"), []byte("c")}},
		{"", map[string]struct{}{"c": {}, "a": {}, "b": {}}, [][]byte{[]byte("a"), []byte("b"), []byte("c")}},
		{"", map[string]struct{}{"c": {}, "b": {}, "a": {}}, [][]byte{[]byte("a"), []byte("b"), []byte("c")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := map2arr(tt.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("map2arr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseDomain(t *testing.T) {
	tests := []struct {
		name       string
		b          []byte
		wantDomain []byte
		wantAllow  bool
	}{
		{"", []byte(`@@||jike.com`), []byte(`jike.com`), true},
		{"", []byte(`@@|http://translate.google.cn`), []byte(`translate.google.cn`), true},
		{"", []byte(`@@|http://www.google.cn/maps`), []byte(`www.google.cn`), true},
		{"", []byte(`@@||http2.golang.org`), []byte(`http2.golang.org`), true},
		{"", []byte(`@@||gov.cn`), []byte(`gov.cn`), true},

		{"", []byte(`google.com`), []byte(`google.com`), false},
		{"", []byte(`.mybet.com`), []byte(`mybet.com`), false},
		{"", []byte(`me.youthwant.com.tw`), []byte(`me.youthwant.com.tw`), false},
		{"", []byte(`|http://*.pimg.tw/`), []byte(`pimg.tw`), false},

		{"", []byte(`|http://cdn*.search.xxx/`), []byte(`search.xxx`), false},
		{"", []byte(`|https://fbcdn*.akamaihd.net/`), []byte(`akamaihd.net`), false},
		{"", []byte(`|http://*2.bahamut.com.tw`), []byte(`bahamut.com.tw`), false},
		{"", []byte(`|https://ss*.4sqi.net`), []byte(`4sqi.net`), false},
		{"", []byte(`|http://hum*.uchicago.edu/faculty/ywang/history`), []byte(`uchicago.edu`), false},
		{"", []byte(`|http://cdn*.xda-developers.com`), []byte(`xda-developers.com`), false},

		{"", []byte(`64memo`), nil, false},
		{"", []byte(`aHR0cHM6Ly95ZWNsLm5ldA`), nil, false},
		{"", []byte(`freenet`), nil, false},
		{"", []byte(`.google.*/falun`), nil, false},
		{"", []byte(`phobos.apple.com*/video`), nil, false},
		{"", []byte(`q=freedom`), nil, false},
		{"", []byte(`q%3Dfreedom`), nil, false},
		{"", []byte(`remembering_tiananmen_20_years`), nil, false},
		{"", []byte(`search*safeweb`), nil, false},
		{"", []byte(`q=triangle`), nil, false},
		{"", []byte(`q%3DTriangle`), nil, false},
		{"", []byte(`ultrareach`), nil, false},
		{"", []byte(`ultrasurf`), nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDomain, gotAllow := parseDomain(tt.b)
			if !reflect.DeepEqual(gotDomain, tt.wantDomain) {
				t.Errorf("parseDomain() gotDomain = %s, want %s", gotDomain, tt.wantDomain)
			}
			if gotAllow != tt.wantAllow {
				t.Errorf("parseDomain() gotAllow = %v, want %v", gotAllow, tt.wantAllow)
			}
		})
	}
}
