package vgin

import (
	_ "unsafe"

	_ "github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin/binding"
)

//go:linkname joinPaths github.com/gin-gonic/gin.joinPaths
func joinPaths(absolutePath, relativePath string) string

//go:linkname mapForm github.com/gin-gonic/gin/binding.mapForm
func mapForm(ptr any, form map[string][]string) error

//go:linkname mapHeader github.com/gin-gonic/gin/binding.mapHeader
func mapHeader(ptr any, form map[string][]string) error

//go:linkname mapURI github.com/gin-gonic/gin/binding.mapURI
func mapURI(ptr any, form map[string][]string) error

//go:linkname mapFormByTag github.com/gin-gonic/gin/binding.mapFormByTag
func mapFormByTag(ptr any, form map[string][]string, tag string) error
