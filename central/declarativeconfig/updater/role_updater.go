package updater

import (
	"context"

	"github.com/gogo/protobuf/proto"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/declarativeconfig/types"
	"github.com/stackrox/rox/central/declarativeconfig/utils"
	groupDataStore "github.com/stackrox/rox/central/group/datastore"
	roleDataStore "github.com/stackrox/rox/central/role/datastore"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/declarativeconfig"
	"github.com/stackrox/rox/pkg/errox"
	"github.com/stackrox/rox/pkg/integrationhealth"
	"github.com/stackrox/rox/pkg/set"
)

type roleUpdater struct {
	roleDS        roleDataStore.DataStore
	groupDS       groupDataStore.DataStore
	reporter      integrationhealth.Reporter
	idExtractor   types.IDExtractor
	nameExtractor types.NameExtractor
}

var _ ResourceUpdater = (*roleUpdater)(nil)

func newRoleUpdater(roleDatastore roleDataStore.DataStore, groupDatastore groupDataStore.DataStore, reporter integrationhealth.Reporter) ResourceUpdater {
	return &roleUpdater{
		roleDS:        roleDatastore,
		groupDS:       groupDatastore,
		reporter:      reporter,
		idExtractor:   types.UniversalIDExtractor(),
		nameExtractor: types.UniversalNameExtractor(),
	}
}

func (u *roleUpdater) Upsert(ctx context.Context, m proto.Message) error {
	role, ok := m.(*storage.Role)
	if !ok {
		return errox.InvariantViolation.Newf("wrong type passed to role updater: %T", role)
	}
	return u.roleDS.UpsertRole(ctx, role)
}

func (u *roleUpdater) DeleteResources(ctx context.Context, resourceIDsToSkip ...string) ([]string, error) {
	rolesToSkip := set.NewFrozenStringSet(resourceIDsToSkip...)

	roles, err := u.roleDS.GetRolesFiltered(ctx, func(role *storage.Role) bool {
		return declarativeconfig.IsDeclarativeOrigin(role) &&
			!rolesToSkip.Contains(role.GetName())
	})
	if err != nil {
		return nil, errors.Wrap(err, "retrieving declarative roles")
	}

	var roleDeletionErr *multierror.Error
	var roleNames []string
	for _, role := range roles {
		if err := u.roleDS.RemoveRole(ctx, role.GetName()); err != nil {
			roleDeletionErr = multierror.Append(roleDeletionErr, err)
			roleNames = append(roleNames, role.GetName())
			u.reporter.UpdateIntegrationHealthAsync(utils.IntegrationHealthForProtoMessage(role, "", err,
				u.idExtractor, u.nameExtractor))
			if errors.Is(err, errox.ReferencedByAnotherObject) {
				role.Traits.Origin = storage.Traits_DECLARATIVE_ORPHANED
				if err = u.roleDS.UpdateRole(ctx, role); err != nil {
					roleDeletionErr = multierror.Append(roleDeletionErr, errors.Wrap(err, "setting origin to orphaned"))
				}
			}
		}
	}
	return roleNames, roleDeletionErr.ErrorOrNil()
}
