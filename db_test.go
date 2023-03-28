package helixddb

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

var (
	testDB    *DB
	testTable = "TestDB"
)

const offlineSkipMsg = "DYNAMO_TEST_REGION not set"

func init() {
	if region := os.Getenv("DYNAMO_TEST_REGION"); region != "" {
		var endpoint *string
		if dte := os.Getenv("DYNAMO_TEST_ENDPOINT"); dte != "" {
			endpoint = aws.String(dte)
		}

		endpointResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if endpoint != nil {
				fmt.Printf("Using endpoint: %s\n", *endpoint)
				return aws.Endpoint{
					URL:           *endpoint,
					SigningRegion: region,
				}, nil
			}
			return aws.Endpoint{}, &aws.EndpointNotFoundError{} //default to the default endpoint resolver
		})

		cfg, err := config.LoadDefaultConfig(context.Background(), func(o *config.LoadOptions) error {
			o.Region = region
			o.EndpointResolverWithOptions = endpointResolver

			if endpoint != nil {
				o.Credentials = credentials.NewStaticCredentialsProvider("fake", "fake", "fake")
			}
			return nil
		})

		if err != nil {
			log.Fatal(err)
		}
		testDB = New(cfg)
	}
	if table := os.Getenv("DYNAMO_TEST_TABLE"); table != "" {
		testTable = table
	}
}

// widget is the data structure used for integration tests
type widget struct {
	UserID int       // PK
	Time   time.Time // RK
	Msg    string
	Count  int
	Meta   map[string]string
	StrPtr *string `dynamo:",allowempty"`
}

func TestListTables(t *testing.T) {
	if testDB == nil {
		t.Skip(offlineSkipMsg)
	}

	tables, err := testDB.ListTables().All()
	if err != nil {
		t.Error(err)
		return
	}

	found := false
	for _, t := range tables {
		if t == testTable {
			found = true
			break
		}
	}

	if !found {
		t.Error("couldn't find testTable", testTable, "in:", tables)
	}
}
