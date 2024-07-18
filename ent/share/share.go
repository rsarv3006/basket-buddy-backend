// Code generated by ent, DO NOT EDIT.

package share

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the share type in the database.
	Label = "share"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldExpiration holds the string denoting the expiration field in the database.
	FieldExpiration = "expiration"
	// FieldShareCode holds the string denoting the share_code field in the database.
	FieldShareCode = "share_code"
	// FieldData holds the string denoting the data field in the database.
	FieldData = "data"
	// FieldCreatorID holds the string denoting the creator_id field in the database.
	FieldCreatorID = "creator_id"
	// Table holds the table name of the share in the database.
	Table = "shares"
)

// Columns holds all SQL columns for share fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldExpiration,
	FieldShareCode,
	FieldData,
	FieldCreatorID,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultExpiration holds the default value on creation for the "expiration" field.
	DefaultExpiration time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the Share queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByExpiration orders the results by the expiration field.
func ByExpiration(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldExpiration, opts...).ToFunc()
}

// ByShareCode orders the results by the share_code field.
func ByShareCode(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldShareCode, opts...).ToFunc()
}

// ByCreatorID orders the results by the creator_id field.
func ByCreatorID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatorID, opts...).ToFunc()
}
