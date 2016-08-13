import Dispatcher from '#app/Dispatcher'

export default {

  fetchInitialState() {
    var settings = {
      username: "marpaia"
    };
    Dispatcher.dispatch("RECEIVE_SETTINGS", settings)
  }

}