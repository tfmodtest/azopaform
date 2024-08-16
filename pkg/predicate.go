package pkg

import (
	"context"
	"fmt"
	"strings"
)

type Predicate interface {
	Evaluate(r SingleRule, ctx context.Context) (result, rules string, err error)
}

var _ Predicate = EmptyPredicate{}

type EmptyPredicate struct{}

func (e EmptyPredicate) Evaluate(r SingleRule, ctx context.Context) (result, rules string, err error) {
	//The same as "where" case below without operator name/value, without adding conditions suffix/prefix
	fieldName := r.Field.(string)
	if string(fieldName[len(fieldName)-3:]) == "[*]" {
		fieldName, _, _ = FieldNameProcessor(fieldName[:len(fieldName)-3], ctx)
	} else {
		fieldName, _, _ = FieldNameProcessor(fieldName, ctx)
	}
	exp := count + "(" + fieldName + ")"
	result = strings.Join([]string{result, exp}, "")
	return result, "", nil
}

var _ Predicate = EqualsPredicate{}

type EqualsPredicate struct{}

func (e EqualsPredicate) Evaluate(r SingleRule, ctx context.Context) (result, rules string, err error) {
	fieldName, _, _ := FieldNameProcessor(r.Field, ctx)
	result = strings.Join([]string{fieldName, "==", fmt.Sprint(r.Operator.Value)}, " ")
	return result, "", nil
}

var _ Predicate = NotEqualsPredicate{}

type NotEqualsPredicate struct{}

func (n NotEqualsPredicate) Evaluate(r SingleRule, ctx context.Context) (result, rules string, err error) {
	fieldName, _, _ := FieldNameProcessor(r.Field, ctx)
	result = strings.Join([]string{fieldName, "!=", fmt.Sprint(r.Operator.Value)}, " ")
	return result, "", nil
}

var _ Predicate = ExistsPredicate{}

type ExistsPredicate struct{}

func (e ExistsPredicate) Evaluate(s SingleRule, ctx context.Context) (result, rules string, err error) {
	if strings.EqualFold(s.Operator.Value.(string), "true") {
		return s.Field.(string), "", nil
	}
	return strings.Join([]string{not, s.Field.(string)}, " "), "", nil
}

var _ Predicate = LikePredicate{}

type LikePredicate struct{}

func (l LikePredicate) Evaluate(s SingleRule, ctx context.Context) (result, rules string, err error) {
	fieldName, _, _ := FieldNameProcessor(s.Field, ctx)
	return strings.Join([]string{regexExp, "(", fmt.Sprint(s.Operator.Value), ",", fieldName, ")"}, ""), "", nil
}

var _ Predicate = NotLikePredicate{}

type NotLikePredicate struct{}

func (n NotLikePredicate) Evaluate(r SingleRule, ctx context.Context) (result, rules string, err error) {
	fieldName, _, _ := FieldNameProcessor(r.Field, ctx)
	return strings.Join([]string{not, " ", regexExp, "(", fmt.Sprint(r.Operator.Value), ",", fieldName, ")"}, ""), "", nil
}

type WherePredicate struct{}

func (w WherePredicate) Evaluate(singleRule SingleRule, ctx context.Context) (result, rules string, err error) {
	//fmt.Printf("here is a where case %+v\n", singleRule)
	var subNames []string
	fieldName, _, _ := FieldNameProcessor(singleRule.Field, ctx)
	switch singleRule.Operator.Value.(type) {
	case SingleRule:
		//fmt.Printf("here is a singlerule case %+v\n", singleRule.Operator.Value)
		operator := singleRule.Operator.Value.(SingleRule)
		operatorSet := RuleSet{
			Flag:        "allOf",
			SingleRules: []SingleRule{operator},
			RuleSets:    nil,
		}
		subsetNames, subRule, err := operatorSet.RuleSetReader("x", ctx)
		if err != nil {
			return "", "", err
		}
		subNames = subsetNames
		rules = subRule
	case RuleSet:
		operator := singleRule.Operator.Value.(RuleSet)
		subsetNames, subRule, err := operator.RuleSetReader(fieldName, ctx)
		if err != nil {
			return "", "", err
		}
		subNames = subsetNames
		rules = subRule
	}

	//fmt.Printf("The rules are %+v\n", rules)
	if string(fieldName[len(fieldName)-3:]) == "[*]" {
		fieldName = fieldName[:len(fieldName)-3] + "[x]"
	}
	exp := count + "(" + "{" + "x" + "|" + fieldName + ";" + subNames[0] + "}" + ")"
	result = exp
	return result, rules, nil
}

var _ Predicate = WherePredicate{}

type InPredicate struct{}

func (i InPredicate) Evaluate(r SingleRule, ctx context.Context) (result, rules string, err error) {
	fieldName, _, _ := FieldNameProcessor(r.Field, ctx)
	return strings.Join([]string{"some", fieldName, "in", SliceConstructor(r.Operator.Value)}, " "), "", nil
}

var _ Predicate = InPredicate{}

type NotInPredicate struct{}

func (n NotInPredicate) Evaluate(r SingleRule, ctx context.Context) (result, rules string, err error) {
	fieldName, _, _ := FieldNameProcessor(r.Field, ctx)
	return strings.Join([]string{"not", fieldName, "in", SliceConstructor(r.Operator.Value)}, " "), "", nil
}

var _ Predicate = NotInPredicate{}
