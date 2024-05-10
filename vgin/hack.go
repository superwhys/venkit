package vgin

import (
	_ "unsafe"

	_ "github.com/gin-gonic/gin"
)

//go:linkname joinPaths github.com/gin-gonic/gin.joinPaths
func joinPaths(absolutePath, relativePath string) string
