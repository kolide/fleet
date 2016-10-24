import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { find } from 'lodash';

import entityGetter from '../../../redux/utilities/entityGetter';
import NewQuery from '../../../components/queries/NewQuery';
import { osqueryTables } from '../../../utilities/osquery_tables';
import QuerySidePanel from '../../../components/side_panels/QuerySidePanel';
import { showRightSidePanel, removeRightSidePanel } from '../../../redux/nodes/app/actions';
import { renderFlash } from '../../../redux/nodes/notifications/actions';
import targetActions from '../../../redux/nodes/entities/targets/actions';
import targetInterface from '../../../interfaces/target';

class NewQueryPage extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
    isLoadingTargets: PropTypes.bool,
    targets: PropTypes.arrayOf(targetInterface),
  };

  componentWillMount () {
    const { dispatch } = this.props;
    const selectedOsqueryTable = find(osqueryTables, { name: 'users' });

    this.state = {
      selectedOsqueryTable,
      textEditorText: 'SELECT * FROM users u JOIN groups g WHERE u.gid = g.gid',
    };

    dispatch(showRightSidePanel);
    dispatch(targetActions.loadAll());

    return false;
  }

  componentWillUnmount () {
    const { dispatch } = this.props;

    dispatch(removeRightSidePanel);

    return false;
  }

  onNewQueryFormSubmit = (formData) => {
    console.log('New Query Form submitted', formData);
  }

  onInvalidQuerySubmit = (errorMessage) => {
    const { dispatch } = this.props;

    dispatch(renderFlash('error', errorMessage));

    return false;
  }

  onOsqueryTableSelect = (tableName) => {
    const selectedOsqueryTable = find(osqueryTables, { name: tableName.toLowerCase() });
    this.setState({ selectedOsqueryTable });

    return false;
  }

  onTextEditorInputChange = (textEditorText) => {
    this.setState({ textEditorText });

    return false;
  }

  render () {
    const {
      onNewQueryFormSubmit,
      onInvalidQuerySubmit,
      onOsqueryTableSelect,
      onTextEditorInputChange,
    } = this;
    const { selectedOsqueryTable, textEditorText } = this.state;
    const { isLoadingTargets, targets } = this.props;

    return (
      <div>
        <NewQuery
          isLoadingTargets={isLoadingTargets}
          onNewQueryFormSubmit={onNewQueryFormSubmit}
          onInvalidQuerySubmit={onInvalidQuerySubmit}
          onOsqueryTableSelect={onOsqueryTableSelect}
          onTextEditorInputChange={onTextEditorInputChange}
          selectedOsqueryTable={selectedOsqueryTable}
          targets={targets}
          textEditorText={textEditorText}
        />
        <QuerySidePanel
          onOsqueryTableSelect={onOsqueryTableSelect}
          onTextEditorInputChange={onTextEditorInputChange}
          selectedOsqueryTable={selectedOsqueryTable}
        />
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  const { entities: targets } = entityGetter(state).get('targets');
  const isLoadingTargets = state.entities.targets.loading;

  return { isLoadingTargets, targets };
};

export default connect(mapStateToProps)(NewQueryPage);
