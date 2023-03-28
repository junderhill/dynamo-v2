package helixddb

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

func TestDescribeTable(t *testing.T) {
	if testDB == nil {
		t.Skip(offlineSkipMsg)
	}
	table := testDB.Table(testTable)

	desc, err := table.Describe().Run()
	if err != nil {
		t.Error(err)
		return
	}

	if desc.Name != testTable {
		t.Error("wrong name:", desc.Name, "≠", testTable)
	}
	if desc.HashKey != "UserID" || desc.RangeKey != "Time" {
		t.Error("bad keys:", desc.HashKey, desc.RangeKey)
	}
}

func TestDescribeTableDescription(t *testing.T) {
	tableDescription := &types.TableDescription{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("PK"),
				AttributeType: "S",
			},
			{
				AttributeName: aws.String("SK"),
				AttributeType: "S",
			},
		},
		TableName: aws.String("test"),
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
				KeyType:       "HASH",
			},
			{
				AttributeName: aws.String("SK"),
				KeyType:       "RANGE",
			},
		},
		TableStatus:      "ACTIVE",
		CreationDateTime: aws.Time(time.Now()),
		ItemCount:        aws.Int64(1001),
		TableSizeBytes:   aws.Int64(987654321),
	}

	desc := newDescription(tableDescription)

	if desc.Name != "test" {
		t.Error("wrong name:", desc.Name, "≠", "test")
	}

	if desc.HashKey != "PK" || desc.RangeKey != "SK" {
		t.Error("bad keys:", desc.HashKey, desc.RangeKey)
	}

	if desc.Items != 1001 {
		t.Error("bad item count:", desc.Items, "≠", 1001)
	}

	if desc.Size != 987654321 {
		t.Error("bad size:", desc.Size, "≠", 987654321)
	}
}
