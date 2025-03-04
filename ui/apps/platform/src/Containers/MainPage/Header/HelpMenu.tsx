import React, { ReactElement, useState } from 'react';
import { Link } from 'react-router-dom';
import {
    ApplicationLauncher,
    ApplicationLauncherGroup,
    ApplicationLauncherItem,
    ApplicationLauncherSeparator,
} from '@patternfly/react-core';
import { QuestionCircleIcon } from '@patternfly/react-icons';

import useMetadata from 'hooks/useMetadata';
import { apidocsPath } from 'routePaths';
import { getVersionedDocs } from 'utils/versioning';

function HelpMenu(): ReactElement {
    const { releaseBuild, version } = useMetadata();
    const [isHelpMenuOpen, setIsHelpMenuOpen] = useState(false);

    function onToggleHelpMenu() {
        setIsHelpMenuOpen(!isHelpMenuOpen);
    }

    // React requires key to render an item in an array of elements.
    const appLauncherItems = [
        <ApplicationLauncherGroup key="">
            <ApplicationLauncherItem
                component={
                    <Link className="pf-c-app-launcher__menu-item" to={apidocsPath}>
                        API Reference
                    </Link>
                }
            />
            {version && (
                <>
                    <ApplicationLauncherItem
                        href={getVersionedDocs(version)}
                        isExternal
                        target="_blank"
                        rel="noopener noreferrer"
                    >
                        Help Center
                    </ApplicationLauncherItem>
                    <ApplicationLauncherSeparator />
                    <ApplicationLauncherItem isDisabled>
                        {`v${version}${releaseBuild ? '' : ' [DEV BUILD]'}`}
                    </ApplicationLauncherItem>
                </>
            )}
        </ApplicationLauncherGroup>,
    ];

    return (
        <ApplicationLauncher
            aria-label="Help menu"
            isGrouped
            isOpen={isHelpMenuOpen}
            items={appLauncherItems}
            onToggle={onToggleHelpMenu}
            position="right"
            toggleIcon={<QuestionCircleIcon alt="" />}
            className="co-app-launcher"
            data-quickstart-id="qs-masthead-utilitymenu"
        />
    );
}

export default HelpMenu;
