import React, { Component, PropTypes } from 'react';
import AceEditor from 'react-ace';
import classnames from 'classnames';
import { noop } from 'lodash';

import Button from '../../buttons/Button';
import { headerClassName } from './helpers';
import hostHelpers from '../../hosts/HostDetails/helpers';
import Modal from '../Modal';
import ShadowBox from '../../ShadowBox';
import targetInterface from '../../../interfaces/target';

const baseClass = 'target-info-modal';

class TargetInfoModal extends Component {
  static propTypes = {
    className: PropTypes.string,
    onAdd: PropTypes.func,
    onExit: PropTypes.func,
    target: targetInterface.isRequired,
  };

  renderHeader = () => {
    const { target } = this.props;
    const { label } = target;
    const className = headerClassName(target);

    return (
      <span className={`${baseClass}__header`}>
        <i className={className} />
        <span>{label}</span>
      </span>
    )
  }

  renderHostModal = () => {
    const { className, onAdd, onExit, target } = this.props;
    const hostBaseClass = `${baseClass}__host`;
    const {
      ip,
      label,
      mac,
      memory,
      platform,
      os_version,
      osquery_version,
      status,
    } = target;
    const isOnline = status === 'online';
    const isOffline = status === 'offline';
    const { renderHeader } = this;
    const statusClassName = classnames(`${hostBaseClass}__status`, {
      'is-online': isOnline,
      'is-offline': isOffline,
    });

    return (
      <Modal
        className={className}
        onExit={onExit}
        title={renderHeader()}
      >
        <p className={statusClassName}>{status}</p>
        <ShadowBox>
          <table className={`${hostBaseClass}__table`}>
            <tbody>
              <tr>
                <th>IP Address</th>
                <td>{ip}</td>
              </tr>
              <tr>
                <th>MAC Address</th>
                <td>{mac}</td>
              </tr>
              <tr>
                <th>Platform</th>
                <td>
                  <i className={hostHelpers.platformIconClass(platform)} />
                  <span className={`${hostBaseClass}__platform-text`}>{platform}</span>
                </td>
              </tr>
              <tr>
                <th>Operating System</th>
                <td>{os_version}</td>
              </tr>
              <tr>
                <th>Osquery Version</th>
                <td>{osquery_version}</td>
              </tr>
              <tr>
                <th>Memory</th>
                <td>{hostHelpers.humanMemory(memory)}</td>
              </tr>
            </tbody>
          </table>
        </ShadowBox>
        <div className={`${hostBaseClass}__labels-wrapper`}>
          <div className={`${hostBaseClass}__labels-wrapper--header`}>
            <i className="kolidecon-label" />
            <span>Labels</span>
          </div>
        </div>
        <div className={`${baseClass}__btn-wrapper`}>
          <Button text="CANCEL" variant="inverse" onClick={onExit} />
          <Button text="ADD TO TARGETS" onClick={onAdd} />
        </div>
      </Modal>
    );
  }

  renderLabelModal = () => {
    const { className, onAdd, onExit, target } = this.props;
    const {
      description,
      hosts,
      label,
      query,
    } = target;
    const { renderHeader } = this;

    return (
      <Modal
        className={className}
        onExit={onExit}
        title={renderHeader()}
      >
        <p>{description}</p>
        <div className={`${baseClass}__text-editor-wrapper`}>
          <AceEditor
            editorProps={{ $blockScrolling: Infinity }}
            mode="kolide"
            minLines={4}
            maxLines={4}
            name="modal-label-query"
            onChange={noop}
            showGutter
            showPrintMargin={false}
            theme="kolide"
            value={query}
            width="100%"
          />
        </div>
        <Button text="ADD" onClick={onAdd} />
      </Modal>
    );
  }

  render () {
    const { renderHostModal, renderLabelModal } = this;
    const { target_type: targetType } = this.props.target;

    if (targetType === 'hosts') {
      return renderHostModal();
    }

    return renderLabelModal();
  }
}

export default TargetInfoModal;
