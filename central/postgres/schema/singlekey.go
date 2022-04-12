// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"github.com/stackrox/rox/pkg/postgres"
)

var (
	// CreateTableSinglekeyStmt holds the create statement for table `Singlekey`.
	CreateTableSinglekeyStmt = &postgres.CreateStmts{
		Table: `
               create table if not exists singlekey (
                   Key varchar,
                   Name varchar UNIQUE,
                   StringSlice text[],
                   Bool bool,
                   Uint64 integer,
                   Int64 integer,
                   Float numeric,
                   Labels jsonb,
                   Timestamp timestamp,
                   Enum integer,
                   Enums int[],
                   serialized bytea,
                   PRIMARY KEY(Key)
               )
               `,
		Indexes: []string{
			"create index if not exists singlekey_Key on singlekey using hash(Key)",
		},
		Children: []*postgres.CreateStmts{},
	}
)
