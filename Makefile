.PHONY: all
integration_tests: scaffold_up
	DYNAMO_TEST_TABLE=TestDB DYNAMO_TEST_REGION=eu-west-1 DYNAMO_TEST_ENDPOINT=http://localhost:8010 go test ./... -cover
	$(MAKE) scaffold_down

.PHONY: all
scaffold_up:
	@command -v docker >/dev/null 2>&1 || { echo >&2 "I require docker but it's not installed.  Aborting."; exit 1; }
	@command -v aws >/dev/null 2>&1 || { echo >&2 "I require aws-cli but it's not installed.  Aborting."; exit 1; }
	@command -v docker-compose >/dev/null 2>&1 || { echo >&2 "I require docker-compose but it's not installed.  Aborting."; exit 1; }

	docker-compose -p "helixddb" up -d

	docker ps
	$(MAKE) create_aws_table

.PHONY: all
scaffold_down:
	docker-compose -p "helixddb"  down

.PHONY: all
create_aws_table: 
	aws dynamodb create-table \
	--table-name TestDB \
	--attribute-definitions \
		AttributeName=UserID,AttributeType=N \
		AttributeName=Time,AttributeType=S \
		AttributeName=Msg,AttributeType=S \
	--key-schema \
		AttributeName=UserID,KeyType=HASH \
		AttributeName=Time,KeyType=RANGE \
	--global-secondary-indexes \
		IndexName=Msg-Time-index,KeySchema=[{'AttributeName=Msg,KeyType=HASH'},{'AttributeName=Time,KeyType=RANGE'}],Projection={'ProjectionType=ALL'} \
	--billing-mode PAY_PER_REQUEST \
	--region eu-west-1 \
	--endpoint-url http://localhost:8010 # using DynamoDB local
