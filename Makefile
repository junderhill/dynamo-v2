.PHONY: all
integration_tests:
	@command -v docker >/dev/null 2>&1 || { echo >&2 "I require docker but it's not installed.  Aborting."; exit 1; }
	@command -v aws >/dev/null 2>&1 || { echo >&2 "I require aws-cli but it's not installed.  Aborting."; exit 1; }

	docker run --name dynamov2-ddb -d -p 8000:8000 amazon/dynamodb-local

	docker ps

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
    --endpoint-url http://localhost:8000 # using DynamoDB local

	DYNAMO_TEST_TABLE=TestDB DYNAMO_TEST_REGION=eu-west-1 go test ./... -cover

	docker stop dynamov2-ddb
	docker rm dynamov2-ddb