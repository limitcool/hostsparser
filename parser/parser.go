package parser

import "hostsparser/lexer"

type HostsEntry struct {
	IP            string
	Domains       []string
	HasComment    bool
	Comment       string
	Line          int
	Col           int
	IsIPv6        bool
	IsEmptyLine   bool
	IsCommentLine bool
}

// IPDomainPair 表示IP地址和对应的域名数组
type IPDomainPair struct {
	IP      string
	Domains []string
}

func ParseHosts(tokens []lexer.Token) []HostsEntry {
	var entries []HostsEntry     // 初始化一个空的 hosts 条目切片
	var currentEntry *HostsEntry // 使用指针，方便判断是否需要创建新条目
	var currentLine int = 0      // 跟踪当前行号

	// 处理新行，决定是否保存当前条目并创建新条目
	saveEntryAndReset := func() {
		// 如果存在当前条目且有有效内容，则保存
		if currentEntry != nil && (currentEntry.IP != "" || len(currentEntry.Domains) > 0 || currentEntry.IsCommentLine) {
			entries = append(entries, *currentEntry)
		}
		// 重置当前条目为nil，下次遇到有效token时会创建新的
		currentEntry = nil
	}

	for _, token := range tokens {
		// 当遇到新行时，保存之前的条目并重置
		if token.Line != currentLine {
			saveEntryAndReset()
			currentLine = token.Line
		}

		// 根据token类型处理
		switch token.Type {
		case lexer.COMMENT:
			// 如果行首是注释
			if currentEntry == nil {
				currentEntry = &HostsEntry{
					Line:          token.Line,
					Col:           token.Col,
					IsCommentLine: true,
					Comment:       token.Value,
				}
			} else {
				// 如果是行末注释
				currentEntry.HasComment = true
				currentEntry.Comment = token.Value
			}

		case lexer.IP:
			// 如果还没有创建当前条目，创建一个新的
			if currentEntry == nil {
				currentEntry = &HostsEntry{
					Line: token.Line,
					Col:  token.Col,
				}
			}
			// 设置IP
			currentEntry.IP = token.Value
			// 检查是否为IPv6
			currentEntry.IsIPv6 = token.Value == "::1" || token.Value[0] == ':' ||
				token.Value[0] == '[' || (len(token.Value) > 1 && token.Value[1] == ':')

		case lexer.DOMAIN:
			// 如果还没有创建当前条目，创建一个新的
			if currentEntry == nil {
				currentEntry = &HostsEntry{
					Line: token.Line,
					Col:  token.Col,
				}
			}
			// 添加域名到当前条目
			currentEntry.Domains = append(currentEntry.Domains, token.Value)

		case lexer.NEWLINE:
			// 处理空行
			if currentEntry == nil {
				currentEntry = &HostsEntry{
					Line:        token.Line,
					Col:         token.Col,
					IsEmptyLine: true,
				}
			}

		case lexer.EOF:
			// 到达文件末尾，确保保存最后一个条目
			saveEntryAndReset()
		}
	}

	return entries
}

// GetIPDomainPairs 从hosts条目中提取所有IP和域名对
// 每个IP地址对应一个或多个域名
func GetIPDomainPairs(entries []HostsEntry) []IPDomainPair {
	// 使用map先收集每个IP对应的所有域名
	ipMap := make(map[string][]string)

	for _, entry := range entries {
		// 跳过注释行、空行和没有域名的条目
		if entry.IsCommentLine || entry.IsEmptyLine || len(entry.Domains) == 0 || entry.IP == "" {
			continue
		}

		// 将域名添加到对应IP的列表中
		ipMap[entry.IP] = append(ipMap[entry.IP], entry.Domains...)
	}

	// 将map转换为IPDomainPair数组
	var pairs []IPDomainPair
	for ip, domains := range ipMap {
		pairs = append(pairs, IPDomainPair{
			IP:      ip,
			Domains: domains,
		})
	}

	return pairs
}

// FilterIPDomainPairs 根据IP或域名筛选IP-域名对
func FilterIPDomainPairs(pairs []IPDomainPair, ip string, domain string) []IPDomainPair {
	var results []IPDomainPair

	for _, pair := range pairs {
		// 如果指定了IP且不匹配，则跳过
		if ip != "" && pair.IP != ip {
			continue
		}

		// 如果指定了域名，筛选包含该域名的条目
		if domain != "" {
			containsDomain := false
			var matchedDomains []string

			for _, d := range pair.Domains {
				if d == domain {
					containsDomain = true
					matchedDomains = append(matchedDomains, d)
				}
			}

			if !containsDomain {
				continue // 不包含指定域名，跳过
			}

			// 如果仅筛选特定域名，则只保留匹配的域名
			if domain != "" {
				results = append(results, IPDomainPair{
					IP:      pair.IP,
					Domains: matchedDomains,
				})
			}
		} else {
			// 未指定域名筛选条件，保留所有域名
			results = append(results, pair)
		}
	}

	return results
}
