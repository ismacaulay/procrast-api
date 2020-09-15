ORG=ismacaulay
APP=procrast-api

build:
	go build ./cmd/procrast-api

run:
	. env.sh && ./procrast-api

