import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

export default {
  getQuestions() {
    $.ajax({
      url: Constants.APIROOT + '/questions',
      method: 'GET',
      dataType: 'json'
    })
      .done(function(data) {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.GET_QUESTIONS_SUCCESS,
          data: data
        });
      })
      .fail(function(res) {
        let data = res;
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.GET_QUESTIONS_FAIL,
          data: data
        });
      });
  },

  answerQuestion(question, agree) {
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
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.QUESTION_ANSWERED_SUCCESS
        });
      })
      .fail(function(res) {
        console.log(res);
      });
  }
};
