// Libraries
import React, {Component} from 'react'
import _ from 'lodash'
import download from 'src/external/download.js'

// Components
import LayoutCellMenu from 'src/shared/components/LayoutCellMenu'
import LayoutCellHeader from 'src/shared/components/LayoutCellHeader'
import View from 'src/shared/components/cells/View'

// Actions
import {notify} from 'src/shared/actions/notifications'

// Utils
import {dataToCSV} from 'src/shared/parsing/dataToCSV'
import {timeSeriesToTableGraph} from 'src/utils/timeSeriesTransformers'

// Constants
import {PREDEFINED_TEMP_VARS} from 'src/shared/constants'
import {csvDownloadFailed} from 'src/shared/copy/notifications'

// Types
import {CellQuery} from 'src/types'
import {Cell} from 'src/types/v2'

import {ErrorHandling} from 'src/shared/decorators/errors'

interface Props {
  cell: Cell
  onDeleteCell: (cell: Cell) => void
  onCloneCell: (cell: Cell) => void
  isEditable: boolean
}

@ErrorHandling
export default class CellComponent extends Component<Props> {
  public render() {
    const {cell, isEditable, onDeleteCell, onCloneCell} = this.props

    return (
      <div className="dash-graph">
        <LayoutCellMenu
          cell={cell}
          queries={this.queries}
          isEditable={isEditable}
          onDelete={onDeleteCell}
          onClone={onCloneCell}
          dataExists={false}
          onEdit={this.handleSummonOverlay}
          onCSVDownload={this.handleCSVDownload}
        />
        <LayoutCellHeader cellName="ima cell" isEditable={isEditable} />
        <div className="dash-graph--container">{this.emptyGraph}</div>
      </div>
    )
  }

  private get queries(): CellQuery[] {
    const {cell} = this.props
    return _.get(cell, ['queries'], [])
  }

  private get emptyGraph(): JSX.Element {
    return (
      <div className="graph-empty">
        <button
          className="no-query--button btn btn-md btn-primary"
          onClick={this.handleSummonOverlay}
        >
          <span className="icon plus" /> Add Data
        </button>
      </div>
    )
  }

  private handleSummonOverlay = (): void => {
    // TODO: add back in once CEO is refactored
  }

  private handleCSVDownload = (): void => {
    // TODO: get data from link
    // const {cellData, cell} = this.props
    // const joinedName = cell.name.split(' ').join('_')
    // const {data} = timeSeriesToTableGraph(cellData)
    // try {
    //   download(dataToCSV(data), `${joinedName}.csv`, 'text/plain')
    // } catch (error) {
    //   notify(csvDownloadFailed())
    //   console.error(error)
    // }
  }
}
