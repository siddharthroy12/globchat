package frontend

import (
	"embed"
)

//go:embed static/*
var Files embed.FS
