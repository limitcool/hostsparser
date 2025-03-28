package lexer

import (
	"testing"
)

func TestIsDomainChar(t *testing.T) {
	tests := []struct {
		name     string
		ch       rune
		expected bool
	}{
		{"字母", 'a', true},
		{"数字", '1', true},
		{"点", '.', true},
		{"连字符", '-', true},
		{"下划线", '_', true},
		{"空格", ' ', false},
		{"特殊字符", '@', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isDomainChar(tt.ch)
			if result != tt.expected {
				t.Errorf("isDomainChar(%c) = %v; 期望 %v", tt.ch, result, tt.expected)
			}
		})
	}
}

func TestIsDomain(t *testing.T) {
	tests := []struct {
		name     string
		domain   string
		expected bool
	}{
		{"有效域名", "example.com", true},
		{"带下划线的域名", "my_domain.com", true},
		{"带连字符的域名", "my-domain.com", true},
		{"下划线开头", "_example.com", false},
		{"下划线结尾", "example_.com", false},
		{"下划线部分", "my._domain.com", false},
		{"IP地址", "192.168.1.1", false},
		{"IPv6地址", "2001:db8::1", false},
		{"空字符串", "", false},
		{"单个部分", "example", false},
		{"无效字符", "example@.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isDomain(tt.domain)
			if result != tt.expected {
				t.Errorf("isDomain(%s) = %v; 期望 %v", tt.domain, result, tt.expected)
			}
		})
	}
}
