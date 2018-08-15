// Libraries
import React, {PureComponent} from 'react'
import {connect} from 'react-redux'
import _ from 'lodash'

// Components
import LineGraph from 'src/shared/components/LineGraph'
import GaugeChart from 'src/shared/components/GaugeChart'
import TableGraph from 'src/shared/components/TableGraph'
import SingleStat from 'src/shared/components/SingleStat'
import TimeSeries from 'src/shared/components/time_series/TimeSeries'

// Constants
import {emptyGraphCopy} from 'src/shared/copy/cell'
import {DEFAULT_TIME_FORMAT} from 'src/dashboards/constants'

// Utils
import {buildQueries} from 'src/utils/buildQueriesForLayouts'

// Actions
import {setHoverTime} from 'src/dashboards/actions/v2'

// Types
import {TimeRange, Template, CellQuery} from 'src/types'
import {V1View, ViewType} from 'src/types/v2/dashboards'

interface Props {
  link: string
  timeRange: TimeRange
  templates: Template[]
  cellID: string
  inView: boolean
  isInCEO: boolean
  timeFormat: string
  cellHeight: number
  autoRefresh: number
  manualRefresh: number
  options: V1View
  staticLegend: boolean
  onZoom: () => void
  editQueryStatus: () => void
  onSetResolution: () => void
  grabDataForDownload: () => void
  handleSetHoverTime: () => void
}

class RefreshingGraph extends PureComponent<Props> {
  public static defaultProps: Partial<Props> = {
    inView: true,
    manualRefresh: 0,
    staticLegend: false,
  }

  public render() {
    const {link, inView, templates} = this.props
    const {queries, type} = this.props.options

    if (!queries.length) {
      return (
        <div className="graph-empty">
          <p data-test="data-explorer-no-results">{emptyGraphCopy}</p>
        </div>
      )
    }

    return (
      <TimeSeries
        link={link}
        inView={inView}
        queries={this.queries}
        templates={templates}
      >
        {({timeSeries, loading}) => {
          switch (type) {
            case ViewType.SingleStat:
              return this.singleStat(timeSeries)
            case ViewType.Table:
              return this.table(timeSeries)
            case ViewType.Gauge:
              return this.gauge(timeSeries)
            default:
              return this.lineGraph(timeSeries, loading)
          }
        }}
      </TimeSeries>
    )
  }

  private singleStat = (data): JSX.Element => {
    const {cellHeight, manualRefresh} = this.props
    const {colors, decimalPlaces} = this.props.options

    return (
      <SingleStat
        data={data}
        colors={colors}
        prefix={this.prefix}
        suffix={this.suffix}
        lineGraph={false}
        key={manualRefresh}
        cellHeight={cellHeight}
        decimalPlaces={decimalPlaces}
      />
    )
  }

  private table = (data): JSX.Element => {
    const {manualRefresh, handleSetHoverTime, grabDataForDownload} = this.props

    const {
      colors,
      fieldOptions,
      tableOptions,
      decimalPlaces,
    } = this.props.options

    return (
      <TableGraph
        data={data}
        colors={colors}
        key={manualRefresh}
        tableOptions={tableOptions}
        fieldOptions={fieldOptions}
        decimalPlaces={decimalPlaces}
        timeFormat={DEFAULT_TIME_FORMAT}
        grabDataForDownload={grabDataForDownload}
        handleSetHoverTime={handleSetHoverTime}
      />
    )
  }

  private gauge = (data): JSX.Element => {
    const {cellID, cellHeight, manualRefresh} = this.props
    const {colors, decimalPlaces} = this.props.options

    return (
      <GaugeChart
        data={data}
        cellID={cellID}
        colors={colors}
        prefix={this.prefix}
        suffix={this.suffix}
        key={manualRefresh}
        cellHeight={cellHeight}
        decimalPlaces={decimalPlaces}
        resizerTopHeight={100}
      />
    )
  }

  private lineGraph = (data, loading): JSX.Element => {
    const {
      onZoom,
      cellID,
      timeRange,
      cellHeight,
      staticLegend,
      manualRefresh,
      handleSetHoverTime,
    } = this.props

    const {decimalPlaces, axes, type, colors, queries} = this.props.options

    return (
      <LineGraph
        data={data}
        type={type}
        axes={axes}
        cellID={cellID}
        colors={colors}
        onZoom={onZoom}
        queries={queries}
        loading={loading}
        key={manualRefresh}
        timeRange={timeRange}
        cellHeight={cellHeight}
        staticLegend={staticLegend}
        decimalPlaces={decimalPlaces}
        handleSetHoverTime={handleSetHoverTime}
      />
    )
  }

  private get queries(): CellQuery[] {
    const {timeRange, options} = this.props
    const {type} = options
    const queries = buildQueries(options.queries, timeRange)

    if (type === ViewType.SingleStat) {
      return [queries[0]]
    }

    if (type === ViewType.Gauge) {
      return [queries[0]]
    }

    return queries
  }

  private get prefix(): string {
    const {axes} = this.props.options

    return _.get(axes, 'y.prefix', '')
  }

  private get suffix(): string {
    const {axes} = this.props.options
    return _.get(axes, 'y.suffix', '')
  }
}

const mapStateToProps = ({annotations: {mode}, sources}) => {
  const source = sources.find(s => s.default)
  const link = source.links.query

  return {
    mode,
    link,
  }
}

const mdtp = {
  handleSetHoverTime: setHoverTime,
}

export default connect(mapStateToProps, mdtp)(RefreshingGraph)
