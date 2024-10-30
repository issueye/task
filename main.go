package main

import (
	"task/internal/initialize"
	"task/pages/home"

	"github.com/ying32/govcl/vcl"
)

func main() {
	initialize.Init()
	vcl.RunApp(&home.Frm_task)
}
