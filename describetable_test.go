package helixddb

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

func TestDescribeTableLSIWithoutArn(t *testing.T) {
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
		LocalSecondaryIndexes: []types.LocalSecondaryIndexDescription{
			{
				IndexName: aws.String("test-index"),
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
				Projection: &types.Projection{
					ProjectionType: "ALL",
				},
				IndexSizeBytes: aws.Int64(123456789),
				ItemCount:      aws.Int64(1001),
			},
		},
	}

	desc := newDescription(tableDescription)

	if len(desc.LSI) != 1 {
		t.Error("wrong index count:", len(desc.LSI), "≠", 1)
	}

	if desc.LSI[0].Name != "test-index" {
		t.Error("wrong index name:", desc.LSI[0].Name, "≠", "test-index")
	}

}

func TestDescribeTableGSIWithoutArn(t *testing.T) {
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
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndexDescription{
			{
				IndexName:   aws.String("test-index"),
				IndexStatus: "ACTIVE",
				ProvisionedThroughput: &types.ProvisionedThroughputDescription{
					ReadCapacityUnits:      aws.Int64(100),
					WriteCapacityUnits:     aws.Int64(100),
					LastIncreaseDateTime:   aws.Time(time.Now()),
					LastDecreaseDateTime:   aws.Time(time.Now()),
					NumberOfDecreasesToday: aws.Int64(1),
				},
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
				Projection: &types.Projection{
					ProjectionType: "ALL",
				},
				IndexSizeBytes: aws.Int64(123456789),
				ItemCount:      aws.Int64(1001),
			},
		},
	}

	desc := newDescription(tableDescription)

	if len(desc.GSI) != 1 {
		t.Error("wrong index count:", len(desc.GSI), "≠", 1)
	}

	if desc.GSI[0].Name != "test-index" {
		t.Error("wrong index name:", desc.GSI[0].Name, "≠", "test-index")
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
