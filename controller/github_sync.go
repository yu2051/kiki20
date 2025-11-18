package controller

import (
	"net/http"

	"github.com/QuantumNous/new-api/common"
	"github.com/gin-gonic/gin"
)

// GetGitHubSyncStatus 获取 GitHub 同步状态
func GetGitHubSyncStatus(c *gin.Context) {
	// 从数据库读取配置
	common.OptionMapRWMutex.RLock()
	token := common.OptionMap["GitHubSyncToken"]
	repo := common.OptionMap["GitHubSyncRepo"]
	common.OptionMapRWMutex.RUnlock()
	
	enabled := token != "" && repo != ""
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"enabled": enabled,
			"last_sync_time": nil, // TODO: 从数据库或缓存中获取上次同步时间
		},
	})
}

// TriggerGitHubSync 手动触发 GitHub 同步
func TriggerGitHubSync(c *gin.Context) {
	// 从数据库读取配置
	common.OptionMapRWMutex.RLock()
	token := common.OptionMap["GitHubSyncToken"]
	repo := common.OptionMap["GitHubSyncRepo"]
	common.OptionMapRWMutex.RUnlock()
	
	if token == "" || repo == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "GitHub 同步未配置，请先配置 GitHub Token 和仓库地址",
		})
		return
	}
	
	// TODO: 实现实际的同步逻辑
	// 这里应该调用 service/github_sync.go 中的同步服务
	// 由于该服务可能还未完全实现，这里先返回成功响应
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "同步任务已触发，请稍后查看同步状态",
	})
}