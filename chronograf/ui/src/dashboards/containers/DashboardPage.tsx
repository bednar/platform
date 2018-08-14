// Libraries
import React, {Component, MouseEvent} from 'react'
import {connect} from 'react-redux'
import {withRouter} from 'react-router'

// Components
import {ErrorHandling} from 'src/shared/decorators/errors'
import DashboardHeader from 'src/dashboards/components/DashboardHeader'
import DashboardComponent from 'src/dashboards/components/Dashboard'
import ManualRefresh from 'src/shared/components/ManualRefresh'

// Actions
import * as dashboardActions from 'src/dashboards/actions/v2'
import * as appActions from 'src/shared/actions/app'
import * as errorActions from 'src/shared/actions/errors'
import * as notifyActions from 'src/shared/actions/notifications'

// Utils
import {getDeep} from 'src/utils/wrappers'
import {updateDashboardLinks} from 'src/dashboards/utils/dashboardSwitcherLinks'
import AutoRefresh from 'src/utils/AutoRefresh'

// APIs
import {loadDashboardLinks} from 'src/dashboards/apis/v2'

// Constants
import {
  interval,
  DASHBOARD_LAYOUT_ROW_HEIGHT,
  TEMP_VAR_DASHBOARD_TIME,
  TEMP_VAR_UPPER_DASHBOARD_TIME,
} from 'src/shared/constants'
import {FORMAT_INFLUXQL, defaultTimeRange} from 'src/shared/data/timeRanges'
import {EMPTY_LINKS} from 'src/dashboards/constants/dashboardHeader'

// Types
import {WithRouterProps} from 'react-router'
import {ManualRefreshProps} from 'src/shared/components/ManualRefresh'
import {Location} from 'history'
import {InjectedRouter} from 'react-router'
import * as AppActions from 'src/types/actions/app'
import * as ColorsModels from 'src/types/colors'
import * as ErrorsActions from 'src/types/actions/errors'
import * as QueriesModels from 'src/types/queries'
import * as NotificationsActions from 'src/types/actions/notifications'
import * as DashboardsModels from 'src/types/dashboards'

import {Links, Source} from 'src/types/v2'
import {Dashboard, Cell} from 'src/types/v2/dashboards'

interface Props extends ManualRefreshProps, WithRouterProps {
  links: Links
  source: Source
  sources: Source[]
  params: {
    dashboardID: string
  }
  location: Location
  dashboard: Dashboard
  handleChooseAutoRefresh: AppActions.SetAutoRefreshActionCreator
  autoRefresh: number
  timeRange: QueriesModels.TimeRange
  zoomedTimeRange: QueriesModels.TimeRange
  showTemplateControlBar: boolean
  inPresentationMode: boolean
  handleClickPresentationButton: AppActions.DelayEnablePresentationModeDispatcher
  cellQueryStatus: {
    queryID: string
    status: object
  }
  errorThrown: ErrorsActions.ErrorThrownActionCreator
  router: InjectedRouter
  notify: NotificationsActions.PublishNotificationActionCreator
  selectedCell: DashboardsModels.Cell
  thresholdsListType: string
  thresholdsListColors: ColorsModels.ColorNumber[]
  gaugeColors: ColorsModels.ColorNumber[]
  lineColors: ColorsModels.ColorString[]
  addCell: typeof dashboardActions.addCellAsync
  getDashboard: typeof dashboardActions.getDashboardAsync
  setDashTimeV1: typeof dashboardActions.setDashTimeV1
  setZoomedTimeRange: typeof dashboardActions.setZoomedTimeRange
  updateDashboard: typeof dashboardActions.updateDashboardAsync
  updateCells: typeof dashboardActions.updateCellsAsync
  putDashboard: typeof dashboardActions.putDashboard
  addDashboardCellAsync: typeof dashboardActions.addDashboardCellAsync
  editCellQueryStatus: typeof dashboardActions.editCellQueryStatus
  updateDashboardCell: typeof dashboardActions.updateDashboardCell
  cloneDashboardCellAsync: typeof dashboardActions.cloneDashboardCellAsync
  deleteDashboardCellAsync: typeof dashboardActions.deleteDashboardCellAsync
  updateQueryParams: typeof dashboardActions.updateQueryParams
}

interface State {
  scrollTop: number
  windowHeight: number
  selectedCell: DashboardsModels.Cell | null
  dashboardLinks: DashboardsModels.DashboardSwitcherLinks
}

@ErrorHandling
class DashboardPage extends Component<Props, State> {
  public constructor(props: Props) {
    super(props)

    this.state = {
      scrollTop: 0,
      selectedCell: null,
      windowHeight: window.innerHeight,
      dashboardLinks: EMPTY_LINKS,
    }
  }

  public async componentDidMount() {
    const {autoRefresh} = this.props

    AutoRefresh.poll(autoRefresh)

    window.addEventListener('resize', this.handleWindowResize, true)

    await this.getDashboard()

    this.getDashboardLinks()
  }

  public componentDidUpdate(prevProps: Props) {
    const {autoRefresh} = this.props

    const prevPath = getDeep(prevProps.location, 'pathname', null)
    const thisPath = getDeep(this.props.location, 'pathname', null)

    if (prevPath && thisPath && prevPath !== thisPath) {
      this.getDashboard()
    }

    if (autoRefresh !== prevProps.autoRefresh) {
      AutoRefresh.poll(autoRefresh)
    }
  }

  public componentWillUnmount() {
    AutoRefresh.stopPolling()

    window.removeEventListener('resize', this.handleWindowResize, true)
  }

  public render() {
    const {
      source,
      sources,
      timeRange,
      timeRange: {lower, upper},
      zoomedTimeRange,
      zoomedTimeRange: {lower: zoomedLower, upper: zoomedUpper},
      showTemplateControlBar,
      dashboard,
      autoRefresh,
      manualRefresh,
      onManualRefresh,
      inPresentationMode,
      handleChooseAutoRefresh,
      handleClickPresentationButton,
    } = this.props
    const low = zoomedLower || lower
    const up = zoomedUpper || upper

    const lowerType = low && low.includes(':') ? 'timeStamp' : 'constant'
    const upperType = up && up.includes(':') ? 'timeStamp' : 'constant'
    const dashboardTime = {
      id: 'dashtime',
      tempVar: TEMP_VAR_DASHBOARD_TIME,
      type: lowerType,
      values: [
        {
          value: low,
          type: lowerType,
          selected: true,
          localSelected: true,
        },
      ],
    }

    const upperDashboardTime = {
      id: 'upperdashtime',
      tempVar: TEMP_VAR_UPPER_DASHBOARD_TIME,
      type: upperType,
      values: [
        {
          value: up || 'now()',
          type: upperType,
          selected: true,
          localSelected: true,
        },
      ],
    }

    let templatesIncludingDashTime
    if (dashboard) {
      templatesIncludingDashTime = [dashboardTime, upperDashboardTime, interval]
    } else {
      templatesIncludingDashTime = []
    }

    const {dashboardLinks} = this.state

    return (
      <div className="page dashboard-page">
        <DashboardHeader
          dashboard={dashboard}
          timeRange={timeRange}
          autoRefresh={autoRefresh}
          isHidden={inPresentationMode}
          onAddCell={this.handleAddCell}
          onManualRefresh={onManualRefresh}
          zoomedTimeRange={zoomedTimeRange}
          onRenameDashboard={this.handleRenameDashboard}
          dashboardLinks={dashboardLinks}
          activeDashboard={dashboard ? dashboard.name : ''}
          showTemplateControlBar={showTemplateControlBar}
          handleChooseAutoRefresh={handleChooseAutoRefresh}
          handleChooseTimeRange={this.handleChooseTimeRange}
          handleClickPresentationButton={handleClickPresentationButton}
        />
        {!!dashboard && (
          <DashboardComponent
            source={source}
            sources={sources}
            setScrollTop={this.setScrollTop}
            inView={this.inView}
            dashboard={dashboard}
            timeRange={timeRange}
            autoRefresh={autoRefresh}
            manualRefresh={manualRefresh}
            onZoom={this.handleZoomedTimeRange}
            inPresentationMode={inPresentationMode}
            onPositionChange={this.handlePositionChange}
            onDeleteCell={this.handleDeleteDashboardCell}
            onCloneCell={this.handleCloneCell}
            templatesIncludingDashTime={templatesIncludingDashTime}
          />
        )}
      </div>
    )
  }

  private handleWindowResize = (): void => {
    this.setState({windowHeight: window.innerHeight})
  }

  private getDashboard = async () => {
    const {params, getDashboard} = this.props

    await getDashboard(params.dashboardID)
    this.updateActiveDashboard()
  }

  private updateActiveDashboard(): void {
    this.setState((prevState, props) => ({
      dashboardLinks: updateDashboardLinks(
        prevState.dashboardLinks,
        props.dashboard
      ),
    }))
  }

  private inView = (cell: DashboardsModels.Cell): boolean => {
    const {scrollTop, windowHeight} = this.state
    const bufferValue = 600
    const cellTop = cell.y * DASHBOARD_LAYOUT_ROW_HEIGHT
    const cellBottom = (cell.y + cell.h) * DASHBOARD_LAYOUT_ROW_HEIGHT
    const bufferedWindowBottom = windowHeight + scrollTop + bufferValue
    const bufferedWindowTop = scrollTop - bufferValue
    const topInView = cellTop < bufferedWindowBottom
    const bottomInView = cellBottom > bufferedWindowTop

    return topInView && bottomInView
  }

  private handleChooseTimeRange = (
    timeRange: QueriesModels.TimeRange
  ): void => {
    const {dashboard, setDashTimeV1, updateQueryParams} = this.props

    setDashTimeV1(dashboard.id, {
      ...timeRange,
      format: FORMAT_INFLUXQL,
    })

    updateQueryParams({
      lower: timeRange.lower,
      upper: timeRange.upper,
    })
  }

  private handlePositionChange = async (cells: Cell[]): Promise<void> => {
    const {dashboard, updateCells} = this.props
    await updateCells(dashboard, cells)
  }

  private handleAddCell = async (): Promise<void> => {
    const {dashboard, addCell} = this.props
    await addCell(dashboard)
  }

  private handleCloneCell = (cell: DashboardsModels.Cell): void => {
    const {dashboard} = this.props
    this.props.cloneDashboardCellAsync(dashboard, cell)
  }

  private handleRenameDashboard = async (name: string): Promise<void> => {
    const {dashboard} = this.props
    const renamedDashboard = {...dashboard, name}

    this.props.updateDashboard(renamedDashboard)
    this.updateActiveDashboard()
  }

  private handleDeleteDashboardCell = (cell: Cell): void => {
    const {dashboard} = this.props
    this.props.deleteDashboardCellAsync(dashboard, cell)
  }

  private handleZoomedTimeRange = (
    zoomedTimeRange: QueriesModels.TimeRange
  ): void => {
    const {setZoomedTimeRange, updateQueryParams} = this.props

    setZoomedTimeRange(zoomedTimeRange)

    updateQueryParams({
      zoomedLower: zoomedTimeRange.lower,
      zoomedUpper: zoomedTimeRange.upper,
    })
  }

  private setScrollTop = (e: MouseEvent<JSX.Element>): void => {
    const target = e.target as HTMLElement

    this.setState({scrollTop: target.scrollTop})
  }

  private getDashboardLinks = async (): Promise<void> => {
    const {links, dashboard: activeDashboard} = this.props

    try {
      const dashboardLinks = await loadDashboardLinks(
        links.dashboards,
        activeDashboard
      )

      this.setState({
        dashboardLinks,
      })
    } catch (error) {
      console.error(error)
    }
  }
}

const mstp = (state, {params: {dashboardID}}) => {
  const {
    links,
    app: {
      ephemeral: {inPresentationMode},
      persisted: {autoRefresh, showTemplateControlBar},
    },
    dashboardUI: {cellQueryStatus, zoomedTimeRange},
    sources,
    ranges,
    cellEditorOverlay: {
      cell,
      thresholdsListType,
      thresholdsListColors,
      gaugeColors,
      lineColors,
    },
    dashboards,
  } = state

  const timeRange =
    ranges.find(r => r.dashboardID === dashboardID) || defaultTimeRange

  const selectedCell = cell

  const dashboard = dashboards.find(d => d.id === dashboardID)

  return {
    links,
    sources,
    timeRange,
    dashboard,
    zoomedTimeRange,
    autoRefresh,
    cellQueryStatus,
    inPresentationMode,
    showTemplateControlBar,
    selectedCell,
    thresholdsListType,
    thresholdsListColors,
    gaugeColors,
    lineColors,
  }
}

const mdtp: Partial<Props> = {
  getDashboard: dashboardActions.getDashboardAsync,
  setDashTimeV1: dashboardActions.setDashTimeV1,
  setZoomedTimeRange: dashboardActions.setZoomedTimeRange,
  updateDashboard: dashboardActions.updateDashboardAsync,
  addCell: dashboardActions.addCellAsync,
  updateCells: dashboardActions.updateCellsAsync,
  editCellQueryStatus: dashboardActions.editCellQueryStatus,
  updateDashboardCell: dashboardActions.updateDashboardCell,
  cloneDashboardCellAsync: dashboardActions.cloneDashboardCellAsync,
  deleteDashboardCellAsync: dashboardActions.deleteDashboardCellAsync,
  updateQueryParams: dashboardActions.updateQueryParams,
  handleChooseAutoRefresh: appActions.setAutoRefresh,
  handleClickPresentationButton: appActions.delayEnablePresentationMode,
  errorThrown: errorActions.errorThrown,
  notify: notifyActions.notify,
}

export default connect(mstp, mdtp)(
  ManualRefresh<Props>(withRouter<Props>(DashboardPage))
)
