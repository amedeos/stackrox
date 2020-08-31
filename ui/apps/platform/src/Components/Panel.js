import React from 'react';
import PropTypes from 'prop-types';

import { Tooltip, TooltipOverlay } from '@stackrox/ui-components';
import CloseButton from './CloseButton';

export const headerClassName = 'flex w-full min-h-14 border-b border-base-400';

const Panel = (props) => (
    <div
        className={`flex flex-col border-r border-base-400 overflow-auto w-full ${
            props.className
        } ${props.short ? '' : 'h-full'}`}
        data-testid={props.id}
    >
        <div className="flex-no-wrap">
            <div className={props.headerClassName}>
                {props.leftButtons && (
                    <div className="flex items-center pr-3 relative border-base-400 border-r hover:bg-primary-300 hover:border-primary-300">
                        {props.leftButtons}
                    </div>
                )}
                {props.headerTextComponent ? (
                    <div className="flex" data-testid={`${props.id}-header`}>
                        {props.headerTextComponent}
                    </div>
                ) : (
                    <div
                        className={`overflow-hidden mx-4 flex text-base-600 items-center tracking-wide leading-normal font-700 ${
                            props.isUpperCase ? 'uppercase' : 'capitalize'
                        }`}
                        data-testid={`${props.id}-header`}
                    >
                        <Tooltip content={<TooltipOverlay>{props.header}</TooltipOverlay>}>
                            <div className="line-clamp break-all">{props.header}</div>
                        </Tooltip>
                    </div>
                )}

                <div
                    className={`flex items-center justify-end relative flex-1 ${
                        props.onClose ? 'pl-3' : 'px-3'
                    }`}
                >
                    {props.headerComponents && props.headerComponents}
                    {props.onClose && (
                        <CloseButton
                            onClose={props.onClose}
                            className={props.closeButtonClassName}
                            iconColor={props.closeButtonIconColor}
                        />
                    )}
                </div>
            </div>
        </div>
        <div className={`h-full overflow-y-auto ${props.bodyClassName}`}>{props.children}</div>
    </div>
);

Panel.propTypes = {
    id: PropTypes.string,
    header: PropTypes.string,
    headerTextComponent: PropTypes.element,
    headerClassName: PropTypes.string,
    bodyClassName: PropTypes.string,
    className: PropTypes.string,
    children: PropTypes.node.isRequired,
    onClose: PropTypes.func,
    closeButtonClassName: PropTypes.string,
    closeButtonIconColor: PropTypes.string,
    headerComponents: PropTypes.element,
    leftButtons: PropTypes.node,
    isUpperCase: PropTypes.bool,
    short: PropTypes.bool,
};

Panel.defaultProps = {
    id: 'panel',
    header: ' ',
    headerTextComponent: null,
    headerClassName,
    bodyClassName: null,
    className: '',
    onClose: null,
    closeButtonClassName: 'border-base-400 border-l',
    closeButtonIconColor: '',
    headerComponents: null,
    leftButtons: null,
    isUpperCase: true,
    short: false,
};

export default Panel;
