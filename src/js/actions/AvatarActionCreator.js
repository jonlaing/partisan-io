/*global FormData */
import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

export default {
  uploadAvatar(files) {
    var request = new FormData();

    console.log(files);
    files.forEach(function(value) {
      console.log(value);
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
        console.log(res);
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
