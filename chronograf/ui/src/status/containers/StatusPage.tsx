// Libraries
import React, {Component} from 'react'

// Components
import FancyScrollbar from 'src/shared/components/FancyScrollbar'
import Grid from 'src/shared/components/Grid'
import PageHeader from 'src/reusable_ui/components/page_layout/PageHeader'

// Constants
import {fixtureStatusPageCells} from 'src/status/fixtures'

// Types
import {Source, Cell} from 'src/types/v2'

import {ErrorHandling} from 'src/shared/decorators/errors'

interface Props {
  source: Source
  params: {
    sourceID: string
  }
}

@ErrorHandling
class StatusPage extends Component<Props> {
  public render() {
    return (
      <div className="page">
        <PageHeader
          titleText="Status"
          fullWidth={true}
          sourceIndicator={true}
        />
        <FancyScrollbar className="page-contents">
          <div className="dashboard container-fluid full-width">
            <Grid cells={fixtureStatusPageCells} />
          </div>
        </FancyScrollbar>
      </div>
    )
  }
}

export default StatusPage
