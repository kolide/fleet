import React, { Component } from 'react';

import Button from 'components/buttons/Button';
import Icon from 'components/icons/Icon';
import queryInterface from 'interfaces/query';
import SecondarySidePanelContainer from 'components/side_panels/SecondarySidePanelContainer';

class QueryDetailsSidePanel extends Component {
  static propTypes = {
    query: queryInterface.isRequired,
  };

  renderPacks = () => {
    const { packs } = this.props.query;

    if (!packs || (packs && !packs.length)) {
      return <p>There are no packs associated with this query</p>;
    }

    return (
      <div>
        {packs.map((pack) => {
          return (
            <div key={`query-side-panel-pack-${pack.id}`}>
              <Icon name="packs" />
              <span>{pack.name}</span>
            </div>
          );
        })}
      </div>
    );
  }

  render () {
    const { query } = this.props;
    const { renderPacks } = this;

    return (
      <SecondarySidePanelContainer>
        <h1>{query.name}</h1>
        <Button variant="inverse">Edit/Run Query</Button>
        <h2>Description</h2>
        <p>{query.description}</p>
        <h2>Packs</h2>
        {renderPacks()}
      </SecondarySidePanelContainer>
    );
  }
}

export default QueryDetailsSidePanel;
