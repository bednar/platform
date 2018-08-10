import AJAX from 'src/utils/ajax'
import {Dashboard} from 'src/types/v2'

// TODO(desa): what to do about getting dashboards from another v2 source
export const getDashboards = async (url: string): Promise<Dashboard[]> => {
  try {
    const {data} = await AJAX({
      url,
    })

    return data.dashboards
  } catch (error) {
    throw error
  }
}

export const getDashboard = async (id: string): Promise<Dashboard> => {
  try {
    const {data} = await AJAX({
      url: `/v2/dashboards/${id}`,
    })

    return data
  } catch (error) {
    throw error
  }
}

export const createDashboard = async (
  url: string,
  dashboard: Partial<Dashboard>
): Promise<Dashboard> => {
  try {
    const {data} = await AJAX({
      method: 'POST',
      url,
      data: dashboard,
    })

    return data
  } catch (error) {
    console.error(error)
    throw error
  }
}

export const deleteDashboard = async (url: string): Promise<void> => {
  try {
    return await AJAX({
      method: 'DELETE',
      url,
    })
  } catch (error) {
    console.error(error)
    throw error
  }
}
