package gorm

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/ragoncsa/todo/domain"
)

const TableName = "Task"

type TaskService struct {
	DynamoDbClient *dynamodb.Client
}

func makeKey(taskId string) (map[string]types.AttributeValue, error) {
	id, err := attributevalue.Marshal(taskId)
	if err != nil {
		return nil, err
	}
	return map[string]types.AttributeValue{"ID": id}, nil
}

func (s *TaskService) Task(id string) (*domain.Task, error) {
	var task domain.Task

	key, err := makeKey(id)
	if err != nil {
		return nil, err
	}
	response, err := s.DynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: key, TableName: aws.String(TableName),
	})
	if err != nil {
		log.Printf("Couldn't get info about %v. Here's why: %v\n", id, err)
	} else if response.Item == nil {
		log.Printf("Item %s does not exist.\n", id)
		return nil, errors.New("NOT FOUND")
	} else {
		err = attributevalue.UnmarshalMap(response.Item, &task)
		if err != nil {
			log.Printf("Couldn't unmarshal response. Here's why: %v\n", err)
		}
	}
	return &task, err
}

func (s *TaskService) CreateTask(t *domain.Task) error {
	id := uuid.New()
	t.ID = id.String()
	fmt.Printf("ID: %s\n", t.ID)
	item, err := attributevalue.MarshalMap(t)
	if err != nil {
		return err
	}
	_, err = s.DynamoDbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(TableName), Item: item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
	}
	return err
}

func (s *TaskService) Tasks() ([]*domain.Task, error) {
	var tasks []domain.Task
	var err error
	var response *dynamodb.ScanOutput
	if err != nil {
		log.Printf("Couldn't build expressions for scan. Here's why: %v\n", err)
	} else {
		response, err = s.DynamoDbClient.Scan(context.TODO(), &dynamodb.ScanInput{
			TableName: aws.String(TableName),
		})
		if err != nil {
			log.Printf("Couldn't scan for tasks. Here's why: %v\n", err)
		} else {
			err = attributevalue.UnmarshalListOfMaps(response.Items, &tasks)
			if err != nil {
				log.Printf("Couldn't unmarshal query response. Here's why: %v\n", err)
			}
		}
	}
	res := []*domain.Task{}
	for _, v := range tasks {
		fmt.Printf("TASK: %#v\n", v)
		temp := v // make a copy to add to the slice
		res = append(res, &temp)
	}
	return res, err
}

func (s *TaskService) DeleteTask(id string) error {
	key, err := makeKey(id)
	if err != nil {
		return err
	}
	output, err := s.DynamoDbClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(TableName), Key: key,
	})
	_ = output
	if err != nil {
		log.Printf("Couldn't delete %s from the table. Here's why: %v\n", id, err)
	}
	return err
}

func (s *TaskService) DeleteTasks() error {
	tasks, err := s.Tasks()
	if err != nil {
		return err
	}
	for _, t := range tasks {
		err = s.DeleteTask(t.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
