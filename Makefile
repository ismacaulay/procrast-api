ORG=ismacaulay
APP=procrast-api

build:
	go build ./cmd/procrast-api

run:
	. env.sh && ./procrast-api

build-admin:
	go build ./cmd/admin

image:
	docker build -t $(ORG)/$(APP) -f Dockerfile .

image-db:
	docker build -t $(ORG)/procrastdb -f ./db/Dockerfile ./db

image-userdb:
	docker build -t $(ORG)/userdb -f ./db/Dockerfile.userdb ./db

