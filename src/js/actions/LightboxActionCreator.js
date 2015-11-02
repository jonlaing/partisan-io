import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

export default {
  open(image) {
    Dispatcher.handleViewAction({
      type: Constants.ActionTypes.LIGHTBOX_OPEN,
      image: image
    });
  },

  close() {
    Dispatcher.handleViewAction({
      type: Constants.ActionTypes.LIGHTBOX_CLOSE
    });
  }
};
