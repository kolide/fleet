import { Table, Column, Cell } from 'fixed-data-table';
import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';

import Button from 'components/buttons/Button';
import entityGetter from 'redux/utilities/entityGetter';
import packActions from 'redux/nodes/entities/packs/actions';
import packInterface from 'interfaces/pack';
import PackInfoSidePanel from 'components/side_panels/PackInfoSidePanel';
import paths from 'router/paths';

const baseClass = 'all-packs-page';

class AllPacksPage extends Component {

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
        <li> {pack.name} </li>
      );
    });
  }

  _onSearchBarChange(e) {
    console.log(e);
  }

  render () {
    const { renderPacks, _onSearchBarChange } = this;
    const { packs, dispatch } = this.props;

    return (
      <div>

        <div className={`${baseClass}__wrapper`}>
          <p className={`${baseClass}__title`}>
            Query Packs
          </p>

          <div className={`${baseClass}__search`}>
            <input
              onChange={_onSearchBarChange}
              placeholder="SEARCH"
            />
          </div>

          <div>
            multi-action
          </div>

          <div>
            <Button
              text={"CREATE NEW PACK"}
              variant={"brand"}
              onClick={evt => { dispatch(push(paths.NEW_PACK)) }}
            />
          </div>

          <Table
            rowHeight={50}
            rowsCount={packs.length}
            width={1000}
            height={300}
            headerHeight={50}>

            <Column
              header={<Cell>Pack Name</Cell>}
              width={250}
              cell={({rowIndex, ...props}) => (
                <Cell {...props}>
                  {packs[rowIndex]["name"]}
                </Cell>
              )}
            />
            <Column
              header={<Cell>Queries</Cell>}
              width={100}
              cell={({rowIndex, ...props}) => (
                <Cell {...props}>
                  10?
                </Cell>
              )}
            />
            <Column
              header={<Cell>Status</Cell>}
              width={150}
              cell={({rowIndex, ...props}) => (
                <Cell {...props}>
                  Enabled?
                </Cell>
              )}
            />
            <Column
              header={<Cell>Author</Cell>}
              width={300}
              cell={({rowIndex, ...props}) => (
                <Cell {...props}>
                  Jason Meller?
                </Cell>
              )}
            />
            <Column
              header={<Cell># Hosts</Cell>}
              width={100}
              cell={({rowIndex, ...props}) => (
                <Cell {...props}>
                  9001?
                </Cell>
              )}
            />
            <Column
              header={<Cell>Last Updated</Cell>}
              width={100}
              cell={({rowIndex, ...props}) => (
                <Cell {...props}>
                  Yesterday?
                </Cell>
              )}
            />
          </Table>

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
