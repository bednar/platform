import {replace, RouterAction} from 'react-router-redux'
import _ from 'lodash'
import qs from 'qs'
import {Dispatch} from 'redux'

import {
  getDashboards as getDashboardsAJAX,
  updateDashboard as updateDashboardAJAX,
  updateDashboardCell as updateDashboardCellAJAX,
  deleteDashboardCell as deleteDashboardCellAJAX,
  createDashboard as createDashboardAJAX,
} from 'src/dashboards/apis'

import {notify} from 'src/shared/actions/notifications'
import {errorThrown} from 'src/shared/actions/errors'
import {stripPrefix} from 'src/utils/basepath'

import {
  cellDeleted,
  dashboardImportFailed,
  dashboardImported,
} from 'src/shared/copy/notifications'

import {getDeep} from 'src/utils/wrappers'

// Types
import {
  Dashboard,
  Cell,
  TimeRange,
  Template,
  TemplateValue,
  TemplateType,
} from 'src/types'

export enum ActionType {
  LoadDashboards = 'LOAD_DASHBOARDS',
  LoadDashboard = 'LOAD_DASHBOARD',
  SetDashboardTimeRange = 'SET_DASHBOARD_TIME_RANGE',
  SetDashboardZoomedTimeRange = 'SET_DASHBOARD_ZOOMED_TIME_RANGE',
  UpdateDashboard = 'UPDATE_DASHBOARD',
  CreateDashboard = 'CREATE_DASHBOARD',
  DeleteDashboard = 'DELETE_DASHBOARD',
  DeleteDashboardFailed = 'DELETE_DASHBOARD_FAILED',
  AddDashboardCell = 'ADD_DASHBOARD_CELL',
  DeleteDashboardCell = 'DELETE_DASHBOARD_CELL',
  SyncDashboardCell = 'SYNC_DASHBOARD_CELL',
  EditCellQueryStatus = 'EDIT_CELL_QUERY_STATUS',
  TemplateVariableLocalSelected = 'TEMPLATE_VARIABLE_LOCAL_SELECTED',
  UpdateTemplates = 'UPDATE_TEMPLATES',
  SetHoverTime = 'SET_HOVER_TIME',
  SetActiveCell = 'SET_ACTIVE_CELL',
  SetDashboardTimeV1 = 'SET_DASHBOARD_TIME_V1',
  RetainRangesDashboardTimeV1 = 'RETAIN_RANGES_DASHBOARD_TIME_V1',
}

interface LoadDashboardsAction {
  type: ActionType.LoadDashboards
  payload: {
    dashboards: Dashboard[]
    dashboardID: number
  }
}

interface LoadDashboardAction {
  type: ActionType.LoadDashboard
  payload: {
    dashboard: Dashboard
  }
}

interface RetainRangesDashTimeV1Action {
  type: ActionType.RetainRangesDashboardTimeV1
  payload: {
    dashboardIDs: number[]
  }
}

interface SetTimeRangeAction {
  type: ActionType.SetDashboardTimeRange
  payload: {
    timeRange: TimeRange
  }
}

interface SetZoomedTimeRangeAction {
  type: ActionType.SetDashboardZoomedTimeRange
  payload: {
    zoomedTimeRange: TimeRange
  }
}

interface UpdateDashboardAction {
  type: ActionType.UpdateDashboard
  payload: {
    dashboard: Dashboard
  }
}

interface CreateDashboardAction {
  type: ActionType.CreateDashboard
  payload: {
    dashboard: Dashboard
  }
}

interface DeleteDashboardAction {
  type: ActionType.DeleteDashboard
  payload: {
    dashboard: Dashboard
  }
}

interface DeleteDashboardFailedAction {
  type: ActionType.DeleteDashboardFailed
  payload: {
    dashboard: Dashboard
  }
}

interface SyncDashboardCellAction {
  type: ActionType.SyncDashboardCell
  payload: {
    dashboard: Dashboard
    cell: Cell
  }
}

interface AddDashboardCellAction {
  type: ActionType.AddDashboardCell
  payload: {
    dashboard: Dashboard
    cell: Cell
  }
}

interface DeleteDashboardCellAction {
  type: ActionType.DeleteDashboardCell
  payload: {
    dashboard: Dashboard
    cell: Cell
  }
}

interface EditCellQueryStatusAction {
  type: ActionType.EditCellQueryStatus
  payload: {
    queryID: string
    status: string
  }
}

interface TemplateVariableLocalSelectedAction {
  type: ActionType.TemplateVariableLocalSelected
  payload: {
    dashboardID: number
    templateID: string
    value: TemplateValue
  }
}

interface UpdateTemplatesAction {
  type: ActionType.UpdateTemplates
  payload: {
    templates: Template[]
  }
}

interface SetHoverTimeAction {
  type: ActionType.SetHoverTime
  payload: {
    hoverTime: string
  }
}

interface SetActiveCellAction {
  type: ActionType.SetActiveCell
  payload: {
    activeCellID: string
  }
}

interface SetDashTimeV1Action {
  type: ActionType.SetDashboardTimeV1
  payload: {
    dashboardID: number
    timeRange: TimeRange
  }
}

export type Action =
  | LoadDashboardsAction
  | LoadDashboardAction
  | RetainRangesDashTimeV1Action
  | SetTimeRangeAction
  | SetZoomedTimeRangeAction
  | UpdateDashboardAction
  | CreateDashboardAction
  | DeleteDashboardAction
  | DeleteDashboardFailedAction
  | SyncDashboardCellAction
  | AddDashboardCellAction
  | DeleteDashboardCellAction
  | EditCellQueryStatusAction
  | TemplateVariableLocalSelectedAction
  | UpdateTemplatesAction
  | SetHoverTimeAction
  | SetActiveCellAction
  | SetDashTimeV1Action

export const loadDashboards = (
  dashboards: Dashboard[],
  dashboardID?: number
): LoadDashboardsAction => ({
  type: ActionType.LoadDashboards,
  payload: {
    dashboards,
    dashboardID,
  },
})

export const loadDashboard = (dashboard: Dashboard): LoadDashboardAction => ({
  type: ActionType.LoadDashboard,
  payload: {dashboard},
})

export const setDashTimeV1 = (
  dashboardID: number,
  timeRange: TimeRange
): SetDashTimeV1Action => ({
  type: ActionType.SetDashboardTimeV1,
  payload: {dashboardID, timeRange},
})

export const retainRangesDashTimeV1 = (
  dashboardIDs: number[]
): RetainRangesDashTimeV1Action => ({
  type: ActionType.RetainRangesDashboardTimeV1,
  payload: {dashboardIDs},
})

export const setTimeRange = (timeRange: TimeRange): SetTimeRangeAction => ({
  type: ActionType.SetDashboardTimeRange,
  payload: {timeRange},
})

export const setZoomedTimeRange = (
  zoomedTimeRange: TimeRange
): SetZoomedTimeRangeAction => ({
  type: ActionType.SetDashboardZoomedTimeRange,
  payload: {zoomedTimeRange},
})

export const updateDashboard = (
  dashboard: Dashboard
): UpdateDashboardAction => ({
  type: ActionType.UpdateDashboard,
  payload: {dashboard},
})

export const createDashboard = (
  dashboard: Dashboard
): CreateDashboardAction => ({
  type: ActionType.CreateDashboard,
  payload: {dashboard},
})

export const deleteDashboard = (
  dashboard: Dashboard
): DeleteDashboardAction => ({
  type: ActionType.DeleteDashboard,
  payload: {dashboard},
})

export const deleteDashboardFailed = (
  dashboard: Dashboard
): DeleteDashboardFailedAction => ({
  type: ActionType.DeleteDashboardFailed,
  payload: {dashboard},
})

export const syncDashboardCell = (
  dashboard: Dashboard,
  cell: Cell
): SyncDashboardCellAction => ({
  type: ActionType.SyncDashboardCell,
  payload: {dashboard, cell},
})

export const addDashboardCell = (
  dashboard: Dashboard,
  cell: Cell
): AddDashboardCellAction => ({
  type: ActionType.AddDashboardCell,
  payload: {dashboard, cell},
})

export const deleteDashboardCell = (
  dashboard: Dashboard,
  cell: Cell
): DeleteDashboardCellAction => ({
  type: ActionType.DeleteDashboardCell,
  payload: {dashboard, cell},
})

export const editCellQueryStatus = (
  queryID: string,
  status: string
): EditCellQueryStatusAction => ({
  type: ActionType.EditCellQueryStatus,
  payload: {queryID, status},
})

export const templateVariableLocalSelected = (
  dashboardID: number,
  templateID: string,
  value: TemplateValue
): TemplateVariableLocalSelectedAction => ({
  type: ActionType.TemplateVariableLocalSelected,
  payload: {dashboardID, templateID, value},
})

export const updateTemplates = (
  templates: Template[]
): UpdateTemplatesAction => ({
  type: ActionType.UpdateTemplates,
  payload: {templates},
})

export const setHoverTime = (hoverTime: string): SetHoverTimeAction => ({
  type: ActionType.SetHoverTime,
  payload: {hoverTime},
})

export const setActiveCell = (activeCellID: string): SetActiveCellAction => ({
  type: ActionType.SetActiveCell,
  payload: {activeCellID},
})

export const updateQueryParams = (updatedQueryParams: object): RouterAction => {
  const {search, pathname} = window.location
  const strippedPathname = stripPrefix(pathname)

  const newQueryParams = _.pickBy(
    {
      ...qs.parse(search, {ignoreQueryPrefix: true}),
      ...updatedQueryParams,
    },
    v => !!v
  )

  const newSearch = qs.stringify(newQueryParams)
  const newLocation = {pathname: strippedPathname, search: `?${newSearch}`}

  return replace(newLocation)
}

const getDashboard = (state, dashboardId: number): Dashboard => {
  const dashboard = state.dashboardUI.dashboards.find(
    d => d.id === +dashboardId
  )

  if (!dashboard) {
    throw new Error(`Could not find dashboard with id '${dashboardId}'`)
  }

  return dashboard
}

// Thunkers

export const getDashboardsAsync = () => async (
  dispatch: Dispatch<Action>
): Promise<Dashboard[]> => {
  try {
    const {
      data: {dashboards},
    } = await getDashboardsAJAX()
    dispatch(loadDashboards(dashboards))
    return dashboards
  } catch (error) {
    console.error(error)
    dispatch(errorThrown(error))
  }
}

export const getChronografVersion = () => async (): Promise<string> => {
  try {
    return Promise.resolve('2.0')
  } catch (error) {
    console.error(error)
  }
}

const removeUnselectedTemplateValues = (dashboard: Dashboard): Template[] => {
  const templates = getDeep<Template[]>(dashboard, 'templates', []).map(
    template => {
      if (
        template.type === TemplateType.CSV ||
        template.type === TemplateType.Map
      ) {
        return template
      }

      const value = template.values.find(val => val.selected)
      const values = value ? [value] : []

      return {...template, values}
    }
  )
  return templates
}

export const putDashboard = (dashboard: Dashboard) => async (
  dispatch: Dispatch<Action>
): Promise<void> => {
  try {
    // save only selected template values to server
    const templatesWithOnlySelectedValues = removeUnselectedTemplateValues(
      dashboard
    )
    const {
      data: dashboardWithOnlySelectedTemplateValues,
    } = await updateDashboardAJAX({
      ...dashboard,
      templates: templatesWithOnlySelectedValues,
    })
    // save all template values to redux
    dispatch(
      updateDashboard({
        ...dashboardWithOnlySelectedTemplateValues,
        templates: dashboard.templates,
      })
    )
  } catch (error) {
    console.error(error)
    dispatch(errorThrown(error))
  }
}

export const putDashboardByID = (dashboardID: number) => async (
  dispatch: Dispatch<Action>,
  getState
): Promise<void> => {
  try {
    const dashboard = getDashboard(getState(), dashboardID)
    const templates = removeUnselectedTemplateValues(dashboard)
    await updateDashboardAJAX({...dashboard, templates})
  } catch (error) {
    console.error(error)
    dispatch(errorThrown(error))
  }
}

export const updateDashboardCell = (dashboard: Dashboard, cell: Cell) => async (
  dispatch: Dispatch<Action>
): Promise<void> => {
  try {
    const {data} = await updateDashboardCellAJAX(cell)
    dispatch(syncDashboardCell(dashboard, data))
  } catch (error) {
    console.error(error)
    dispatch(errorThrown(error))
  }
}

export const deleteDashboardCellAsync = (
  dashboard: Dashboard,
  cell: Cell
) => async (dispatch: Dispatch<Action>): Promise<void> => {
  try {
    await deleteDashboardCellAJAX(cell)
    dispatch(deleteDashboardCell(dashboard, cell))
    dispatch(notify(cellDeleted(cell.name)))
  } catch (error) {
    console.error(error)
    dispatch(errorThrown(error))
  }
}

export const importDashboardAsync = (dashboard: Dashboard) => async (
  dispatch: Dispatch<Action>
): Promise<void> => {
  try {
    // save only selected template values to server
    const templatesWithOnlySelectedValues = removeUnselectedTemplateValues(
      dashboard
    )

    const results = await createDashboardAJAX({
      ...dashboard,
      templates: templatesWithOnlySelectedValues,
    })

    const dashboardWithOnlySelectedTemplateValues = _.get(results, 'data')

    // save all template values to redux
    dispatch(
      createDashboard({
        ...dashboardWithOnlySelectedTemplateValues,
        templates: dashboard.templates,
      })
    )

    const {
      data: {dashboards},
    } = await getDashboardsAJAX()

    dispatch(loadDashboards(dashboards))

    dispatch(notify(dashboardImported(name)))
  } catch (error) {
    const errorMessage = _.get(
      error,
      'data.message',
      'Could not upload dashboard'
    )
    dispatch(notify(dashboardImportFailed('', errorMessage)))
    console.error(error)
    dispatch(errorThrown(error))
  }
}
