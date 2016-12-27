import React, { Component, PropTypes } from 'react';
import { includes, pull } from 'lodash';

import Checkbox from 'components/forms/fields/Checkbox';
import packInterface from 'interfaces/pack';
import Row from 'components/packs/PacksList/Row';

const baseClass = 'packs-list';

class PacksList extends Component {
  static propTypes = {
    className: PropTypes.string,
    packs: PropTypes.arrayOf(packInterface),
  };

  static defaultProps = {
    packs: [],
  };

  constructor (props) {
    super(props);

    this.state = { allPacksChecked: false, checkedPackIDs: [] };
  }

  onCheckAllPacks = (shouldCheck) => {
    this.setState({ allPacksChecked: shouldCheck });

    return false;
  }

  onCheckPack = (checked, id) => {
    const { checkedPackIDs } = this.state;
    const newCheckedPackIDs = checked ? checkedPackIDs.concat(id) : pull(checkedPackIDs, id);

    this.setState({ checkedPackIDs: newCheckedPackIDs });

    return false;
  }

  renderPack = (pack) => {
    const { allPacksChecked, checkedPackIDs } = this.state;
    const { onCheckPack } = this;
    const checked = allPacksChecked || includes(checkedPackIDs, pack.id);

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
    const { allPacksChecked } = this.state;
    const { className, packs } = this.props;
    const { onCheckAllPacks, renderPack } = this;

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
