#!/usr/bin/env bash

env CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o goSafeChatBackend ./src