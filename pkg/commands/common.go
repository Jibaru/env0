package commands

import (
	"log"
	"os"

	"github.com/Jibaru/env0/pkg/client"
)

var logger = log.New(os.Stdout, "", 0)
var apiClient = client.New("")
