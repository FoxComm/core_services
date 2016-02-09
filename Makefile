# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test
GOFMT = $(GOCMD) fmt
GOGET = $(GOCMD) get
PKG_NAME = github.com/FoxComm/FoxComm
PSQL_CMD = psql -t template1 -c
PSQL_FOXCOMM = psql foxcomm -c
PSQL_FOXCOMM_TEST = psql foxcomm_test -c

# Lists of services and packages to perform operations against.
SERVICE_LIST := backups catalog_cache deployer feature_manager user ui
BUILD_LIST = $(foreach int, $(SERVICE_LIST), $(int)_build)
FMT_LIST = $(foreach int, $(SERVICE_LIST), $(int)_fmt)
GET_LIST = $(foreach int, $(SERVICE_LIST), $(int)_get)

# List actions
$(BUILD_LIST): %_build:
	cd $* && $(GOBUILD)

$(GET_LIST): %_get:
	cd $* && $(GOGET)

$(FMT_LIST): %_fmt:
	$(GOFMT) $(PKG_NAME)/$*/...

# Targets
.PHONY: all test clean setup gpm

dependencies:
	#brew install mercurial git bazaar
	cd $(GOPATH)
	go get -u github.com/pilu/fresh
	go get -u github.com/ddollar/forego
	#bundle install

configure: dependencies
	mkdir -p bin
	curl -s https://raw.githubusercontent.com/pote/gpm/master/bin/gpm > bin/gpm
	chmod +x bin/gpm

build: $(BUILD_LIST)

get: $(GET_LIST)

gpm:
	bin/gpm

etcd:
	cd deployer && ./key_injector_dev.sh -o

etcd-test:
	cd deployer && ./key_injector_dev.sh -o -t

node-setup:
	cd ui && npm install

db-setup:
	$(PSQL_CMD) 'drop database if exists foxcomm;'
	$(PSQL_CMD) 'drop database if exists foxcomm_user;'
	$(PSQL_CMD) 'drop database if exists fc_socialshopping;'
	$(PSQL_CMD) 'create database foxcomm;'
	$(PSQL_CMD) 'create database foxcomm_user;'
	$(PSQL_CMD) 'create database fc_socialshopping;'

db-test-setup:
	$(PSQL_CMD) 'drop database if exists foxcomm_test;'
	$(PSQL_CMD) 'drop database if exists foxcomm_user_test;'
	$(PSQL_CMD) 'drop database if exists fc_socialshopping_test;'
	$(PSQL_CMD) 'create database foxcomm_test;'
	$(PSQL_CMD) 'create database foxcomm_user_test;'
	$(PSQL_CMD) 'create database fc_socialshopping_test;'

seeds:
	StoreID=1 FC_ENV=development go run db/main.go migrations duh
	FC_ENV=development go run db/main.go seeds stores
	$(PSQL_FOXCOMM) "UPDATE store_features SET datasource = 'dbname=fc_socialshopping sslmode=disable' WHERE feature_id = (SELECT id FROM features WHERE features.name = 'social_shopping');"
	$(PSQL_FOXCOMM) "UPDATE store_features SET datasource = 'localhost#social_analytics' WHERE feature_id = (SELECT id FROM features WHERE features.name = 'social_analytics');"
	$(PSQL_FOXCOMM) "UPDATE store_features SET datasource = 'localhost#social_analytics' WHERE feature_id = (SELECT id FROM features WHERE features.name = 'loyalty_engine');"
	$(PSQL_FOXCOMM) "UPDATE store_features SET datasource='dbname=foxcomm_user sslmode=disable' WHERE feature_id = (SELECT id FROM features WHERE features.name = 'user');"
	StoreID=1 go run db/main.go migrations
	StoreID=1 go run db/main.go seeds all

test-seeds:
	StoreID=1 FC_ENV=test go run db/main.go migrations duh
	FC_ENV=test go run db/main.go seeds stores
	$(PSQL_FOXCOMM_TEST) "UPDATE store_features SET datasource = 'dbname=fc_socialshopping_test sslmode=disable' WHERE feature_id = (SELECT id FROM features WHERE features.name = 'social_shopping');"
	$(PSQL_FOXCOMM_TEST) "UPDATE store_features SET datasource = 'localhost#social_analytics_test' WHERE feature_id = (SELECT id FROM features WHERE features.name = 'social_analytics');"
	$(PSQL_FOXCOMM_TEST) "UPDATE store_features SET datasource = 'localhost#social_analytics_test' WHERE feature_id = (SELECT id FROM features WHERE features.name = 'loyalty_engine');"
	$(PSQL_FOXCOMM_TEST) "UPDATE store_features SET datasource='dbname=foxcomm_user_test sslmode=disable' WHERE feature_id = (SELECT id FROM features WHERE features.name = 'user');"
	StoreID=1 FC_ENV=test go run db/main.go migrations
	StoreID=1 FC_ENV=test go run db/main.go seeds all

setup: dependencies configure gpm node-setup db-setup etcd seeds

test-setup: etcd-test db-test-setup test-seeds

fmt: $(FMT_LIST)

test:
	FC_ENV=test LOG_LEVEL=2 go test -v ./... -cover
