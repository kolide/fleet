import Base from './base';
import endpoints from './endpoints';
import { appendTargetTypeToTargets } from '../redux/nodes/entities/targets/helpers';

class Kolide extends Base {
  forgotPassword ({ email }) {
    const { FORGOT_PASSWORD } = endpoints;
    const forgotPasswordEndpoint = this.baseURL + FORGOT_PASSWORD;

    return Base.post(forgotPasswordEndpoint, JSON.stringify({ email }));
  }

  getConfig = () => {
    const { CONFIG } = endpoints;

    return this.authenticatedGet(this.endpoint(CONFIG))
      .then((response) => { return response.org_info; });
  }

  getInvites = () => {
    const { INVITES } = endpoints;

    return this.authenticatedGet(this.endpoint(INVITES))
      .then((response) => { return response.invites; });
  }

  getHosts = () => {
    const { HOSTS } = endpoints;

    return this.authenticatedGet(this.endpoint(HOSTS))
      .then((response) => { return response.hosts; });
  }

  getTargets = (options = {}) => {
    console.log(options);

    const stubbedResponse = {
      targets: {
        hosts: [
          {
            id: 3,
            ip: '192.168.1.10',
            label: 'OS X El Capitan 10.11',
            name: 'osx-10.11',
            platform: 'darwin',
          },
          {
            id: 4,
            ip: '192.168.1.11',
            label: 'Jason Meller\'s Windows Note',
            name: 'jmeller.local',
            platform: 'windows',
          },
        ],
        labels: [
          {
            id: 4,
            label: 'All Macs',
            name: 'macs',
            count: 1234,
          },
        ],
      },
      selected_targets_count: 1234,
    };

    return Promise.resolve(stubbedResponse)
      .then((response) => { return appendTargetTypeToTargets(response); });
  }

  getUsers = () => {
    const { USERS } = endpoints;

    return this.authenticatedGet(this.endpoint(USERS))
      .then((response) => { return response.users; });
  }

  inviteUser = (formData) => {
    const { INVITES } = endpoints;

    return this.authenticatedPost(this.endpoint(INVITES), JSON.stringify(formData))
      .then((response) => { return response.invite; });
  }

  loginUser ({ username, password }) {
    const { LOGIN } = endpoints;
    const loginEndpoint = this.baseURL + LOGIN;

    return Base.post(loginEndpoint, JSON.stringify({ username, password }));
  }

  logout () {
    const { LOGOUT } = endpoints;
    const logoutEndpoint = this.baseURL + LOGOUT;

    return this.authenticatedPost(logoutEndpoint);
  }

  me () {
    const { ME } = endpoints;
    const meEndpoint = this.baseURL + ME;

    return this.authenticatedGet(meEndpoint);
  }

  resetPassword (formData) {
    const { RESET_PASSWORD } = endpoints;
    const resetPasswordEndpoint = this.baseURL + RESET_PASSWORD;

    return Base.post(resetPasswordEndpoint, JSON.stringify(formData));
  }

  revokeInvite = ({ entityID }) => {
    const { INVITES } = endpoints;
    const endpoint = `${this.endpoint(INVITES)}/${entityID}`;

    return this.authenticatedDelete(endpoint);
  }

  updateUser = (user, formData) => {
    const { USERS } = endpoints;
    const updateUserEndpoint = `${this.baseURL}${USERS}/${user.id}`;

    return this.authenticatedPatch(updateUserEndpoint, JSON.stringify(formData))
      .then((response) => { return response.user; });
  }
}

export default new Kolide();
