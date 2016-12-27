import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { filter, includes } from 'lodash';
import { push } from 'react-router-redux';

import Button from 'components/buttons/Button';
import entityGetter from 'redux/utilities/entityGetter';
import InputField from 'components/forms/fields/InputField';
import NumberPill from 'components/NumberPill';
import packActions from 'redux/nodes/entities/packs/actions';
import PackInfoSidePanel from 'components/side_panels/PackInfoSidePanel';
import packInterface from 'interfaces/pack';
import PacksList from 'components/packs/PacksList';
import paths from 'router/paths';

const baseClass = 'all-packs-page';

class AllPacksPage extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
    packs: PropTypes.arrayOf(packInterface),
  }

  constructor (props) {
    super(props);

    this.state = {
      packFilter: '',
      selectedPack: undefined,
    };
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

  renderSidePanel = () => {
    const { selectedPack } = this.state;

    if (!selectedPack) {
      return <PackInfoSidePanel />;
    }

    // TODO: render PackDetailSidePanel
    return false;
  }

  render () {
    const {
      getPacks,
      goToNewPackPage,
      onFilterPacks,
      renderSidePanel,
    } = this;
    const { packFilter } = this.state;
    const packs = getPacks();
    const packsCount = packs.length;

    return (
      <div className={`${baseClass} has-sidebar`}>
        <div className={`${baseClass}__wrapper body-wrap`}>
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
          <PacksList className={`${baseClass}__table`} packs={packs} />
        </div>
        {renderSidePanel()}
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  const { entities: packs } = entityGetter(state).get('packs');

  return { packs };
};

export default connect(mapStateToProps)(AllPacksPage);
