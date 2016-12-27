import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { filter, includes } from 'lodash';
import moment from 'moment';
import { push } from 'react-router-redux';

import Button from 'components/buttons/Button';
import entityGetter from 'redux/utilities/entityGetter';
import InputField from 'components/forms/fields/InputField';
import NumberPill from 'components/NumberPill';
import packActions from 'redux/nodes/entities/packs/actions';
import packInterface from 'interfaces/pack';
import paths from 'router/paths';

const baseClass = 'all-packs-page';

class AllPacksPage extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
    packs: PropTypes.arrayOf(packInterface),
  }

  constructor (props) {
    super(props);

    this.state = { packFilter: '' };
  }

  componentWillMount() {
    const { dispatch, packs } = this.props;

    if (!packs.length) {
      dispatch(packActions.loadAll());
    }

    return false;
  }

  onFilterPacks = (packFilter) => {
    this.setState({ packFilter });

    return false;
  }

  getPacks = () => {
    const { packFilter } = this.state;
    const { packs } = this.props;

    if (!packFilter) {
      return packs;
    }

    const lowerPackFilter = packFilter.toLowerCase();

    return filter(packs, (pack) => {
      if (!pack.name) {
        return false;
      }

      const lowerPackName = pack.name.toLowerCase();

      return includes(lowerPackName, lowerPackFilter);
    });
  }

  goToNewPackPage = () => {
    const { dispatch } = this.props;
    const { NEW_PACK } = paths;

    dispatch(push(NEW_PACK));

    return false;
  }

  renderPack = (pack) => {
    const updatedTime = moment(pack.updated_at);

    return (
      <tr key={`pack-${pack.id}-table`}>
        <td>{pack.name}</td>
        <td>{pack.query_count}</td>
        <td>Enabled?</td>
        <td>Jason Meller?</td>
        <td>{pack.hosts_count}</td>
        <td>{updatedTime.fromNow()}</td>
      </tr>
    );
  }

  render () {
    const { getPacks, goToNewPackPage, onFilterPacks, renderPack } = this;
    const { packFilter } = this.state;
    const packs = getPacks();
    const packsCount = packs.length;

    return (
      <div className={`${baseClass} body-wrap`}>
        <div className={`${baseClass}__wrapper`}>
          <p className={`${baseClass}__title`}>
            <NumberPill number={packsCount} /> Query Packs
          </p>
          <div className={`${baseClass}__search-create-section`}>
            <InputField
              name="pack-filter"
              onChange={onFilterPacks}
              placeholder="Search Packs"
              value={packFilter}
            />
            <Button variant="brand" onClick={goToNewPackPage}>
              CREATE NEW PACK
            </Button>
            <Button
              text="CREATE NEW PACK"
              variant="brand"
              onClick={goToNewPackPage}
            />
          </div>
          <table className={`${baseClass}__table`}>
            <thead>
              <tr>
                <th>Name</th>
                <th>Queries</th>
                <th>Status</th>
                <th>Author</th>
                <th>Number of Hosts</th>
                <th>Last Updated</th>
              </tr>
            </thead>
            <tbody>
              {packs.map((pack) => {
                return renderPack(pack);
              })}
            </tbody>
          </table>
        </div>
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  const { entities: packs } = entityGetter(state).get('packs');

  return { packs };
};

export default connect(mapStateToProps)(AllPacksPage);
