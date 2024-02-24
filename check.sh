#!/bin/bash


set -e


go run main.go
diff foo/example.go bar/example.go
