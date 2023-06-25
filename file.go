package main

import "embed"

var (
	//go:embed all:openapi
	openAPI embed.FS
)
