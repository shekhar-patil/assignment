package storage

import (
	"shekhar-patil/assignment/api/models"
	"sync"
)

var (
	PipelineData []models.PipelineRecord
	Mu           sync.Mutex
)

var ValidToken = "s3cr3t-token"
