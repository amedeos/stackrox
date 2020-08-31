import React, { Component } from 'react';
import {
    FlexibleXYPlot,
    XAxis,
    YAxis,
    VerticalGridLines,
    HorizontalGridLines,
    VerticalBarSeries,
} from 'react-vis';
import PropTypes from 'prop-types';
import ReactRouterPropTypes from 'react-router-prop-types';
import { withRouter, Link } from 'react-router-dom';
import DiscreteColorLegend from 'react-vis/dist/legends/discrete-color-legend';
import merge from 'deepmerge';

import { standardBaseTypes } from 'constants/entityTypes';
import colors from 'constants/visuals/colors';
import { DetailedTooltipOverlay, HoverHint } from '@stackrox/ui-components';

class VerticalClusterBar extends Component {
    static propTypes = {
        id: PropTypes.string,
        history: ReactRouterPropTypes.history.isRequired,
        data: PropTypes.shape({}).isRequired,
        colors: PropTypes.arrayOf(PropTypes.string),
        containerProps: PropTypes.shape({}),
        plotProps: PropTypes.shape({}),
        seriesProps: PropTypes.shape({}),
        tickValues: PropTypes.arrayOf(PropTypes.number),
        tickFormat: PropTypes.func,
        labelLinks: PropTypes.shape({}),
    };

    static defaultProps = {
        id: '',
        colors,
        containerProps: {},
        plotProps: {},
        seriesProps: {},
        tickValues: [25, 50, 75, 100],
        tickFormat: (x) => `${x}%`,
        labelLinks: {},
    };

    constructor(props) {
        super(props);

        this.state = { hintInfo: null };
    }

    getLegendData = () => {
        const { data, colors: colorRange } = this.props;
        return Object.keys(data)
            .sort()
            .map((key, i) => ({
                title: standardBaseTypes[key],
                color: colorRange[i % colorRange.length],
            }));
    };

    render() {
        const { id, data, colors: colorRange, tickValues, tickFormat, labelLinks } = this.props;

        // Default props
        const defaultPlotProps = {
            xType: 'ordinal',
            yDomain: [0, 100],
            height: 270,
        };

        const defaultContainerProps = {
            className: 'relative chart-container w-full horizontal-bar-responsive',
        };

        const defaultSeriesProps = {
            // animation: true, //causes onValueMouseOut to fail https://github.com/uber/react-vis/issues/381
            barWidth: 0.5,
            style: {
                opacity: '.85',
                width: '10px',
                ry: '2px',
                cursor: 'pointer',
            },
            onValueMouseOver: (datum, e) => {
                this.setState({
                    hintInfo: { data: datum.hint, target: e.event.target },
                });
            },
            onValueMouseOut: () => {
                this.setState({ hintInfo: null });
            },
            onValueClick: (datum) => {
                if (datum.link) {
                    this.props.history.push(datum.link);
                }
            },
        };

        // Merge props
        const containerProps = merge(defaultContainerProps, this.props.containerProps);
        const plotProps = merge(defaultPlotProps, this.props.plotProps);
        const seriesProps = merge(defaultSeriesProps, this.props.seriesProps);

        function formatTicks(value) {
            let inner = value;
            if (labelLinks[value]) {
                inner = (
                    <Link
                        style={{ fill: 'currentColor' }}
                        className="underline text-base-600 hover:text-primary-700"
                        to={labelLinks[value]}
                    >
                        {value}
                    </Link>
                );
            }

            return <tspan>{inner}</tspan>;
        }

        // Calculate unique cluster names
        let clusterNames = new Set();
        Object.keys(data).forEach((dataSetKey) => {
            const dataSet = data[dataSetKey];
            dataSet.forEach((datum) => {
                clusterNames.add(datum.x);
            });
        });
        clusterNames = Array.from(clusterNames);

        // Create Barseries for each data set
        const series = [];
        Object.keys(data)
            .sort()
            .forEach((key, i) => {
                series.push(
                    <VerticalBarSeries
                        data={data[key]}
                        color={colorRange[i % colorRange.length]}
                        className={`vertical-cluster-bar-${standardBaseTypes[key]}`}
                        {...seriesProps}
                        key={key}
                    />
                );
            });

        const { hintInfo } = this.state;

        return (
            <div {...containerProps} data-testid={id}>
                <div className="flex flex-col h-full">
                    <FlexibleXYPlot {...plotProps}>
                        <VerticalGridLines
                            left={330 / clusterNames.length / 2 + 30}
                            tickValues={clusterNames.slice(0, clusterNames.length - 1)}
                        />
                        <HorizontalGridLines tickValues={tickValues} />
                        <YAxis tickValues={tickValues} tickSize={0} tickFormat={tickFormat} />
                        {series}

                        <XAxis tickSize={0} tickFormat={formatTicks} />
                    </FlexibleXYPlot>
                    <div>
                        <DiscreteColorLegend
                            orientation="horizontal"
                            items={this.getLegendData()}
                            colors={colors}
                            className="horizontal-bar-legend"
                        />
                    </div>
                    {hintInfo?.target && (
                        <HoverHint target={hintInfo.target}>
                            <DetailedTooltipOverlay
                                title={hintInfo.data.title}
                                body={hintInfo.data.body}
                            />
                        </HoverHint>
                    )}
                </div>
            </div>
        );
    }
}

export default withRouter(VerticalClusterBar);
