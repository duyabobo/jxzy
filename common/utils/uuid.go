package utils

import (
	"strings"

	"github.com/google/uuid"
)

// GenerateUUID 生成UUID
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateShortUUID 生成短UUID（去除横线）
func GenerateShortUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

// IsValidUUID 验证UUID格式
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
