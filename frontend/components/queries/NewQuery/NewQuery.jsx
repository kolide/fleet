import React, { Component, PropTypes } from 'react';
import AceEditor from 'react-ace';
import 'brace/mode/sql';
import 'brace/ext/linking';

import Button from 'components/buttons/Button';
import debounce from 'utilities/debounce';
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
    onInvalidQuerySubmit: PropTypes.func,
    onNewQueryFormSubmit: PropTypes.func,
    onOsqueryTableSelect: PropTypes.func,
    onRemoveMoreInfoTarget: PropTypes.func,
    onTargetSelectInputChange: PropTypes.func,
    onTargetSelectMoreInfo: PropTypes.func,
    onTextEditorInputChange: PropTypes.func,
    selectedTargetsCount: PropTypes.number,
    targets: PropTypes.arrayOf(targetInterface),
    textEditorText: PropTypes.string,
  };

  constructor (props) {
    super(props);

    this.state = {
      isSaveQueryForm: false,
      selectedTargets: [],
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

  onRunQuery = (evt) => {
    evt.preventDefault();

    const { onInvalidQuerySubmit, textEditorText } = this.props;
    const { error } = validateQuery(textEditorText);
    const { onSaveQueryFormSubmit } = this;
    const { selectedTargets } = this.state;

    if (error) {
      return onInvalidQuerySubmit(error);
    }

    return onSaveQueryFormSubmit({ selectedTargets });
  }

  onSaveQueryFormCancel = (evt) => {
    evt.preventDefault();

    this.setState({ isSaveQueryForm: false });

    return false;
  }

  onSaveQueryFormSubmit = debounce((formData) => {
    const { onInvalidQuerySubmit, onNewQueryFormSubmit, textEditorText } = this.props;
    const { error } = validateQuery(textEditorText);

    this.setState({ isSaveQueryForm: false });

    if (error) {
      return onInvalidQuerySubmit(error);
    }

    return onNewQueryFormSubmit({
      ...formData,
      query: textEditorText,
    });
  })

  onTargetSelect = (selectedTargets) => {
    this.setState({ selectedTargets });
    return false;
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
      onTargetSelectInputChange,
      onTargetSelectMoreInfo,
      onTextEditorInputChange,
      selectedTargetsCount,
      targets,
      textEditorText,
    } = this.props;
    const { selectedTargets } = this.state;
    const {
      onLoad,
      onLoadSaveQueryModal,
      onRunQuery,
      onTargetSelect,
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
