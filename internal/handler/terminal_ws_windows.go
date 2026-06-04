//go:build windows

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RunCommandWS is Unix-only (PTY); on Windows use stream or run endpoints.
func (h *TerminalHandler) RunCommandWS(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "interactive WebSocket terminal is not supported on Windows; use POST /terminal/run/stream",
	})
}
