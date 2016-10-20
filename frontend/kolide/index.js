import Base from './base';
import endpoints from './endpoints';

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

  getLabels = () => {
    return Promise.resolve({
      labels: [
        { id: 1, label: 'All Hosts', name: 'all', type: 'all', count: 22 },
        { id: 4, label: 'OFFLINE', name: 'offline', type: 'status', count: 2 },
        { id: 5, label: 'ONLINE', name: 'online', type: 'status', count: 20 },
        { id: 6, label: 'MAC OS', name: 'macs', type: 'platform', count: 1 },
        { id: 7, label: 'CENTOS', name: 'centos', type: 'platform', count: 10 },
        { id: 8, label: 'UBUNTU', name: 'ubuntu', type: 'platform', count: 10 },
        { id: 10, label: 'WINDOWS', name: 'windows', type: 'platform', count: 1 },
      ],
    })
      .then(response => { return response.labels; });
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
