package cargo_query

import (
	"fmt"
	"strings"
)

type CargoQuery struct {
	Tables  []string `json:"tables"`
	Fields  []string `json:"fields"`
	Where   string   `json:"where"`
	JoinOn  string   `json:"join_on"`
	GroupBy string   `json:"group_by"`
	Having  string   `json:"having"`
	OrderBy string   `json:"order_by"`
	Offset  int      `json:"offset"`
	Limit   int      `json:"limit"`
}

func NewCargoQuery(tables []string, fields []string, where string, joinOn string, groupBy string, having string, orderBy string, offset int, limit int) *CargoQuery {
	return &CargoQuery{
		Tables:  tables,
		Fields:  fields,
		Where:   where,
		JoinOn:  joinOn,
		GroupBy: groupBy,
		Having:  having,
		OrderBy: orderBy,
		Offset:  offset,
		Limit:   limit,
	}
}

func (q *CargoQuery) ToQuery() string {
	query := "action=cargoquery"
	if len(q.Tables) > 0 {
		query += fmt.Sprintf("&tables=%s", strings.Join(q.Tables, ","))
	}
	if len(q.Fields) > 0 {
		query += fmt.Sprintf("&fields=%s", strings.Join(q.Fields, ","))
	}
	if q.Where != "" {
		query += fmt.Sprintf("&where=%s", q.Where)
	}
	if q.JoinOn != "" {
		query += fmt.Sprintf("&join_on=%s", q.JoinOn)
	}
	if q.GroupBy != "" {
		query += fmt.Sprintf("&group_by=%s", q.GroupBy)
	}
	if q.Having != "" {
		query += fmt.Sprintf("&having=%s", q.Having)
	}
	if q.OrderBy != "" {
		query += fmt.Sprintf("&order_by=%s", q.OrderBy)
	}
	if q.Offset != 0 {
		query += fmt.Sprintf("&offset=%d", q.Offset)
	}
	if q.Limit != 0 {
		query += fmt.Sprintf("&limit=%d", q.Limit)
	}
	return query
}
