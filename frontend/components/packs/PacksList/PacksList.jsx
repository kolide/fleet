import React, { Component, PropTypes } from 'react';
import { includes } from 'lodash';

import Checkbox from 'components/forms/fields/Checkbox';
import packInterface from 'interfaces/pack';
import Row from 'components/packs/PacksList/Row';

const baseClass = 'packs-list';

class PacksList extends Component {
  static propTypes = {
    allPacksChecked: PropTypes.bool,
    checkedPackIDs: PropTypes.arrayOf(PropTypes.number),
    className: PropTypes.string,
    onCheckAllPacks: PropTypes.func.isRequired,
    onCheckPack: PropTypes.func.isRequired,
    packs: PropTypes.arrayOf(packInterface),
  };

  static defaultProps = {
    checkedPackIDs: [],
    packs: [],
  };

  renderPack = (pack) => {
    const { checkedPackIDs, onCheckPack } = this.props;
    const checked = includes(checkedPackIDs, pack.id);

    return (
      <Row
        checked={checked}
        key={`pack-row-${pack.id}`}
        onCheck={onCheckPack}
        pack={pack}
      />
    );
  }

  render () {
    const { allPacksChecked, className, onCheckAllPacks, packs } = this.props;
    const { renderPack } = this;

    return (
      <table className={`${baseClass} ${className}`}>
        <thead>
          <tr>
            <th>
              <Checkbox
                name="select-all-packs"
                onChange={onCheckAllPacks}
                value={allPacksChecked}
              />
            </th>
            <th>Pack Name</th>
            <th>Queries</th>
            <th>Status</th>
            <th>Hosts</th>
            <th>Last Modified</th>
          </tr>
        </thead>
        <tbody>
          {packs.map(pack => renderPack(pack))}
        </tbody>
      </table>
    );
  }
}

export default PacksList;
