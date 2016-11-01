import React, { Component, PropTypes } from 'react';
import AceEditor from 'react-ace';
import 'brace/mode/sql';
import 'brace/ext/linking';

import Button from 'components/buttons/Button';
import SaveQueryFormModal from 'components/modals/SaveQueryFormModal';
import SelectTargetsInput from 'components/queries/SelectTargetsInput';
import SelectTargetsMenu from 'components/queries/SelectTargetsMenu';
import targetInterface from 'interfaces/target';
import { validateQuery } from 'components/queries/NewQuery/helpers';
import './mode';
import './theme';

const baseClass = 'new-query';

class NewQuery extends Component {
  static propTypes = {
    isLoadingTargets: PropTypes.bool,
    moreInfoTarget: targetInterface,
    onOsqueryTableSelect: PropTypes.func,
    onRemoveMoreInfoTarget: PropTypes.func,
    onRunQuery: PropTypes.func,
    onSaveQueryFormSubmit: PropTypes.func,
    onTargetSelect: PropTypes.func,
    onTargetSelectInputChange: PropTypes.func,
    onTargetSelectMoreInfo: PropTypes.func,
    onTextEditorInputChange: PropTypes.func,
    selectedTargets: PropTypes.arrayOf(targetInterface),
    selectedTargetsCount: PropTypes.number,
    targets: PropTypes.arrayOf(targetInterface),
    textEditorText: PropTypes.string,
  };

  constructor (props) {
    super(props);

    this.state = {
      isSaveQueryForm: false,
    };
  }

  onLoad = (editor) => {
    editor.setOptions({
      enableLinking: true,
    });

    editor.on('linkClick', (data) => {
      const { type, value } = data.token;
      const { onOsqueryTableSelect } = this.props;

      if (type === 'osquery-token') {
        return onOsqueryTableSelect(value);
      }

      return false;
    });
  }

  onLoadSaveQueryModal = () => {
    this.setState({ isSaveQueryForm: true });

    return false;
  }

  onSaveQueryFormCancel = (evt) => {
    evt.preventDefault();

    this.setState({ isSaveQueryForm: false });

    return false;
  }

  onSaveQueryFormSubmit = (formData) => {
    this.setState({ isSaveQueryForm: false });

    return this.props.onSaveQueryFormSubmit(formData);
  }

  renderSaveQueryFormModal = () => {
    const { isSaveQueryForm } = this.state;
    const {
      onSaveQueryFormSubmit,
      onSaveQueryFormCancel,
    } = this;

    if (!isSaveQueryForm) {
      return false;
    }

    return (
      <SaveQueryFormModal
        onCancel={onSaveQueryFormCancel}
        onSubmit={onSaveQueryFormSubmit}
      />
    );
  }

  render () {
    const {
      isLoadingTargets,
      moreInfoTarget,
      onRemoveMoreInfoTarget,
      onRunQuery,
      onTargetSelect,
      onTargetSelectInputChange,
      onTargetSelectMoreInfo,
      onTextEditorInputChange,
      selectedTargets,
      selectedTargetsCount,
      targets,
      textEditorText,
    } = this.props;
    const {
      onLoad,
      onLoadSaveQueryModal,
      renderSaveQueryFormModal,
    } = this;
    const menuRenderer = SelectTargetsMenu(onTargetSelectMoreInfo, onRemoveMoreInfoTarget, moreInfoTarget);

    return (
      <div className={`${baseClass}__wrapper`}>
        <div className={`${baseClass}__text-editor-wrapper`}>
          <AceEditor
            enableBasicAutocompletion
            enableLiveAutocompletion
            editorProps={{ $blockScrolling: Infinity }}
            mode="kolide"
            minLines={4}
            maxLines={4}
            name="query-editor"
            onLoad={onLoad}
            onChange={onTextEditorInputChange}
            setOptions={{ enableLinking: true }}
            showGutter
            showPrintMargin={false}
            theme="kolide"
            value={textEditorText}
            width="100%"
          />
        </div>
        <div>
          <p>
            <span className={`${baseClass}__select-targets`}>Select Targets</span>
            <span className={`${baseClass}__targets-count`}> {selectedTargetsCount} unique hosts</span>
          </p>
          <SelectTargetsInput
            isLoading={isLoadingTargets}
            menuRenderer={menuRenderer}
            onTargetSelect={onTargetSelect}
            onTargetSelectInputChange={onTargetSelectInputChange}
            selectedTargets={selectedTargets}
            targets={targets}
          />
        </div>
        <div className={`${baseClass}__btn-wrapper`}>
          <Button
            className={`${baseClass}__save-query-btn`}
            onClick={onLoadSaveQueryModal}
            text="Save Query"
            variant="inverse"
          />
          <Button
            className={`${baseClass}__run-query-btn`}
            disabled={!selectedTargets.length}
            onClick={onRunQuery}
            text="Run Query"
          />
        </div>
        {renderSaveQueryFormModal()}
      </div>
    );
  }
}

export default NewQuery;
