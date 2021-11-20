package repository

import (
	"fmt"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/allisson/go-env"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/allisson/hammer"
)

func init() {
	txdb.Register("pgx", "postgres", env.GetString("HAMMER_TEST_DATABASE_URL", ""))
}

type txnTestHelper struct {
	db                  *sqlx.DB
	topicRepo           *Topic
	subscriptionRepo    *Subscription
	messageRepo         *Message
	deliveryRepo        *Delivery
	deliveryAttemptRepo *DeliveryAttempt
}

func newTxnTestHelper() txnTestHelper {
	cName := fmt.Sprintf("connection_%d", time.Now().UnixNano())
	db, _ := sqlx.Open("pgx", cName)
	return txnTestHelper{
		db:                  db,
		topicRepo:           NewTopic(db),
		subscriptionRepo:    NewSubscription(db),
		messageRepo:         NewMessage(db),
		deliveryRepo:        NewDelivery(db),
		deliveryAttemptRepo: NewDeliveryAttempt(db),
	}
}

func findQuery(tableName string, findOptions hammer.FindOptions) (sql string, args []interface{}) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("*").From(tableName)

	// Pagination
	if findOptions.FindPagination != nil {
		sb.Limit(int(findOptions.FindPagination.Limit)).Offset(int(findOptions.FindPagination.Offset))
	}

	// Filters
	for _, findFilter := range findOptions.FindFilters {
		switch findFilter.Operator {
		case "=":
			sb.Where(sb.Equal(findFilter.FieldName, findFilter.Value))
		case "gt":
			sb.Where(sb.GreaterThan(findFilter.FieldName, findFilter.Value))
		case "gte":
			sb.Where(sb.GreaterEqualThan(findFilter.FieldName, findFilter.Value))
		case "lt":
			sb.Where(sb.LessThan(findFilter.FieldName, findFilter.Value))
		case "lte":
			sb.Where(sb.LessEqualThan(findFilter.FieldName, findFilter.Value))
		}
	}

	// Order by
	if findOptions.FindOrderBy != nil {
		sb.OrderBy(findOptions.FindOrderBy.FieldName)
		switch findOptions.FindOrderBy.Order {
		case "", "asc", "ASC":
			sb.Asc()
		case "desc", "DESC":
			sb.Desc()
		}
	}

	// Return result
	return sb.Build()
}

func insertQuery(tableName string, structValue interface{}) (string, []interface{}) {
	theStruct := sqlbuilder.NewStruct(structValue).For(sqlbuilder.PostgreSQL)
	ib := theStruct.InsertInto(tableName, structValue)
	return ib.Build()
}

func updateQuery(tableName string, id string, structValue interface{}) (string, []interface{}) {
	theStruct := sqlbuilder.NewStruct(structValue).For(sqlbuilder.PostgreSQL)
	ib := theStruct.Update(tableName, structValue)
	ib.Where(ib.Equal("id", id))
	return ib.Build()
}

func deleteQuery(tableName string, id string) (string, []interface{}) {
	sb := sqlbuilder.PostgreSQL.NewDeleteBuilder()
	sb.DeleteFrom(tableName).Where(sb.Equal("id", id))
	return sb.Build()
}

func rollback(tx *sqlx.Tx) {
	if err := tx.Rollback(); err != nil {
		zap.L().Error("unable-to-rollback", zap.Error(err))
	}
}
