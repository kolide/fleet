import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { noop } from 'lodash';
import classnames from 'classnames';

import { renderFlash } from 'redux/nodes/notifications/actions';
import Icon from 'components/icons/Icon';
import { copyText } from './helpers';

const HOST_TABS = {
  FIRST: 'What Does This Script Do?',
  SECOND: 'Additional Script Options',
};

const baseClass = 'new-host';

export class NewHostPage extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
  };

  static defaultProps = {
    dispatch: noop,
  };

  constructor (props) {
    super(props);

    this.state = {
      osqueryCommandText: `
osqueryd
 --enroll_secret_env=OSQUERY_ENROLL_SECRET
 --tls_server_certs=/etc/osquery/kolide.crt

 --tls_hostname=acme.kolide.co

 --host_identifier=hostname
 --enroll_tls_endpoint=/api/v1/osquery/enroll

 --config_plugin=tls
 --config_tls_endpoint=/api/v1/osquery/config
 --config_tls_refresh=10

 --disable_distributed=false
 --distributed_plugin=tls
 --distributed_interval=10
 --distributed_tls_max_attempts=3
 --distributed_tls_read_endpoint=/api/v1/osquery/distributed/read
 --distributed_tls_write_endpoint=/api/v1/osquery/distributed/write

 --logger_plugin=tls
 --logger_tls_endpoint=/api/v1/osquery/log
 --logger_tls_period=10
      `,
      osqueryCommandTextCopied: false,
      selectedTab: HOST_TABS.FIRST,
    };
  }

  onCopyText = (text, elementId) => {
    return (evt) => {
      evt.preventDefault();

      const { dispatch } = this.props;
      const { osqueryCommandText } = this.state;

      if (copyText(elementId)) {
        dispatch(renderFlash('success', 'Text copied to clipboard'));
      } else {
        dispatch(renderFlash('error', 'Text not copied. Use CMD + C to copy text'));
      }

      if (text === osqueryCommandText) {
        this.setState({
          osqueryCommandTextCopied: true,
        });
      }

      setTimeout(() => {
        this.setState({
          osqueryCommandTextCopied: false,
        });

        return false;
      }, 1500);

      return false;
    };
  }

  render () {
    const { osqueryCommandText, osqueryCommandTextCopied } = this.state;
    const { onCopyText } = this;

    const osqueryCommandIconClasses = classnames(
      `${baseClass}__clipboard-icon`,
      {
        [`${baseClass}__clipboard-icon--copied`]: osqueryCommandTextCopied,
      }
    );

    return (
      <div className={baseClass}>
        <section className={`${baseClass}__section-wrap body-wrap`}>
          <h1 className={`${baseClass}__title`}>Kolide Installation Instructions</h1>

          <div className={`${baseClass}__text`}>
            <p>To use Kolide, you must install the open source osquery tool on the hosts which you wish to monitor. You can find various ways to install osquery on a variety of platforms at <a href="https://osquery.io/downloads">https://osquery.io/downloads</a>.</p>
            <br />
            <p>Once you have installed osquery, you need to do two things:</p>
            <ol className="kolide-ol">
              <li>Set an environment variable with an agent enrollment secret</li>
              <li>Deploy the TLS certificate that osquery will use to communicate with Kolide</li>
            </ol>
            <br />
            <p>The enrollment secret is a value that osquery uses to ensure a level of confidence that the host running osquery is actually a host that you would like to hear from. Morbi id varius velit. Phasellus risus arcu, lacinia non cursus a, tempor quis dolor. Nam tempor quam orci, eget semper augue rutrum quis. Fusce eget volutpat ipsum, et rhoncus sapien. Aenean luctus, nulla sit amet facilisis dictum, lacus lorem rutrum tellus, et ultricies nisl nunc vitae ligula. Duis a arcu efficitur, porta nisi sed, placerat justo. Nulla volutpat mollis purus, vel ultricies odio molestie at. Integer vitae mollis nulla, at pellentesque urna. For more information on configuring and deploying enrollment secrets, see the <a href="https://osquery.readthedocs.io/en/stable/deployment/remote/#simple-shared-secret-enrollment">osquery documentation</a></p>
            <br />
            <p>The TLS certificate that osquery will use to communicate with Kolide is the same certificate that your browser is using to communicate with Kolide right now. Nullam sollicitudin odio vitae ipsum consequat commodo. Etiam sodales tempus erat, ut faucibus lorem aliquam et. Nam ac ligula venenatis, ultrices orci vel, suscipit odio. Nam ullamcorper euismod pellentesque. Aenean dapibus risus nec mollis imperdiet. Aenean ac faucibus elit. Nam fringilla vel eros ac pellentesque. Curabitur vulputate sollicitudin posuere. Mauris aliquet non ante eu commodo. For more information on configuring the TLS settings in osquery, see the <a href="https://osquery.readthedocs.io/en/stable/deployment/remote/#remote-authentication">osquery documentation</a></p>
            <br />
            <p>Assuming that you arere deploying your enrollment secret as the environment variable OSQUERY_ENROLL_SECRET and your osquery server certificate is at /etc/osquery/kolide.crt, you could copy and paste the following command with the following flags:</p>
            <br />
          </div>

          <div className={`${baseClass}__input-wrap`}>
            <input id="osqueryCommand" className={`${baseClass}__input`} value={osqueryCommandText} readOnly />
            {osqueryCommandTextCopied && <span className={`${baseClass}__clipboard-text`}>copied!</span>}
            <a href="#copyosqueryCommand" onClick={onCopyText(osqueryCommandText, '#osqueryCommand')}><Icon name="clipboard" className={osqueryCommandIconClasses} /></a>
          </div>
        </section>
      </div>
    );
  }
}

export default connect()(NewHostPage);
