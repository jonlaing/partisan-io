/*global FormData */
import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

export default {
  uploadAvatar(files) {
    var request = new FormData();

    files.forEach(function(value) {
      request.append('avatar', value);
    });

    $.ajax({
      url: Constants.APIROOT + '/users/avatar_upload',
      data: request,
      cache: false,
      method: 'POST',
      processData: false,
      contentType: false
    })
      .done(function(res) {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.UPLOAD_AVATAR_SUCCESS,
          data: res
        });
      })
      .fail(function(res) {
        console.log(res);
      });
  }
};
