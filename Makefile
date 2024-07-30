AWS_LOCAL_ENV := AWS_ACCESS_KEY_ID='abc' AWS_SECRET_ACCESS_KEY='def' AWS_ENDPOINT_URL_DYNAMODB=http://localhost:4566 AWS_REGION=eu-west-2

AWSLOCAL := $(AWS_LOCAL_ENV) aws --no-cli-pager $@

.DEFAULT: help
help: ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

data_generation: ## Runs the data generator
	$(AWS_LOCAL_ENV) go run cmd/data-generator/main.go

run_api_local: ## Runs the data generator
	$(AWS_LOCAL_ENV) go run cmd/api/main.go

dev_env: ## Creates the local dev environment
	docker compose up -d localstack
	terraform init
	terraform apply -auto-approve
	# Duplicated code in reset_dev_env, could be reused
	$(AWSLOCAL) dynamodb list-tables
	$(AWSLOCAL) dynamodb describe-table --table-name cicd-audit

reset_dev_env: ## Resets the local dev environment
	terraform init
	terraform destroy -auto-approve
	terraform apply -auto-approve
	$(AWSLOCAL) dynamodb list-tables
	$(AWSLOCAL) dynamodb describe-table --table-name cicd-audit

delete_dev_env: ## Deletes the local dev environment
	docker compose down

run_all_docker_compose: dev_env ## Runs the whole environment in Docker compose
	docker compose up -d api

lint: ## Runs a linter over the codebase
	revive -formatter stylish -config revive.toml ./...
	goimports -local 'github.com/muzzapp/interviewtask' -w .