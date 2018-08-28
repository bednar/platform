import {
  Source,
  SourceAuthenticationMethod,
  Template,
  Dashboard,
  Cell,
  CellType,
  SourceLinks,
  TemplateType,
  TemplateValueType,
} from 'src/types'

export const role = {
  name: '',
  organization: '',
}

export const currentOrganization = {
  name: '',
  defaultRole: '',
  id: '',
  links: {
    self: '',
  },
}

export const me = {
  currentOrganization,
  role,
}

export const sourceLinks: SourceLinks = {
  query: '/chronograf/v1/sources/16/query',
  services: '/chronograf/v1/sources/16/services',
  self: '/chronograf/v1/sources/16',
  kapacitors: '/chronograf/v1/sources/16/kapacitors',
  proxy: '/chronograf/v1/sources/16/proxy',
  queries: '/chronograf/v1/sources/16/queries',
  write: '/chronograf/v1/sources/16/write',
  permissions: '/chronograf/v1/sources/16/permissions',
  users: '/chronograf/v1/sources/16/users',
  databases: '/chronograf/v1/sources/16/dbs',
  annotations: '/chronograf/v1/sources/16/annotations',
  health: '/chronograf/v1/sources/16/health',
}

export const source: Source = {
  id: '16',
  name: 'ssl',
  type: 'influx',
  username: 'admin',
  url: 'https://localhost:9086',
  insecureSkipVerify: true,
  default: false,
  telegraf: 'telegraf',
  links: sourceLinks,
  authentication: SourceAuthenticationMethod.Basic,
}

export const timeRange = {
  lower: 'now() - 15m',
  upper: null,
}

export const query = {
  id: '0',
  database: 'db1',
  measurement: 'm1',
  retentionPolicy: 'r1',
  fill: 'null',
  fields: [
    {
      value: 'f1',
      type: 'field',
      alias: 'foo',
      args: [],
    },
  ],
  tags: {
    tk1: ['tv1', 'tv2'],
  },
  groupBy: {
    time: null,
    tags: [],
  },
  areTagsAccepted: true,
  rawText: null,
  status: null,
  shifts: [],
}

export const kapacitor = {
  url: '/foo/bar/baz',
  name: 'kapa',
  username: 'influx',
  password: '',
  active: false,
  insecureSkipVerify: false,
  links: {
    self: '/kapa/1',
    proxy: '/proxy/kapacitor/1',
  },
}

export const service = {
  id: '1',
  sourceID: '1',
  url: 'localhost:8082',
  type: 'flux',
  name: 'Flux',
  username: '',
  password: '',
  active: false,
  insecureSkipVerify: false,
  links: {
    source: '/chronograf/v1/sources/1',
    proxy: '/chronograf/v1/sources/1/services/2/proxy',
    self: '/chronograf/v1/sources/1/services/2',
  },
  metadata: {},
}

export const layout = {
  id: '6dfb4d49-20dc-4157-9018-2b1b1cb75c2d',
  app: 'apache',
  measurement: 'apache',
  autoflow: false,
  cells: [
    {
      x: 0,
      y: 0,
      w: 4,
      h: 4,
      i: '0246e457-916b-43e3-be99-211c4cbc03e8',
      name: 'Apache Bytes/Second',
      queries: [
        {
          query:
            'SELECT non_negative_derivative(max("BytesPerSec")) AS "bytes_per_sec" FROM ":db:".":rp:"."apache"',
          groupbys: ['"server"'],
          label: 'bytes/s',
        },
      ],
      axes: {
        x: {
          bounds: [],
          label: '',
          prefix: '',
          suffix: '',
          base: '',
          scale: '',
        },
        y: {
          bounds: [],
          label: '',
          prefix: '',
          suffix: '',
          base: '',
          scale: '',
        },
        y2: {
          bounds: [],
          label: '',
          prefix: '',
          suffix: '',
          base: '',
          scale: '',
        },
      },
      type: '',
      colors: [],
    },
    {
      x: 4,
      y: 0,
      w: 4,
      h: 4,
      i: '37f2e4bb-9fa5-4891-a424-9df5ce7458bb',
      name: 'Apache - Requests/Second',
      queries: [
        {
          query:
            'SELECT non_negative_derivative(max("ReqPerSec")) AS "req_per_sec" FROM ":db:".":rp:"."apache"',
          groupbys: ['"server"'],
          label: 'requests/s',
        },
      ],
      axes: {
        x: {
          bounds: [],
          label: '',
          prefix: '',
          suffix: '',
          base: '',
          scale: '',
        },
        y: {
          bounds: [],
          label: '',
          prefix: '',
          suffix: '',
          base: '',
          scale: '',
        },
        y2: {
          bounds: [],
          label: '',
          prefix: '',
          suffix: '',
          base: '',
          scale: '',
        },
      },
      type: '',
      colors: [],
    },
    {
      x: 8,
      y: 0,
      w: 4,
      h: 4,
      i: 'ea9174b3-2b56-4e80-a37d-064507c6775a',
      name: 'Apache - Total Accesses',
      queries: [
        {
          query:
            'SELECT non_negative_derivative(max("TotalAccesses")) AS "tot_access" FROM ":db:".":rp:"."apache"',
          groupbys: ['"server"'],
          label: 'accesses/s',
        },
      ],
      axes: {
        x: {
          bounds: [],
          label: '',
          prefix: '',
          suffix: '',
          base: '',
          scale: '',
        },
        y: {
          bounds: [],
          label: '',
          prefix: '',
          suffix: '',
          base: '',
          scale: '',
        },
        y2: {
          bounds: [],
          label: '',
          prefix: '',
          suffix: '',
          base: '',
          scale: '',
        },
      },
      type: '',
      colors: [],
    },
  ],
  link: {
    href: '/chronograf/v1/layouts/6dfb4d49-20dc-4157-9018-2b1b1cb75c2d',
    rel: 'self',
  },
}

export const hosts = {
  'MacBook-Pro.local': {
    name: 'MacBook-Pro.local',
    deltaUptime: -1,
    cpu: 0,
    load: 0,
  },
}

// Dashboards
export const template: Template = {
  id: '1',
  type: TemplateType.TagKeys,
  label: 'test query',
  tempVar: ':region:',
  query: {
    db: 'db1',
    rp: 'rp1',
    tagKey: 'tk1',
    fieldKey: 'fk1',
    measurement: 'm1',
    influxql: 'SHOW TAGS WHERE CHRONOGIRAFFE = "friend"',
  },
  values: [
    {
      value: 'us-west',
      type: TemplateValueType.TagKey,
      selected: false,
      localSelected: false,
    },
    {
      value: 'us-east',
      type: TemplateValueType.TagKey,
      selected: true,
      localSelected: true,
    },
    {
      value: 'us-mount',
      type: TemplateValueType.TagKey,
      selected: false,
      localSelected: false,
    },
  ],
}

export const dashboard: Dashboard = {
  id: 1,
  cells: [],
  name: 'd1',
  templates: [],
  organization: 'thebestorg',
}

export const cell: Cell = {
  x: 0,
  y: 0,
  w: 4,
  h: 4,
  i: '0246e457-916b-43e3-be99-211c4cbc03e8',
  name: 'Apache Bytes/Second',
  queries: [],
  axes: {
    x: {
      bounds: ['', ''],
      label: '',
      prefix: '',
      suffix: '',
      base: '',
      scale: '',
    },
    y: {
      bounds: ['', ''],
      label: '',
      prefix: '',
      suffix: '',
      base: '',
      scale: '',
    },
  },
  type: CellType.Line,
  colors: [],
  tableOptions: {
    verticalTimeAxis: true,
    sortBy: {
      internalName: '',
      displayName: '',
      visible: true,
    },
    fixFirstColumn: true,
  },
  fieldOptions: [],
  timeFormat: '',
  decimalPlaces: {
    isEnforced: false,
    digits: 1,
  },
  links: {
    self:
      '/chronograf/v1/dashboards/10/cells/8b3b7897-49b1-422c-9443-e9b778bcbf12',
  },
  legend: {},
  inView: true,
}
