import Dispatcher from '../Dispatcher';
import Constants from '../Constants';
import BaseStore from './BaseStore';
import assign from 'object-assign';

// data storage
let _show = false;
let _image = "";

// Facebook style store creation.
const LightboxStore = assign({}, BaseStore, {
  // public methods used by Controller-View to operate on data
  getState() {
    return {
      show: _show,
      image: _image
    };
  },

  // register store with dispatcher, allowing actions to flow through
  dispatcherIndex: Dispatcher.register(function(payload) {
    let action = payload.action;

    switch(action.type) {
      case Constants.ActionTypes.LIGHTBOX_OPEN:
        let image = action.image.trim();
        if (image !== '') {
          _image = image;
          _show = true;
          LightboxStore.emitChange();
        }
        break;
      case Constants.ActionTypes.LIGHTBOX_CLOSE:
        _image = "";
        _show = false;
        LightboxStore.emitChange();
        break;

    }
  })
});

export default LightboxStore;
