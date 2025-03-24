package hosts

import (
	"errors"
	"hostsparser/lexer"
	"hostsparser/logger"
	"hostsparser/parser"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type HostsFile struct {
	Entries  []parser.HostsEntry
	Filepath string
}

// NewHostsFile 创建一个空的HostsFile实例
func NewHostsFile() *HostsFile {
	return &HostsFile{
		Entries: []parser.HostsEntry{},
	}
}

// LoadHostsFile 从指定路径加载hosts文件
func LoadHostsFile(path string) (*HostsFile, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		logger.Errorf("获取绝对路径失败: %v", err)
		return nil, errors.New("获取绝对路径失败")
	}

	content, err := os.ReadFile(absPath)
	if err != nil {
		logger.Errorf("读取文件失败: %v", err)
		return nil, errors.New("读取文件失败")
	}

	// 使用词法分析和语法分析解析hosts文件内容
	tokens := lexer.Lex(string(content))
	entries := parser.ParseHosts(tokens)

	hostsFile := &HostsFile{
		Entries:  entries,
		Filepath: absPath,
	}

	return hostsFile, nil
}

// SaveHostsFile 保存hosts文件到指定路径
func (h *HostsFile) SaveHostsFile(path string) error {
	if path == "" {
		path = h.Filepath
	}

	if path == "" {
		logger.Error("未指定保存路径")
		return errors.New("未指定保存路径")
	}

	// 将条目转换为字符串
	content := h.String()

	// 写入文件
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		logger.Errorf("保存文件失败: %v", err)
		return errors.New("保存文件失败")
	}

	// 更新文件路径
	h.Filepath = path
	return nil
}

// String 返回hosts文件的字符串表示
func (h *HostsFile) String() string {
	var lines []string

	for _, entry := range h.Entries {
		if entry.IsCommentLine {
			// 保留完整注释行
			lines = append(lines, entry.Comment)
		} else if entry.IsEmptyLine {
			// 保留空行
			lines = append(lines, "")
		} else {
			// 构建hosts条目行
			line := entry.IP
			for _, domain := range entry.Domains {
				line += "\t" + domain
			}
			if entry.HasComment {
				line += "\t" + entry.Comment
			}
			lines = append(lines, line)
		}
	}

	return strings.Join(lines, "\n") + "\n"
}

// AddEntry 添加一个新条目
func (h *HostsFile) AddEntry(ip string, domains []string, comment string) {
	entry := parser.HostsEntry{
		IP:            ip,
		Domains:       domains,
		Comment:       comment,
		HasComment:    comment != "",
		IsEmptyLine:   false,
		IsCommentLine: false,
	}

	h.Entries = append(h.Entries, entry)
}

// SetHostIP 设置或更新指定主机名对应的IP地址
// 如果主机名已存在，则更新IP地址
// 如果主机名不存在，则添加新条目
func (h *HostsFile) SetHostIP(hostname, ip string) error {
	// 验证主机名
	if hostname == "" {
		logger.Error("主机名不能为空")
		return errors.New("主机名不能为空")
	}

	// 验证IP地址
	if !isValidIP(ip) {
		logger.Errorf("IP地址格式无效: %s", ip)
		return errors.New("IP地址格式无效")
	}

	// 查找该主机名是否已存在
	var found bool
	for i, entry := range h.Entries {
		if !entry.IsCommentLine && !entry.IsEmptyLine {
			for j, domain := range entry.Domains {
				if strings.EqualFold(domain, hostname) {
					// 如果已存在且IP相同，不做任何更改
					if entry.IP == ip {
						return nil
					}

					// 如果域名列表中只有这一个域名，直接更新IP
					if len(entry.Domains) == 1 {
						h.Entries[i].IP = ip
						found = true
					} else {
						// 如果有多个域名，从当前条目中移除该域名
						h.Entries[i].Domains = append(
							entry.Domains[:j],
							entry.Domains[j+1:]...,
						)
						// 并创建一个新条目
						h.AddEntry(ip, []string{hostname}, "")
						found = true
					}
					break
				}
			}
		}
		if found {
			break
		}
	}

	// 如果没有找到该主机名，添加新条目
	if !found {
		h.AddEntry(ip, []string{hostname}, "")
	}

	return nil
}

// RemoveHost 删除指定主机名的所有映射
func (h *HostsFile) RemoveHost(hostname string) (bool, error) {
	if hostname == "" {
		logger.Error("主机名不能为空")
		return false, errors.New("主机名不能为空")
	}

	var found bool
	var modified bool

	// 遍历条目查找主机名
	for i := 0; i < len(h.Entries); i++ {
		entry := h.Entries[i]
		if !entry.IsCommentLine && !entry.IsEmptyLine {
			for j, domain := range entry.Domains {
				if strings.EqualFold(domain, hostname) {
					found = true

					// 如果域名列表中只有这一个域名，删除整个条目
					if len(entry.Domains) == 1 {
						h.Entries = append(h.Entries[:i], h.Entries[i+1:]...)
						i-- // 因为删除了一个元素，需要调整索引
					} else {
						// 如果有多个域名，只删除该域名
						h.Entries[i].Domains = append(
							entry.Domains[:j],
							entry.Domains[j+1:]...,
						)
					}
					modified = true
					break
				}
			}
		}
	}

	if !found {
		return false, nil // 未找到主机名
	}

	return modified, nil // 返回是否进行了修改
}

// GetHostIP 获取指定主机名对应的IP地址
// 如果主机名不存在或对应多个IP，返回错误
func (h *HostsFile) GetHostIP(hostname string) (string, error) {
	if hostname == "" {
		logger.Error("主机名不能为空")
		return "", errors.New("主机名不能为空")
	}

	var ips []string

	// 查找该主机名对应的所有IP
	for _, entry := range h.Entries {
		if !entry.IsCommentLine && !entry.IsEmptyLine {
			for _, domain := range entry.Domains {
				if strings.EqualFold(domain, hostname) {
					ips = append(ips, entry.IP)
					break
				}
			}
		}
	}

	// 根据找到的IP数量返回结果
	switch len(ips) {
	case 0:
		logger.Errorf("未找到主机名 %s", hostname)
		return "", errors.New("未找到主机名")
	case 1:
		return ips[0], nil
	default:
		logger.Errorf("主机名 %s 对应多个IP: %s", hostname, strings.Join(ips, ", "))
		return "", errors.New("主机名对应多个IP")
	}
}

// GetHostsByIP 获取指定IP地址对应的所有主机名
func (h *HostsFile) GetHostsByIP(ip string) ([]string, error) {
	if ip == "" {
		logger.Error("IP地址不能为空")
		return nil, errors.New("IP地址不能为空")
	}

	var hostnames []string

	// 查找该IP对应的所有主机名
	for _, entry := range h.Entries {
		if !entry.IsCommentLine && !entry.IsEmptyLine && entry.IP == ip {
			hostnames = append(hostnames, entry.Domains...)
		}
	}

	return hostnames, nil
}

// SetMultipleHostIPs 批量设置多个主机名到同一IP
func (h *HostsFile) SetMultipleHostIPs(hostnames []string, ip string) error {
	if len(hostnames) == 0 {
		logger.Error("主机名列表不能为空")
		return errors.New("主机名列表不能为空")
	}

	// 验证IP地址
	if !isValidIP(ip) {
		logger.Errorf("IP地址格式无效: %s", ip)
		return errors.New("IP地址格式无效")
	}

	// 先删除所有指定的主机名
	for _, hostname := range hostnames {
		h.RemoveHost(hostname)
	}

	// 添加新条目
	h.AddEntry(ip, hostnames, "")

	return nil
}

// GetSystemHostsPath 获取系统hosts文件路径
func GetSystemHostsPath() string {
	if os.PathSeparator == '\\' { // Windows
		return filepath.Join(os.Getenv("SystemRoot"), "System32", "drivers", "etc", "hosts")
	}
	// Unix/Linux/MacOS
	return "/etc/hosts"
}

// ParseHostsContent 解析hosts文件内容，返回条目列表
func ParseHostsContent(content string) ([]parser.HostsEntry, error) {
	tokens := lexer.Lex(content)
	entries := parser.ParseHosts(tokens)
	return entries, nil
}

// 检查是否是有效的IP地址，使用net包进行更准确的验证
func isValidIP(ip string) bool {
	if ip == "" {
		return false
	}

	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}

// GetAllIPDomainPairs 获取所有IP和域名对
// 返回的每个IPDomainPair包含一个IP和与之关联的所有域名
func (h *HostsFile) GetAllIPDomainPairs() []parser.IPDomainPair {
	// 使用map收集每个IP对应的所有域名
	ipMap := make(map[string][]string)

	for _, entry := range h.Entries {
		// 跳过注释行、空行和没有域名的条目
		if entry.IsCommentLine || entry.IsEmptyLine || len(entry.Domains) == 0 || entry.IP == "" {
			continue
		}

		// 将域名添加到对应IP的列表中
		ipMap[entry.IP] = append(ipMap[entry.IP], entry.Domains...)
	}

	// 将map转换为IPDomainPair数组
	var pairs []parser.IPDomainPair
	for ip, domains := range ipMap {
		pairs = append(pairs, parser.IPDomainPair{
			IP:      ip,
			Domains: domains,
		})
	}

	return pairs
}

// FilterIPDomainPairs 根据IP或域名筛选IP-域名对
func (h *HostsFile) FilterIPDomainPairs(ip string, domain string) []parser.IPDomainPair {
	pairs := h.GetAllIPDomainPairs()
	var results []parser.IPDomainPair

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
				if strings.EqualFold(d, domain) {
					containsDomain = true
					matchedDomains = append(matchedDomains, d)
				}
			}

			if !containsDomain {
				continue // 不包含指定域名，跳过
			}

			// 只保留匹配的域名
			results = append(results, parser.IPDomainPair{
				IP:      pair.IP,
				Domains: matchedDomains,
			})
		} else {
			// 未指定域名筛选条件，保留所有域名
			results = append(results, pair)
		}
	}

	return results
}
