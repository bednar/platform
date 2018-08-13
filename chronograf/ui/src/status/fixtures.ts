import {DEFAULT_AXIS} from 'src/dashboards/constants/cellEditor'
import {CellQuery, Axes, CellType} from 'src/types'
import {Cell} from 'src/types/v2'

const emptyQuery: CellQuery = {
  query: '',
  source: '',
  queryConfig: {
    database: '',
    measurement: '',
    retentionPolicy: '',
    fields: [],
    tags: {},
    groupBy: {},
    areTagsAccepted: false,
    rawText: null,
    range: null,
  },
}

const emptyAxes: Axes = {
  x: DEFAULT_AXIS,
  y: DEFAULT_AXIS,
  y2: DEFAULT_AXIS,
}

export const fixtureStatusPageCells: Cell[] = [
  {
    ref: 'news-feed',
    x: 0,
    y: 0,
    w: 8.5,
    h: 10,
  },
  {
    ref: 'getting-started',
    x: 8.5,
    y: 0,
    w: 3.5,
    h: 10,
  },
]
