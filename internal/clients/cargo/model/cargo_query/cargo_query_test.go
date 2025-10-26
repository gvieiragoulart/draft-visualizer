package cargo_query

import (
	"testing"
)

func TestCargoQuery_ToQuery(t *testing.T) {
	tests := []struct {
		name     string
		tables   []string
		fields   []string
		where    string
		joinOn   string
		groupBy  string
		having   string
		orderBy  string
		offset   int
		limit    int
		expected string
	}{
		{
			name:     "complete query with all parameters",
			tables:   []string{"test_table", "other_table"},
			fields:   []string{"id", "name"},
			where:    "id = 1",
			joinOn:   "other_table.id = test_table.id",
			groupBy:  "test_table.id",
			having:   "COUNT(*) > 0",
			orderBy:  "test_table.id",
			offset:   10,
			limit:    50,
			expected: "action=cargoquery&tables=test_table,other_table&fields=id,name&where=id = 1&join_on=other_table.id = test_table.id&group_by=test_table.id&having=COUNT(*) > 0&order_by=test_table.id&offset=10&limit=50",
		},
		{
			name:     "query with only where clause",
			tables:   []string{"users"},
			fields:   []string{"id", "email"},
			where:    "active = 1",
			joinOn:   "",
			groupBy:  "",
			having:   "",
			orderBy:  "",
			offset:   0,
			limit:    0,
			expected: "action=cargoquery&tables=users&fields=id,email&where=active = 1",
		},
		{
			name:     "query with pagination",
			tables:   []string{"products"},
			fields:   []string{"id", "name", "price"},
			where:    "",
			joinOn:   "",
			groupBy:  "",
			having:   "",
			orderBy:  "price DESC",
			offset:   10,
			limit:    20,
			expected: "action=cargoquery&tables=products&fields=id,name,price&order_by=price DESC&offset=10&limit=20",
		},
		{
			name:     "query with group by and having",
			tables:   []string{"orders"},
			fields:   []string{"customer_id", "COUNT(*)"},
			where:    "",
			joinOn:   "",
			groupBy:  "customer_id",
			having:   "COUNT(*) > 5",
			orderBy:  "",
			offset:   0,
			limit:    0,
			expected: "action=cargoquery&tables=orders&fields=customer_id,COUNT(*)&group_by=customer_id&having=COUNT(*) > 5",
		},
		{
			name:     "minimal query with only tables",
			tables:   []string{"articles"},
			fields:   []string{},
			where:    "",
			joinOn:   "",
			groupBy:  "",
			having:   "",
			orderBy:  "",
			offset:   0,
			limit:    0,
			expected: "action=cargoquery&tables=articles",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := NewCargoQuery(tt.tables, tt.fields, tt.where, tt.joinOn, tt.groupBy, tt.having, tt.orderBy, tt.offset, tt.limit)
			result := query.ToQuery()
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
