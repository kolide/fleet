import React, { Component, PropTypes } from 'react';
import AceEditor from 'react-ace';
import classnames from 'classnames';

import hostHelpers from 'components/hosts/HostDetails/helpers';
import targetInterface from 'interfaces/target';

const baseClass = 'target-details';

class TargetDetails extends Component {
  static propTypes = {
    target: targetInterface,
    className: PropTypes.string,
  };

  renderHost = () => {
    const { className, target } = this.props;
    const {
      display_text,
      ip,
      mac,
      memory,
      osqueryVersion,
      osVersion,
      platform,
      status,
    } = target;
    const hostBaseClass = 'host-target';
    const isOnline = status === 'online';
    const isOffline = status === 'offline';
    const statusClassName = classnames(
      `${hostBaseClass}__status`,
      { [`${hostBaseClass}__status--is-online`]: isOnline },
      { [`${hostBaseClass}__status--is-offline`]: isOffline },
    );

    return (
      <div className={`${hostBaseClass} ${className}`}>
        <p className={`${hostBaseClass}__display-text`}>
          <i className={`${hostBaseClass}__icon kolidecon-fw kolidecon-single-host`} />
          <span>{display_text}</span>
        </p>
        <p className={statusClassName}>
          {isOnline && <i className={`${hostBaseClass}__icon ${hostBaseClass}__icon--online kolidecon-fw kolidecon-success-check`} />}
          {isOffline && <i className={`${hostBaseClass}__icon ${hostBaseClass}__icon--offline kolidecon-fw kolidecon-offline`} />}
          <span>{status}</span>
        </p>
          <table className={`${baseClass}__table`}>
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
              <td>{osVersion}</td>
            </tr>
            <tr>
              <th>Osquery Version</th>
              <td>{osqueryVersion}</td>
            </tr>
            <tr>
              <th>Memory</th>
              <td>{hostHelpers.humanMemory(memory)}</td>
            </tr>
          </tbody>
        </table>
        <div className={`${hostBaseClass}__labels-wrapper`}>
          <p className={`${hostBaseClass}__labels-header`}>
            <i className={`${hostBaseClass}__icon kolidecon-fw kolidecon-label`} />
            <span>Labels</span>
          </p>
        </div>
      </div>
    );
  }

  renderLabel = () => {
    const { className, target } = this.props;
    const {
      display_text,
      hosts,
      query,
    } = target;
    const labelBaseClass = 'label-target';

    return (
      <div className={`${labelBaseClass} ${className}`}>
      <p className={`${labelBaseClass}__display-text`}><i className={`${labelBaseClass}__icon kolidecon-fw kolidecon-label`} /> {display_text}</p>
        <div className={`${labelBaseClass}__text-editor-wrapper`}>
          <AceEditor
            editorProps={{ $blockScrolling: Infinity }}
            mode="kolide"
            minLines={4}
            maxLines={4}
            name="label-query"
            readOnly
            setOptions={{ wrap: true }}
            showGutter={false}
            showPrintMargin={false}
            theme="kolide"
            value={query}
            width="100%"
          />
        </div>
        <div className={`${labelBaseClass}__search-section`}>
          <div className={`${labelBaseClass}__num-hosts-section`}>
            <span className="num-hosts">{hosts.length} HOSTS</span>
          </div>
        </div>
        <table className={`${baseClass}__table`}>
          <thead>
            <tr>
              <th>Hostname</th>
              <th>Status</th>
              <th>Platform</th>
              <th>Location</th>
              <th>MAC</th>
            </tr>
          </thead>
          <tbody>
            {hosts.map((host) => {
              return (
                <tr className={`${baseClass}__label-row`} key={`host-${host.id}`}>
                  <td>{host.hostname}</td>
                  <td>{host.status}</td>
                  <td><i className={hostHelpers.platformIconClass(host.platform)} /></td>
                  <td>{host.ip}</td>
                  <td>{host.mac}</td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    );
  }

  render () {
    const { target } = this.props;

    if (!target) {
      return false;
    }

    const { target_type: targetType } = target;
    const { renderHost, renderLabel } = this;

    if (targetType === 'labels') {
      return renderLabel();
    }

    return renderHost();
  }
}

export default TargetDetails;
