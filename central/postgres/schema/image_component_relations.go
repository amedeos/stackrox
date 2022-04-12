// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"github.com/stackrox/rox/pkg/postgres"
)

var (
	// CreateTableImageComponentRelationsStmt holds the create statement for table `ImageComponentRelations`.
	CreateTableImageComponentRelationsStmt = &postgres.CreateStmts{
		Table: `
               create table if not exists image_component_relations (
                   Id varchar,
                   Location varchar,
                   ImageId varchar,
                   ImageComponentId varchar,
                   serialized bytea,
                   PRIMARY KEY(Id, ImageId, ImageComponentId),
                   CONSTRAINT fk_parent_table_0 FOREIGN KEY (ImageId) REFERENCES images(Id) ON DELETE CASCADE
               )
               `,
		Indexes:  []string{},
		Children: []*postgres.CreateStmts{},
	}
)
