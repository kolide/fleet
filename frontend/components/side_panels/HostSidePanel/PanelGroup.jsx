import React, { Component, PropTypes } from 'react';
import { isEqual, noop } from 'lodash';
import radium from 'radium';

import PanelGroupItem from './PanelGroupItem';

class PanelGroup extends Component {
  static propTypes = {
    groupItems: PropTypes.array,
    onLabelClick: PropTypes.func,
    selectedLabel: PropTypes.object,
  };

  static defaultProps = {
    onLabelClick: noop,
  };

  renderGroupItem = (item) => {
    const {
      onLabelClick,
      selectedLabel,
    } = this.props;
    const selected = isEqual(selectedLabel, item);

    return (
      <PanelGroupItem
        item={item}
        key={item.title}
        onLabelClick={onLabelClick(item)}
        selected={selected}
      />
    );
  }

  render () {
    const { groupItems } = this.props;
    const { renderGroupItem } = this;

    return (
      <div>
        {groupItems.map(item => {
          return renderGroupItem(item);
        })}
      </div>
    );
  }
}

export default radium(PanelGroup);
