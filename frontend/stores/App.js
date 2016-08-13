import { Store, toImmutable } from 'nuclear-js'

 export default new Store({

  getInitialState() {
    return toImmutable({
      id: 0,
      username: "",
      email: "",
      name: ""
    })
  },

  initialize() {
    this.on("RECEIVE_SETTINGS", receiveSettings)
  }

})

function receiveSettings(state, settings) {
  return toImmutable({
    username: settings.username
  })

}

export var AppGetters = {
  id: ['app', 'id'],
  username: ['app', 'username'],
  email: ['app', 'email'],
  name: ['app', 'name']
}
