import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

export default {
  getQuestion() {
    $.ajax({
      url: Constants.APIROOT + '/questions',
      method: 'GET',
      dataType: 'json'
    })
      .done(function(data) {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.GET_QUESTION_SUCCESS,
          data: data
        });
      })
      .fail(function(res) {
        let data = res;
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.GET_QUESTION_FAIL,
          data: data
        });
      });
  },

  answerQuestion(question, agree) {
    let self = this;

    console.log(question.map);

    $.ajax({
      url: Constants.APIROOT + '/answers',
      data: JSON.stringify({
          "map": question.map,
          "agree": agree
      }),
      method: 'PATCH',
      dataType: 'json'
    })
      .done(function() {
        self.getQuestion();
      })
      .fail(function(res) {
        console.log(res);
      });
  }
};
