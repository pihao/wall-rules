package main

import (
	"testing"
)

func Test_validDomain(t *testing.T) {
	tests := []struct {
		name string
		b    []byte
		want bool
	}{
		{"", []byte(`google.com`), true},
		{"", []byte(`www.google.com`), true},
		{"", []byte(`me.youthwant.com.tw`), true},
		{"", []byte(`*.pimg.tw`), true},

		{"", []byte(`64memo`), false},
		{"", []byte(`aHR0cHM6Ly95ZWNsLm5ldA`), false},
		{"", []byte(`freenet`), false},
		{"", []byte(`.google.*`), false},
		{"", []byte(`phobos.apple.com*`), false},
		{"", []byte(`q=freedom`), false},
		{"", []byte(`q%3Dfreedom`), false},
		{"", []byte(`remembering_tiananmen_20_years`), false},
		{"", []byte(`search*safeweb`), false},
		{"", []byte(`q=triangle`), false},
		{"", []byte(`q%3DTriangle`), false},
		{"", []byte(`ultrareach`), false},
		{"", []byte(`ultrasurf`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validDomain(tt.b); got != tt.want {
				t.Errorf("validDomain(%s) = %v, want %v", tt.b, got, tt.want)
			}
		})
	}
}
