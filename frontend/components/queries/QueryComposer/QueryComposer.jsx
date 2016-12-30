import React, { Component, PropTypes } from 'react';
import 'brace/mode/sql';
import 'brace/ext/linking';

import QueryForm from 'components/forms/queries/QueryForm';
import queryInterface from 'interfaces/query';
import targetInterface from 'interfaces/target';
import './mode';
import './theme';

const baseClass = 'query-composer';

class QueryComposer extends Component {
  static propTypes = {
    onFetchTargets: PropTypes.func,
    onFormCancel: PropTypes.func,
    onOsqueryTableSelect: PropTypes.func,
    onRunQuery: PropTypes.func,
    onSave: PropTypes.func,
    onStopQuery: PropTypes.func,
    onTargetSelect: PropTypes.func,
    onUpdate: PropTypes.func,
    query: queryInterface,
    queryIsRunning: PropTypes.bool,
    queryType: PropTypes.string,
    selectedTargets: PropTypes.arrayOf(targetInterface),
    targetsCount: PropTypes.number,
  };

  static defaultProps = {
    queryType: 'query',
    targetsCount: 0,
  };

  renderForm = () => {
    const {
      errors,
      onFetchTargets,
      onFormCancel,
      onOsqueryTableSelect,
      onRunQuery,
      onSave,
      onStopQuery,
      onTargetSelect,
      onUpdate,
      query,
      queryIsRunning,
      queryType,
      selectedTargets,
      targetsCount,
    } = this.props;

    return (
      <QueryForm
        onCancel={onFormCancel}
        onFetchTargets={onFetchTargets}
        onOsqueryTableSelect={onOsqueryTableSelect}
        onRunQuery={onRunQuery}
        handleSubmit={onSave}
        onStopQuery={onStopQuery}
        onTargetSelect={onTargetSelect}
        onUpdate={onUpdate}
        formData={query}
        query={query}
        queryIsRunning={queryIsRunning}
        queryType={queryType}
        selectedTargets={selectedTargets}
        serverErrors={errors}
        targetsCount={targetsCount}
      />
    );
  }

  render () {
    const { queryIsRunning, queryText, queryType } = this.props;
    const { onLoad, renderForm, renderTargetsInput } = this;

    return (
      <div className={`${baseClass}__wrapper body-wrap`}>
        {renderForm()}
      </div>
    );
  }
}

export default QueryComposer;
