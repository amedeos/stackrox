// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"github.com/stackrox/rox/pkg/postgres"
)

var (
	// CreateTableProcessbaselinesStmt holds the create statement for table `Processbaselines`.
	CreateTableProcessbaselinesStmt = &postgres.CreateStmts{
		Table: `
               create table if not exists processbaselines (
                   Id varchar,
                   Key_DeploymentId varchar,
                   Key_ClusterId varchar,
                   Key_Namespace varchar,
                   serialized bytea,
                   PRIMARY KEY(Id)
               )
               `,
		Indexes:  []string{},
		Children: []*postgres.CreateStmts{},
	}
)
