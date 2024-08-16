package pkg

import (
	"context"
	"fmt"
	"strings"
)

type OperationValue string
type OperationField string

func (o OperationValue) Rego(ctx context.Context) (string, error) {
	processed, _, err := FieldNameProcessor(string(o), ctx)
	return processed, err
}

func (o OperationField) Rego(ctx context.Context) (string, error) {
	processed, _, err := FieldNameProcessor(string(o), ctx)
	return processed, err
}

type Rego interface {
	Rego(ctx context.Context) (string, error)
}

type operation struct {
	Subject Rego
}

var _ Rego = EqualsOperation{}

type EqualsOperation struct {
	operation
	Value string
}

func (e EqualsOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := e.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{fieldName, "==", fmt.Sprint(e.Value)}, " "), nil
}

var _ Rego = NotEqualsOperation{}

type NotEqualsOperation struct {
	operation
	Value string
}

func (n NotEqualsOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := n.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{fieldName, "!=", fmt.Sprint(n.Value)}, " "), nil
}

var _ Rego = LikeOperation{}

type LikeOperation struct {
	operation
	Value string
}

func (l LikeOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := l.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{regexExp, "(", fmt.Sprint(l.Value), ",", fieldName, ")"}, ""), nil
}

var _ Rego = NotLikeOperation{}

type NotLikeOperation struct {
	operation
	Value string
}

func (n NotLikeOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := n.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{not, " ", regexExp, "(", fmt.Sprint(n.Value), ",", fieldName, ")"}, ""), nil
}

var _ Rego = MatchOperation{}

type MatchOperation struct {
	operation
	Value string
}

func (m MatchOperation) Rego(ctx context.Context) (string, error) {
	return "", fmt.Errorf("`match` condition is not supported, yet")
}

var _ Rego = MatchInsensitivelyOperation{}

type MatchInsensitivelyOperation struct {
	operation
	Value string
}

func (m MatchInsensitivelyOperation) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`matchInsensitively` condition is not supported, yet")
}

var _ Rego = NotMatchOperation{}

type NotMatchOperation struct {
	operation
	Value string
}

func (n NotMatchOperation) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`notMatch` condition is not supported, yet")
}

var _ Rego = NotMatchInsensitivelyOperation{}

type NotMatchInsensitivelyOperation struct {
	operation
	Value string
}

func (n NotMatchInsensitivelyOperation) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`notMatchInsensitively` condition is not supported, yet")
}

var _ Rego = ContainsOperation{}

type ContainsOperation struct {
	operation
	Value string
}

func (c ContainsOperation) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`contains` condition is not supported, yet")
}

var _ Rego = NotContainsOperation{}

type NotContainsOperation struct {
	operation
	Value string
}

func (n NotContainsOperation) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`notContains` condition is not supported, yet")
}

var _ Rego = InOperation{}

type InOperation struct {
	operation
	Values []string
}

func (i InOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := i.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{"some", fieldName, "in", SliceConstructor(fieldName)}, " "), nil
}

var _ Rego = NotInOperation{}

type NotInOperation struct {
	operation
	Values []string
}

func (n NotInOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := n.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{"not", fieldName, "in", SliceConstructor(fieldName)}, " "), nil
}

var _ Rego = ContainsKeyOperation{}

type ContainsKeyOperation struct {
	operation
	KeyName string
}

func (c ContainsKeyOperation) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`containsKey` condition is not supported, yet")
}

var _ Rego = NotContainsKeyOperation{}

type NotContainsKeyOperation struct {
	operation
	KeyName string
}

func (n NotContainsKeyOperation) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`notContainsKey` condition is not supported, yet")
}

var _ Rego = LessOperation{}

type LessOperation struct {
	operation
	Value any
}

func (l LessOperation) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`less` condition is not supported, yet")
}

var _ Rego = LessOrEqualsOperation{}

type LessOrEqualsOperation struct {
	operation
	Value any
}

func (l LessOrEqualsOperation) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`lessOrEquals` condition is not supported, yet")
}

var _ Rego = GreaterOperation{}

type GreaterOperation struct {
	operation
	Value any
}

func (g GreaterOperation) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`greater` condition is not supported, yet")
}

var _ Rego = GreaterOrEqualsOperation{}

type GreaterOrEqualsOperation struct {
	operation
	Value any
}

func (g GreaterOrEqualsOperation) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`greaterOrEquals` condition is not supported, yet")
}

var _ Rego = ExistsOperation{}

type ExistsOperation struct {
	operation
	Value bool
}

func (e ExistsOperation) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`exists` condition is not supported, yet")
}
