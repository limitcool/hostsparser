package lexer

import (
	"strings"
	"unicode"
)

type TokenType string

const (
	IP         TokenType = "IP"
	DOMAIN     TokenType = "DOMAIN"
	COMMENT    TokenType = "COMMENT"
	WHITESPACE TokenType = "WHITESPACE"
	NEWLINE    TokenType = "NEWLINE"
	EOF        TokenType = "EOF"
)

type Token struct {
	Type  TokenType
	Value string
	Line  int
	Col   int
}

func NewToken(tokenType TokenType, value string, line, col int) Token {
	return Token{Type: tokenType, Value: value, Line: line, Col: col}
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t'
}

func isIPChar(ch rune) bool {
	return unicode.IsDigit(ch) || ch == '.' || ch == ':' || ch == '[' || ch == ']'
}

func isDomainChar(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '.' || ch == '-'
}

// 检查是否是IPv4地址
func isIPv4(s string) bool {
	parts := strings.Split(s, ".")
	if len(parts) != 4 {
		return false
	}

	for _, part := range parts {
		if part == "" {
			return false
		}

		for _, ch := range part {
			if !unicode.IsDigit(ch) {
				return false
			}
		}
	}
	return true
}

// 检查是否是IPv6地址
func isIPv6(s string) bool {
	return strings.Contains(s, ":") || strings.Contains(s, "[") || strings.Contains(s, "]")
}

// 检查是否是有效的域名
func isDomain(s string) bool {
	if strings.Contains(s, ".") && !isIPv4(s) && !isIPv6(s) {
		parts := strings.Split(s, ".")
		// 至少有两部分且每部分不为空
		if len(parts) >= 2 {
			for _, part := range parts {
				if part == "" {
					return false
				}
			}
			return true
		}
	}
	return false
}

func Lex(input string) []Token {
	tokens := []Token{}
	line := 1
	col := 1
	inputRunes := []rune(input)
	pos := 0

	for pos < len(inputRunes) {
		ch := inputRunes[pos]

		// 处理注释
		if ch == '#' {
			startCol := col
			commentStr := string(ch)
			pos++
			col++

			// 读取整行注释直到换行符
			for pos < len(inputRunes) && inputRunes[pos] != '\n' {
				commentStr += string(inputRunes[pos])
				pos++
				col++
			}

			tokens = append(tokens, NewToken(COMMENT, commentStr, line, startCol))
			continue
		}

		// 处理空白字符
		if isWhitespace(ch) {
			startCol := col
			whitespace := string(ch)
			pos++
			col++

			// 连续的空白字符合并为一个Token
			for pos < len(inputRunes) && isWhitespace(inputRunes[pos]) {
				whitespace += string(inputRunes[pos])
				pos++
				col++
			}

			tokens = append(tokens, NewToken(WHITESPACE, whitespace, line, startCol))
			continue
		}

		// 处理换行符
		if ch == '\n' {
			tokens = append(tokens, NewToken(NEWLINE, string(ch), line, col))
			pos++
			line++
			col = 1
			continue
		}

		// 处理IP地址
		if unicode.IsDigit(ch) || ch == ':' {
			startCol := col
			ipStr := string(ch)
			pos++
			col++

			// 收集IP地址的字符
			for pos < len(inputRunes) && (isIPChar(inputRunes[pos]) || inputRunes[pos] == '/') {
				ipStr += string(inputRunes[pos])
				pos++
				col++
			}

			// 判断是否是IP地址或域名
			if isIPv4(ipStr) || isIPv6(ipStr) {
				tokens = append(tokens, NewToken(IP, ipStr, line, startCol))
			} else {
				// 如果不符合IP格式规则，尝试作为域名处理
				tokens = append(tokens, NewToken(DOMAIN, ipStr, line, startCol))
			}
			continue
		}

		// 处理域名
		if unicode.IsLetter(ch) {
			startCol := col
			domainStr := string(ch)
			pos++
			col++

			// 收集域名的字符
			for pos < len(inputRunes) && isDomainChar(inputRunes[pos]) {
				domainStr += string(inputRunes[pos])
				pos++
				col++
			}

			tokens = append(tokens, NewToken(DOMAIN, domainStr, line, startCol))
			continue
		}

		// 处理无法识别的字符
		tokens = append(tokens, NewToken(TokenType("UNKNOWN"), string(ch), line, col))
		pos++
		col++
	}

	// 添加文件结束Token
	tokens = append(tokens, NewToken(EOF, "", line, col))
	return tokens
}
