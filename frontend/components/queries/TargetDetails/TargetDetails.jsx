import React, { Component } from 'react';
import AceEditor from 'react-ace';
import classnames from 'classnames';

import hostHelpers from 'components/hosts/HostDetails/helpers';
import ShadowBox from 'components/ShadowBox';
import ShadowBoxInput from 'components/forms/fields/ShadowBoxInput';
import targetInterface from 'interfaces/target';

const baseClass = 'target-details';

class TargetDetails extends Component {
  static propTypes = {
    target: targetInterface.isRequired,
  };

  render () {
    const { target } = this.props;
    const { target_type: targetType } = target;
    const labelBaseClass = 'label-target';

    if (targetType === 'labels') {
      return (
        <div>
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
              value={target.query}
              width="100%"
            />
          </div>
          <div className={`${labelBaseClass}__search-section`}>
            <ShadowBoxInput
              iconClass="kolidecon-search"
              name="search-hosts"
              placeholder="SEARCH HOSTS"
            />
            <div className={`${labelBaseClass}__num-hosts-section`}>
              <span className="num-hosts">{target.hosts.length} HOSTS</span>
            </div>
          </div>
          <ShadowBox>
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
                {target.hosts.map((host) => {
                  return (
                    <tr className="__label-row" key={`host-${host.id}`}>
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
          </ShadowBox>
        </div>
      );
    }

    const hostBaseClass = 'host-target';
    const isOnline = target.status === 'online';
    const isOffline = target.status === 'offline';
    const statusClassName = classnames(
      `${hostBaseClass}__status`,
      { [`${hostBaseClass}__status--is-online`]: isOnline },
      { [`${hostBaseClass}__status--is-offline`]: isOffline },
    );

    return (
      <div>
        <p className={statusClassName}>{target.status}</p>
        <ShadowBox>
          <table className={`${baseClass}__table`}>
            <tbody>
              <tr>
                <th>IP Address</th>
                <td>{target.ip}</td>
              </tr>
              <tr>
                <th>MAC Address</th>
                <td>{target.mac}</td>
              </tr>
              <tr>
                <th>Platform</th>
                <td>
                  <i className={hostHelpers.platformIconClass(target.platform)} />
                  <span className={`${hostBaseClass}__platform-text`}>{target.platform}</span>
                </td>
              </tr>
              <tr>
                <th>Operating System</th>
                <td>{target.osVersion}</td>
              </tr>
              <tr>
                <th>Osquery Version</th>
                <td>{target.osqueryVersion}</td>
              </tr>
              <tr>
                <th>Memory</th>
                <td>{hostHelpers.humanMemory(target.memory)}</td>
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
      </div>
    );
  }
}

export default TargetDetails;
