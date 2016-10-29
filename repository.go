package usersvc

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Repository interface {
	GetUser(id string) (User, error)
	PutUser(u User) error
	DeleteUser(u User) error
}

type dynamoRepository struct {
	dynamodb *dynamodb.DynamoDB
	err      error
}

type inmemRepository struct {
	mtx sync.RWMutex
	m   map[string]User
}

type key struct {
	ID string `json:"id"`
}

func NewRepository() Repository {
	return &inmemRepository{
		m: map[string]User{},
	}
}

func NewDynamoRepository() Repository {
	s := session.New(&aws.Config{Region: aws.String("us-west-2")})

	return &dynamoRepository{
		dynamodb: dynamodb.New(s),
	}
}

func (r *inmemRepository) GetUser(id string) (User, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	u, ok := r.m[id]

	if !ok {
		return User{}, ErrNotFound
	}

	return u, nil
}

func (r *inmemRepository) PutUser(u User) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.m[u.ID] = u
	return nil
}

func (r *inmemRepository) DeleteUser(u User) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if _, ok := r.m[u.ID]; !ok {
		return ErrNotFound
	}
	delete(r.m, u.ID)
	return nil
}

func (r *dynamoRepository) GetUser(id string) (User, error) {
	r.err = nil

	k := r.marshalKey(id)
	resp := r.getItem(k)
	u := r.unmarshalUser(resp.Item)

	if r.err != nil {
		return User{}, r.err
	}

	return u, nil
}

func (r *dynamoRepository) PutUser(u User) error {
	r.err = nil

	item := r.marshalUser(u)
	r.putItem(item)

	if r.err != nil {
		return r.err
	}

	return nil
}

func (r *dynamoRepository) DeleteUser(u User) error {
	r.err = nil

	k := r.marshalKey(u.ID)
	r.deleteItem(k)

	if r.err != nil {
		return r.err
	}

	return nil
}

func (r *dynamoRepository) marshalKey(id string) map[string]*dynamodb.AttributeValue {
	if r.err != nil {
		return nil
	}

	k, e := dynamodbattribute.MarshalMap(key{ID: id})
	r.err = e

	return k
}

func (r *dynamoRepository) getItem(key map[string]*dynamodb.AttributeValue) *dynamodb.GetItemOutput {
	if r.err != nil {
		return nil
	}

	params := &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String("Users"),
	}

	resp, e := r.dynamodb.GetItem(params)

	if e != nil {
		r.err = e
		return nil
	}

	return resp
}

func (r *dynamoRepository) unmarshalUser(item map[string]*dynamodb.AttributeValue) User {
	if r.err != nil {
		return User{}
	}

	u := User{}
	e := dynamodbattribute.UnmarshalMap(item, &u)

	if e != nil {
		r.err = e
		return User{}
	}

	return u
}

func (r *dynamoRepository) marshalUser(u User) map[string]*dynamodb.AttributeValue {
	if r.err != nil {
		return nil
	}

	item, e := dynamodbattribute.MarshalMap(u)

	if e != nil {
		r.err = e
		return nil
	}

	return item
}

func (r *dynamoRepository) putItem(item map[string]*dynamodb.AttributeValue) {
	if r.err != nil {
		return
	}

	params := &dynamodb.PutItemInput{
		ConditionExpression: aws.String("attribute_not_exists(id)"),
		Item:                item,
		TableName:           aws.String("Users"),
	}

	_, e := r.dynamodb.PutItem(params)

	if e != nil {
		r.err = e
	}
}

func (r *dynamoRepository) deleteItem(key map[string]*dynamodb.AttributeValue) {
	if r.err != nil {
		return
	}

	params := &dynamodb.DeleteItemInput{
		Key:       key,
		TableName: aws.String("Users"),
	}

	_, e := r.dynamodb.DeleteItem(params)

	if e != nil {
		r.err = e
	}
}
