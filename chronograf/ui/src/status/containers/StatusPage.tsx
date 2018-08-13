// Libraries
import React, {Component} from 'react'

// Components
import FancyScrollbar from 'src/shared/components/FancyScrollbar'
import Grid from 'src/shared/components/Grid'
import PageHeader from 'src/reusable_ui/components/page_layout/PageHeader'

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
            <Grid cells={this.cells} />
          </div>
        </FancyScrollbar>
      </div>
    )
  }

  private get cells(): Cell[] {
    return [
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
  }
}

export default StatusPage
