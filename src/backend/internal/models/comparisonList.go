package models

type ComparisonListField int
type ComparisonListFieldsToUpdate map[ComparisonListField]any

const (
	ComparisonListFieldUserId = ComparisonListField(iota)
	ComparisonListFieldTotalPrice
	ComparisonListFieldAmount
)

type ComparisonList struct {
	ComparisonListId uint64
	UserId           uint64
	TotalPrice       uint64
	Amount           uint64
}
