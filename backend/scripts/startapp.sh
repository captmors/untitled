#!/bin/bash

APP_NAME=$1

mkdir -p internal/$APP_NAME/mdl
touch internal/$APP_NAME/mdl/mdl.go

mkdir -p internal/$APP_NAME/repo
touch internal/$APP_NAME/repo/repo.go

mkdir -p internal/$APP_NAME/svc
touch internal/$APP_NAME/svc/svc.go

mkdir -p internal/$APP_NAME/urls
touch internal/$APP_NAME/urls/urls.go

mkdir -p internal/$APP_NAME/ctl
touch internal/$APP_NAME/ctl/ctl.go
