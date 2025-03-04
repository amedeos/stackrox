import qs from 'qs';

import { vulnerabilitiesWorkloadCvesPath } from 'routePaths';
import { VulnerabilitySeverity, vulnerabilitySeverities } from 'types/cve.proto';
import { SearchFilter } from 'types/search';
import { getQueryString } from 'utils/queryStringUtils';
import { searchValueAsArray } from 'utils/searchUtils';
import { ensureExhaustive } from 'utils/type.utils';

import { CveStatusTab, isValidCveStatusTab, QuerySearchFilter } from './types';

export type EntityTab = 'CVE' | 'Image' | 'Deployment';

export type WorkloadCvesSearch = {
    cveStatusTab: CveStatusTab;
    entityTab?: EntityTab;
    s?: SearchFilter;
};

export function parseWorkloadCvesOverviewSearchString(search: string): WorkloadCvesSearch {
    const { cveStatusTab } = qs.parse(search, { ignoreQueryPrefix: true });

    return {
        cveStatusTab: isValidCveStatusTab(cveStatusTab) ? cveStatusTab : 'Observed',
    };
}

export function getOverviewCvesPath(workloadCvesSearch: WorkloadCvesSearch): string {
    return `${vulnerabilitiesWorkloadCvesPath}${getQueryString(workloadCvesSearch)}`;
}

export function getEntityPagePath(workloadCveEntity: EntityTab, id: string): string {
    switch (workloadCveEntity) {
        case 'CVE':
            return `${vulnerabilitiesWorkloadCvesPath}/cves/${id}`;
        case 'Image':
            return `${vulnerabilitiesWorkloadCvesPath}/images/${id}`;
        case 'Deployment':
            return `${vulnerabilitiesWorkloadCvesPath}/deployments/${id}`;
        default:
            return ensureExhaustive(workloadCveEntity);
    }
}

/**
 * Parses an open `SearchFilter` obtained from the URL into a restricted `SearchFilter` that
 * matches the fields and values expected by the backend.
 */
export function parseQuerySearchFilter(rawSearchFilter: SearchFilter): QuerySearchFilter {
    const cleanSearchFilter: QuerySearchFilter = {};

    // SearchFilter values that can be directly translated over to the backend equivalent
    const unprocessedSearchKeys = ['IMAGE', 'DEPLOYMENT', 'NAMESPACE', 'CLUSTER'] as const;
    unprocessedSearchKeys.forEach((key) => {
        if (rawSearchFilter[key]) {
            cleanSearchFilter[key] = searchValueAsArray(rawSearchFilter[key]);
        }
    });

    if (rawSearchFilter.Fixable) {
        const rawFixable = searchValueAsArray(rawSearchFilter.Fixable);
        const cleanFixable: ('true' | 'false')[] = [];

        rawFixable.forEach((status) => {
            if (status === 'Fixable') {
                cleanFixable.push('true');
            } else if (status === 'Not fixable') {
                cleanFixable.push('false');
            }
        });

        cleanSearchFilter.Fixable = cleanFixable;
    }

    if (rawSearchFilter.Severity) {
        const rawSeverities = searchValueAsArray(rawSearchFilter.Severity);
        cleanSearchFilter.Severity = [];

        rawSeverities.forEach((rs) => {
            if (rs === 'Critical') {
                cleanSearchFilter.Severity?.push('CRITICAL_VULNERABILITY_SEVERITY');
            } else if (rs === 'Important') {
                cleanSearchFilter.Severity?.push('IMPORTANT_VULNERABILITY_SEVERITY');
            } else if (rs === 'Moderate') {
                cleanSearchFilter.Severity?.push('MODERATE_VULNERABILITY_SEVERITY');
            } else if (rs === 'Low') {
                cleanSearchFilter.Severity?.push('LOW_VULNERABILITY_SEVERITY');
            }
        });
    }

    return cleanSearchFilter;
}

// Given a search filter, determine which severities should be hidden from the user
export function getHiddenSeverities(
    querySearchFilter: QuerySearchFilter
): Set<VulnerabilitySeverity> {
    return querySearchFilter.Severity
        ? new Set(vulnerabilitySeverities.filter((s) => !querySearchFilter.Severity?.includes(s)))
        : new Set([]);
}
