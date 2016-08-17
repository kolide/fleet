import { Store, toImmutable } from 'nuclear-js'

 var UserStore = new Store({

  getInitialState() {
    return toImmutable({
      logged_in: false,
      is_authenticating: false,
    })
  },

  initialize() {
    this.on("RECEIVE_USER_INFO", receiveUserInfo)
    this.on("IS_AUTHENTICATING", setAuthenticatingState)
  }

})

function receiveUserInfo(state, user) {
  return state.merge(toImmutable({
      id: user.id,
      username: user.username,
      email: user.email,
      name: user.name,
      admin: user.admin,
      needs_password_reset: user.needs_password_reset,
      logged_in: true,
  }))
}

function setAuthenticatingState(state, isAuthenticating) {
  return state.merge(toImmutable({
    is_authenticating: isAuthenticating,
  }))
}

var UserGetters = {
  id: ['user', 'id'],
  username: ['user', 'username'],
  email: ['user', 'email'],
  name: ['user', 'name'],
  admin: ['user', 'admin'],
  needs_password_reset: ['user', 'needs_password_reset'],
  logged_in: ['user', 'logged_in'],
  is_authenticating: ['user', 'is_authenticating'],
}


export default UserStore;
exports.UserGetters = UserGetters;