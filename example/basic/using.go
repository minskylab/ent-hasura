package basic

type OperatorName string

const EqualOperator OperatorName = "$eq"

type BoolExpression struct{}

type ColumnExp struct {
	BoolExpression
	ColumnName string
	Operator   OperatorName
	Value      string
}

type ColumnExpression interface {
	Name() string
	Operator() OperatorName
	Val() string
}

type ColumnEqualsExpr struct {
	ColumnName string
	Value      string
}

func (c *ColumnEqualsExpr) Name() string {
	return c.ColumnName
}

func (c *ColumnEqualsExpr) Operator() OperatorName {
	return EqualOperator
}

func (c *ColumnEqualsExpr) Val() string {
	return c.Value
}

type AndExpression []BoolExpression

type OrExpression []BoolExpression

type NotExpression BoolExpression

type ExistsExpression struct {
	Table string
	Where BoolExpression
}

// type ColumnEqualsExpr struct {
// 	BoolExpression

// 	ColumnName string
// 	Value      string
// }

func encodeBoolExpression(expr BoolExpression) string {
	return ""
}

func usingAPI() {
	// client := ent.NewClient()

	// // client.Note.Get

	// note.HasAuthorsWith(user.EmailEQ(""))

	// client.Note.Create()
	// client.User.Create().
	// fmt.Println(user.FieldID == "X-Hasura-User-Id")

	// expr1 := ColumnExp{
	// 	ColumnName: "id",
	// 	Operator:   EqualOperator,
	// 	Value:      "X-Hasura-User-Id",
	// }

	// expr2 := ColumnEqualsExpr{
	// 	ColumnName: "id",
	// 	Value:      "X-Hasura-User-Id",
	// }

	// encodeBoolExpression(expr1)

	// expr
}
