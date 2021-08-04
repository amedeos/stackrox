package role

import (
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/auth/permissions"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/uuid"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	// permissionSetIDPrefix should be prepended to every human-hostile ID of a
	// permission set for readability, e.g.,
	//     "io.stackrox.authz.permissionset.94ac7bfe-f9b2-402e-b4f2-bfda480e1a13".
	permissionSetIDPrefix = "io.stackrox.authz.permissionset."

	// accessScopeIDPrefix should be prepended to every human-hostile ID of an
	// access scope for readability, e.g.,
	//     "io.stackrox.authz.accessscope.94ac7bfe-f9b2-402e-b4f2-bfda480e1a13".
	accessScopeIDPrefix = "io.stackrox.authz.accessscope."
)

// GeneratePermissionSetID returns a random valid permission set ID.
func GeneratePermissionSetID() string {
	return permissionSetIDPrefix + uuid.NewV4().String()
}

// EnsureValidPermissionSetID converts id to the correct format if necessary.
func EnsureValidPermissionSetID(id string) string {
	if strings.HasPrefix(id, permissionSetIDPrefix) {
		return id
	}
	return permissionSetIDPrefix + id
}

// GenerateAccessScopeID returns a random valid access scope ID.
func GenerateAccessScopeID() string {
	return accessScopeIDPrefix + uuid.NewV4().String()
}

// EnsureValidAccessScopeID converts id to the correct format if necessary.
func EnsureValidAccessScopeID(id string) string {
	if strings.HasPrefix(id, accessScopeIDPrefix) {
		return id
	}
	return accessScopeIDPrefix + id
}

// ValidateRole checks whether the supplied protobuf message is a valid role.
func ValidateRole(role *storage.Role, permissionSetRequired bool) error {
	var multiErr error

	if role.GetName() == "" {
		err := errors.New("role name field must be set")
		multiErr = multierror.Append(multiErr, err)
	}
	if role.GetGlobalAccess() != storage.Access_NO_ACCESS {
		err := errors.Errorf("role name=%q: globalAccess should not be set, but is set to %s", role.GetName(), role.GetGlobalAccess())
		multiErr = multierror.Append(multiErr, err)
	}

	if permissionSetRequired {
		if len(role.GetResourceToAccess()) != 0 {
			err := errors.Errorf("role name=%q: must not have resourceToAccess, use a permission set instead", role.GetName())
			multiErr = multierror.Append(multiErr, err)
		}
		if role.GetPermissionSetId() == "" {
			err := errors.New("role permission_set_id field must be set")
			multiErr = multierror.Append(multiErr, err)
		}
		return multiErr
	}

	if role.GetPermissionSetId() != "" {
		err := errors.Errorf(
			"role name=%q: permission sets are not supported without the scoped access control feature", role.GetName())
		multiErr = multierror.Append(multiErr, err)
	}
	if role.GetAccessScopeId() != "" {
		err := errors.Errorf(
			"role name=%q: access scopes are not supported without the scoped access control feature", role.GetName())
		multiErr = multierror.Append(multiErr, err)
	}

	for resource := range role.GetResourceToAccess() {
		if _, ok := resources.MetadataForResource(permissions.Resource(resource)); !ok {
			multiErr = multierror.Append(multiErr, errors.Errorf(
				"role name=%q: resource %q does not exist", role.GetName(), resource))
		}
	}

	return multiErr
}

// ValidatePermissionSet checks whether the supplied protobuf message is a
// valid permission set.
func ValidatePermissionSet(ps *storage.PermissionSet) error {
	var multiErr error

	if !strings.HasPrefix(ps.GetId(), permissionSetIDPrefix) {
		multiErr = multierror.Append(multiErr, errors.Errorf("id field must be in '%s*' format", permissionSetIDPrefix))
	}
	if ps.GetName() == "" {
		multiErr = multierror.Append(multiErr, errors.New("name field must be set"))
	}
	for resource := range ps.GetResourceToAccess() {
		if _, ok := resources.MetadataForResource(permissions.Resource(resource)); !ok {
			multiErr = multierror.Append(multiErr, errors.Errorf(
				"resource %q does not exist", resource))
		}
	}

	return multiErr
}

// ValidateSimpleAccessScope checks whether the supplied protobuf message is a
// valid simple access scope.
func ValidateSimpleAccessScope(scope *storage.SimpleAccessScope) error {
	var multiErr error

	if !strings.HasPrefix(scope.GetId(), accessScopeIDPrefix) {
		multiErr = multierror.Append(multiErr, errors.Errorf("id field must be in '%s*' format", accessScopeIDPrefix))
	}
	if scope.GetName() == "" {
		multiErr = multierror.Append(multiErr, errors.New("name field must be set"))
	}

	err := ValidateSimpleAccessScopeRules(scope.GetRules())
	if err != nil {
		multiErr = multierror.Append(err)
	}

	return multiErr
}

// ValidateSimpleAccessScopeRules checks whether the supplied protobuf message
// represents valid simple access scope rules.
func ValidateSimpleAccessScopeRules(scopeRules *storage.SimpleAccessScope_Rules) error {
	var multiErr error

	for _, ns := range scopeRules.GetIncludedNamespaces() {
		if ns.GetClusterName() == "" || ns.GetNamespaceName() == "" {
			multiErr = multierror.Append(multiErr, errors.Errorf(
				"both cluster_name and namespace_name fields must be set in namespace rule <%s, %s>",
				ns.GetClusterName(), ns.GetNamespaceName()))
		}
	}
	for _, labelSelector := range scopeRules.GetClusterLabelSelectors() {
		err := validateSelectorRequirement(labelSelector)
		if err != nil {
			multiErr = multierror.Append(multiErr, err)
		}
	}
	for _, labelSelector := range scopeRules.GetNamespaceLabelSelectors() {
		err := validateSelectorRequirement(labelSelector)
		if err != nil {
			multiErr = multierror.Append(multiErr, err)
		}
	}

	return multiErr
}

func validateSelectorRequirement(labelSelector *storage.SetBasedLabelSelector) error {
	var multiErr error
	for _, requirement := range labelSelector.GetRequirements() {
		op := sac.ConvertLabelSelectorOperatorToSelectionOperator(requirement.GetOp())
		_, err := labels.NewRequirement(requirement.GetKey(), op, requirement.Values)
		if err != nil {
			multiErr = multierror.Append(multiErr, err)
		}
	}
	return multiErr
}
