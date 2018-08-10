// Types
import {Dispatch} from 'redux'
import {Dashboard} from 'src/types/v2'
import {replace} from 'react-router-redux'

// APIs
import {
  getDashboard as getDashboardAJAX,
  getDashboards as getDashboardsAJAX,
  createDashboard as createDashboardAJAX,
  deleteDashboard as deleteDashboardAJAX,
} from 'src/dashboards/apis/v2'

// Actions
import {notify} from 'src/shared/actions/notifications'
import {
  deleteTimeRange,
  updateTimeRangeFromQueryParams,
} from 'src/dashboards/actions/v2/ranges'

// Constants
import * as copy from 'src/shared/copy/notifications'

export enum ActionTypes {
  LoadDashboards = 'LOAD_DASHBOARDS',
  LoadDashboard = 'LOAD_DASHBOARD',
  DeleteDashboard = 'DELETE_DASHBOARD',
  DeleteDashboardFailed = 'DELETE_DASHBOARD_FAILED',
}

export type Action =
  | LoadDashboardsAction
  | DeleteDashboardAction
  | LoadDashboardAction

interface LoadDashboardsAction {
  type: ActionTypes.LoadDashboards
  payload: {
    dashboards: Dashboard[]
  }
}

interface DeleteDashboardAction {
  type: ActionTypes.DeleteDashboard
  payload: {
    dashboardID: string
  }
}

interface DeleteDashboardFailedAction {
  type: ActionTypes.DeleteDashboardFailed
  payload: {
    dashboard: Dashboard
  }
}

interface LoadDashboardAction {
  type: ActionTypes.LoadDashboard
  payload: {
    dashboard: Dashboard
  }
}

// Action Creators

export const loadDashboards = (
  dashboards: Dashboard[]
): LoadDashboardsAction => ({
  type: ActionTypes.LoadDashboards,
  payload: {
    dashboards,
  },
})

export const loadDashboard = (dashboard: Dashboard): LoadDashboardAction => ({
  type: ActionTypes.LoadDashboard,
  payload: {dashboard},
})

export const deleteDashboard = (
  dashboardID: string
): DeleteDashboardAction => ({
  type: ActionTypes.DeleteDashboard,
  payload: {dashboardID},
})

export const deleteDashboardFailed = (
  dashboard: Dashboard
): DeleteDashboardFailedAction => ({
  type: ActionTypes.DeleteDashboardFailed,
  payload: {dashboard},
})

// Thunks

export const getDashboardsAsync = (url: string) => async (
  dispatch: Dispatch<Action>
): Promise<Dashboard[]> => {
  try {
    const dashboards = await getDashboardsAJAX(url)
    dispatch(loadDashboards(dashboards))
    return dashboards
  } catch (error) {
    console.error(error)
    throw error
  }
}

export const importDashboardAsync = (
  url: string,
  dashboard: Dashboard
) => async (dispatch: Dispatch<Action>): Promise<void> => {
  try {
    await createDashboardAJAX(url, dashboard)
    const dashboards = await getDashboardsAJAX(url)

    dispatch(loadDashboards(dashboards))
    dispatch(notify(copy.dashboardImported(name)))
  } catch (error) {
    dispatch(
      notify(copy.dashboardImportFailed('', 'Could not upload dashboard'))
    )
    console.error(error)
  }
}

export const deleteDashboardAsync = (dashboard: Dashboard) => async (
  dispatch: Dispatch<Action>
): Promise<void> => {
  dispatch(deleteDashboard(dashboard.id))
  dispatch(deleteTimeRange(dashboard.id))

  try {
    await deleteDashboardAJAX(dashboard.links.self)
    dispatch(notify(copy.dashboardDeleted(dashboard.name)))
  } catch (error) {
    dispatch(
      notify(copy.dashboardDeleteFailed(dashboard.name, error.data.message))
    )

    dispatch(deleteDashboardFailed(dashboard))
  }
}

export const getDashboardAsync = (dashboardID: string) => async (
  dispatch
): Promise<void> => {
  try {
    const dashboard = await getDashboardAJAX(dashboardID)
    dispatch(loadDashboard(dashboard))
  } catch {
    dispatch(replace(`/dashboards`))
    dispatch(notify(copy.dashboardNotFound(dashboardID)))

    return
  }

  // TODO: Notify if any of the supplied query params were invalid
  dispatch(updateTimeRangeFromQueryParams(dashboardID))
}
