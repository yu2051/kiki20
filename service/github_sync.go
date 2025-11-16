package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type GitHubSyncService struct {
	client       *github.Client
	owner        string
	repo         string
	enabled      bool
	syncInterval time.Duration
	mu           sync.RWMutex
	ctx          context.Context
	cancel       context.CancelFunc
}

var (
	ghSyncService     *GitHubSyncService
	ghSyncServiceOnce sync.Once
)

// InitGitHubSyncService 初始化 GitHub 同步服务
func InitGitHubSyncService() error {
	var initErr error
	ghSyncServiceOnce.Do(func() {
		token := common.GetEnvOrDefaultString("GITHUB_SYNC_TOKEN", "")
		repoURL := common.GetEnvOrDefaultString("GITHUB_SYNC_REPO", "")
		
		if token == "" || repoURL == "" {
			common.SysLog("GitHub sync not configured, skipping...")
			return
		}

		// 解析仓库 URL: https://github.com/owner/repo
		owner, repo, err := parseGitHubRepo(repoURL)
		if err != nil {
			initErr = fmt.Errorf("failed to parse GitHub repo URL: %v", err)
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		tc := oauth2.NewClient(ctx, ts)
		client := github.NewClient(tc)

		syncInterval := time.Duration(common.GetEnvOrDefault("GITHUB_SYNC_INTERVAL", 300)) * time.Second

		ghSyncService = &GitHubSyncService{
			client:       client,
			owner:        owner,
			repo:         repo,
			enabled:      true,
			syncInterval: syncInterval,
			ctx:          ctx,
			cancel:       cancel,
		}

		common.SysLog(fmt.Sprintf("GitHub sync service initialized: %s/%s, interval: %v", owner, repo, syncInterval))
	})

	return initErr
}

// parseGitHubRepo 解析 GitHub 仓库 URL
func parseGitHubRepo(repoURL string) (owner, repo string, err error) {
	// 支持格式: https://github.com/owner/repo 或 owner/repo
	repoURL = trimPrefix(repoURL, "https://github.com/")
	repoURL = trimPrefix(repoURL, "http://github.com/")
	repoURL = trimSuffix(repoURL, ".git")
	repoURL = trimSuffix(repoURL, "/")
	
	parts := splitString(repoURL, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid repo format, expected 'owner/repo'")
	}
	
	return parts[0], parts[1], nil
}

// GetGitHubSyncService 获取 GitHub 同步服务实例
func GetGitHubSyncService() *GitHubSyncService {
	return ghSyncService
}

// StartAutoSync 启动自动同步
func (s *GitHubSyncService) StartAutoSync() {
	if s == nil || !s.enabled {
		return
	}

	common.SysLog("Starting GitHub auto sync...")

	// 启动时立即加载一次
	if err := s.LoadFromGitHub(); err != nil {
		common.SysLog(fmt.Sprintf("Initial load from GitHub failed: %v", err))
	}

	// 启动定时同步
	go func() {
		ticker := time.NewTicker(s.syncInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := s.SyncToGitHub(); err != nil {
					common.SysLog(fmt.Sprintf("Auto sync to GitHub failed: %v", err))
				}
			case <-s.ctx.Done():
				common.SysLog("GitHub auto sync stopped")
				return
			}
		}
	}()
}

// Stop 停止同步服务
func (s *GitHubSyncService) Stop() {
	if s != nil && s.cancel != nil {
		s.cancel()
	}
}

// SyncToGitHub 同步数据到 GitHub
func (s *GitHubSyncService) SyncToGitHub() error {
	if s == nil || !s.enabled {
		return fmt.Errorf("GitHub sync service not enabled")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	common.SysLog("Syncing data to GitHub...")

	// 同步渠道数据
	if err := s.syncChannels(); err != nil {
		return fmt.Errorf("sync channels failed: %v", err)
	}

	// 同步用户数据
	if err := s.syncUsers(); err != nil {
		return fmt.Errorf("sync users failed: %v", err)
	}

	// 同步令牌数据
	if err := s.syncTokens(); err != nil {
		return fmt.Errorf("sync tokens failed: %v", err)
	}

	common.SysLog("Data synced to GitHub successfully")
	return nil
}

// LoadFromGitHub 从 GitHub 加载数据
func (s *GitHubSyncService) LoadFromGitHub() error {
	if s == nil || !s.enabled {
		return fmt.Errorf("GitHub sync service not enabled")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	common.SysLog("Loading data from GitHub...")

	// 加载渠道数据
	if err := s.loadChannels(); err != nil {
		common.SysLog(fmt.Sprintf("Load channels failed: %v", err))
	}

	// 加载用户数据
	if err := s.loadUsers(); err != nil {
		common.SysLog(fmt.Sprintf("Load users failed: %v", err))
	}

	// 加载令牌数据
	if err := s.loadTokens(); err != nil {
		common.SysLog(fmt.Sprintf("Load tokens failed: %v", err))
	}

	common.SysLog("Data loaded from GitHub successfully")
	return nil
}

// syncChannels 同步渠道到 GitHub
func (s *GitHubSyncService) syncChannels() error {
	channels, err := model.GetAllChannels(0, 0, true, false)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(channels, "", "  ")
	if err != nil {
		return err
	}

	return s.updateFile("channels.json", data, "Update channels data")
}

// syncUsers 同步用户到 GitHub
func (s *GitHubSyncService) syncUsers() error {
	var users []*model.User
	if err := model.DB.Find(&users).Error; err != nil {
		return err
	}

	// 清除敏感信息
	for _, user := range users {
		user.Password = ""
		user.AccessToken = nil
	}

	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}

	return s.updateFile("users.json", data, "Update users data")
}

// syncTokens 同步令牌到 GitHub
func (s *GitHubSyncService) syncTokens() error {
	var tokens []*model.Token
	if err := model.DB.Find(&tokens).Error; err != nil {
		return err
	}

	data, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		return err
	}

	return s.updateFile("tokens.json", data, "Update tokens data")
}

// loadChannels 从 GitHub 加载渠道
func (s *GitHubSyncService) loadChannels() error {
	data, err := s.getFile("channels.json")
	if err != nil {
		return err
	}

	var channels []*model.Channel
	if err := json.Unmarshal(data, &channels); err != nil {
		return err
	}

	// 这里可以选择性地更新数据库
	common.SysLog(fmt.Sprintf("Loaded %d channels from GitHub", len(channels)))
	return nil
}

// loadUsers 从 GitHub 加载用户
func (s *GitHubSyncService) loadUsers() error {
	data, err := s.getFile("users.json")
	if err != nil {
		return err
	}

	var users []*model.User
	if err := json.Unmarshal(data, &users); err != nil {
		return err
	}

	common.SysLog(fmt.Sprintf("Loaded %d users from GitHub", len(users)))
	return nil
}

// loadTokens 从 GitHub 加载令牌
func (s *GitHubSyncService) loadTokens() error {
	data, err := s.getFile("tokens.json")
	if err != nil {
		return err
	}

	var tokens []*model.Token
	if err := json.Unmarshal(data, &tokens); err != nil {
		return err
	}

	common.SysLog(fmt.Sprintf("Loaded %d tokens from GitHub", len(tokens)))
	return nil
}

// updateFile 更新或创建文件
func (s *GitHubSyncService) updateFile(path string, content []byte, message string) error {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	// 获取文件信息（如果存在）
	fileContent, _, resp, err := s.client.Repositories.GetContents(
		ctx, s.owner, s.repo, path, nil,
	)

	var sha *string
	if err == nil && fileContent != nil {
		sha = fileContent.SHA
	} else if resp != nil && resp.StatusCode != 404 {
		return fmt.Errorf("get file failed: %v", err)
	}

	// 创建或更新文件
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(message),
		Content: content,
		SHA:     sha,
	}

	_, _, err = s.client.Repositories.CreateFile(ctx, s.owner, s.repo, path, opts)
	if err != nil {
		return fmt.Errorf("update file failed: %v", err)
	}

	return nil
}

// getFile 获取文件内容
func (s *GitHubSyncService) getFile(path string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	fileContent, _, _, err := s.client.Repositories.GetContents(
		ctx, s.owner, s.repo, path, nil,
	)
	if err != nil {
		return nil, err
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return nil, err
	}

	return []byte(content), nil
}

// 辅助函数
func trimPrefix(s, prefix string) string {
	if len(s) >= len(prefix) && s[:len(prefix)] == prefix {
		return s[len(prefix):]
	}
	return s
}

func trimSuffix(s, suffix string) string {
	if len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix {
		return s[:len(s)-len(suffix)]
	}
	return s
}

func splitString(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}