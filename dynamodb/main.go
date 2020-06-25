package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func main() {
	// need env AWS_ACCESS_KEY_ID AWS_ACCESS_KEY
	// credentials.NewEnvCredentials()
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region:      aws.String("cn-north-1"),
			Credentials: credentials.NewSharedCredentials("", "default"),
		},
	}))

	svc := dynamodb.New(sess)
	ListTables(svc)
	CreateTable(svc)
	ListItems(svc)
	PutOne(svc)
	PutBatch(svc)
	ListItems(svc)
	DeleteTable(svc)
}

func ListTables(svc *dynamodb.DynamoDB) {
	input := &dynamodb.ListTablesInput{}

	fmt.Printf("Tables:\n")

	for {
		// Get the list of tables
		result, err := svc.ListTables(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeInternalServerError:
					fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return
		}

		for _, n := range result.TableNames {
			fmt.Println(*n)
		}

		// assign the last read tablename as the start for our next call to the ListTables function
		// the maximum number of table names returned in a call is 100 (default), which requires us to make
		// multiple calls to the ListTables function to retrieve all table names
		input.ExclusiveStartTableName = result.LastEvaluatedTableName

		if result.LastEvaluatedTableName == nil {
			break
		}
	}
}

func CreateTable(svc *dynamodb.DynamoDB) {

	_, err := svc.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String("Music"),
	})
	if err == nil {
		fmt.Printf("describe table: %s, error: %v\n", "Music", err)
		return
	}

	input := &dynamodb.CreateTableInput{
		TableName: aws.String("Music"),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				// A name for the attribute.
				//
				// AttributeName is a required field
				AttributeName: aws.String("Artist"),

				// The data type for the attribute, where:
				//
				//    * S - the attribute is of type String
				//
				//    * N - the attribute is of type Number
				//
				//    * B - the attribute is of type Binary
				//
				// AttributeType is a required field
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("SongTitle"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Artist"),

				// The role that this key attribute will assume:
				//
				//    * HASH - partition key
				//
				//    * RANGE - sort key
				//
				// The partition key of an item is also known as its hash attribute. The term
				// "hash attribute" derives from DynamoDB's usage of an internal hash function
				// to evenly distribute data items across partitions, based on their partition
				// key values.
				//
				// The sort key of an item is also known as its range attribute. The term "range
				// attribute" derives from the way DynamoDB stores items with the same partition
				// key physically close together, in sorted order by the sort key value.
				//
				// KeyType is a required field
				KeyType: aws.String("HASH"),
			},
			{
				AttributeName: aws.String("SongTitle"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{ // The maximum number of strongly consistent reads consumed per second before
			// DynamoDB returns a ThrottlingException. For more information, see Specifying
			// Read and Write Requirements (https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/WorkingWithTables.html#ProvisionedThroughput)
			// in the Amazon DynamoDB Developer Guide.
			//
			// If read/write capacity mode is PAY_PER_REQUEST the value is set to 0.
			//
			// ReadCapacityUnits is a required field
			ReadCapacityUnits: aws.Int64(10),

			// The maximum number of writes consumed per second before DynamoDB returns
			// a ThrottlingException. For more information, see Specifying Read and Write
			// Requirements (https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/WorkingWithTables.html#ProvisionedThroughput)
			// in the Amazon DynamoDB Developer Guide.
			//
			// If read/write capacity mode is PAY_PER_REQUEST the value is set to 0.
			//
			// WriteCapacityUnits is a required field
			WriteCapacityUnits: aws.Int64(5),
		},
	}

	_, err = svc.CreateTable(input)
	if err != nil {
		panic(err)
	}
}

type MusicItem struct {
	Artist    string
	SongTitle string
	Size      int64
	Data      []byte
}

func ListItems(svc *dynamodb.DynamoDB) {
	result, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String("Music"),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() != dynamodb.ErrCodeResourceNotFoundException {
				panic(fmt.Sprintf("code: %s, message: %s\n", awsErr.Code(), awsErr.Message()))
			} else {
				fmt.Printf("Nothing Found!\n")
				return
			}
		} else {
			panic(err)
		}
	}
	for _, i := range result.Items {
		var item MusicItem
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			panic(err)
		}
		fmt.Printf("item: %#v\n", item)
	}
}

func PutOne(svc *dynamodb.DynamoDB) {
	item := MusicItem{
		Artist:    "Jay Chou",
		SongTitle: "Hello World",
		Size:      6 * 1024 * 1024,
		Data:      []byte{'\xff', '\xab'},
	}
	data, err := dynamodbattribute.MarshalMap(&item)
	if err != nil {
		panic(err)
	}
	_, err = svc.PutItem(&dynamodb.PutItemInput{
		Item:      data,
		TableName: aws.String("Music"),
	})
	if err != nil {
		panic(err)
	}
}

func PutBatch(svc *dynamodb.DynamoDB) {
	_, err := svc.BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			"Music": {
				{
					// DeleteRequest: nil,

					// THIS SHOULD BE FAILED
					// panic: Batch Write. error: ValidationException: Supplied AttributeValue has more than
					// one datatypes set, must contain exactly one of the supported datatypes status code: 400,
					// request id: 4V7MHTV5KFMMV6ULERTFGHHPF3VV4KQNSO5AEMVJF66Q9ASUAAJG
					// PutRequest: &dynamodb.PutRequest{
					// 	Item: map[string]*dynamodb.AttributeValue{"Artist": {
					// 		B:    []byte("Bob"),
					// 		NULL: aws.Bool(true),
					// 	}},
					// },

					// THIS SHOULD BE FAILED
					// panic: Batch Write. error: ValidationException: The provided key element does not match
					// the schema status code: 400, request id: KSUN6C9GR0A1C0VP2PBL3M2TSJVV4KQNSO5AEMVJF66Q9ASUAAJG
					// PutRequest: &dynamodb.PutRequest{
					// 	Item: map[string]*dynamodb.AttributeValue{
					// 		"Artist": {
					// 			B: []byte("Bob"),
					// 		},
					// 		"SongTitle": {
					// 			N: aws.String("3.1415"),
					// 		},
					// 	},
					// },

					PutRequest: &dynamodb.PutRequest{
						Item: map[string]*dynamodb.AttributeValue{
							"Artist": {
								S: aws.String("Bob"),
							},
							"SongTitle": {
								S: aws.String("Bob's Song"),
							},
							"Size": {
								N: aws.String("60000000"),
							},
							"Data": {
								B: []byte("Content"),
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		panic(fmt.Sprintf("Batch Write. error: %v\n", err))
	}
}

func DeleteTable(svc *dynamodb.DynamoDB) {
	_, err := svc.DeleteTable(&dynamodb.DeleteTableInput{
		TableName: aws.String("Music"),
	})
	if err != nil {
		panic(err)
	}
}
