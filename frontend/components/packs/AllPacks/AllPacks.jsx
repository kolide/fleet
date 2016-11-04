import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';

import entityGetter from 'redux/utilities/entityGetter';
import packActions from 'redux/nodes/entities/packs/actions';
import packInterface from 'interfaces/pack';

const baseClass = 'all-packs';

class AllPacks extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
    packs: PropTypes.arrayOf(packInterface),
  }

  componentWillMount() {
    const { dispatch, packs } = this.props;
    if (!packs.length) {
      dispatch(packActions.loadAll());
    }

    return false;
  }

  renderPacks = () => {
    const { packs } = this.props;

    return packs.map((pack) => {
      return (
        <li> { pack.name } </li>
      );
    });
  }

  render () {
    const { renderPacks } = this;

    return (
      <div className={`${baseClass}__wrapper`}>
        <p className={`${baseClass}__title`}>
          Query Packs
        </p>
        <ul>
          {renderPacks()}
        </ul>
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  const { entities: packs } = entityGetter(state).get('packs');

  return { packs };
};

export default connect(mapStateToProps)(AllPacks);
