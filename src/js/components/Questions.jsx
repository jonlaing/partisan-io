import React from 'react/addons';
import QuestionsActionCreator from '../actions/QuestionsActionCreator';
import QuestionsStore from '../stores/QuestionsStore';
import Card from './Card.jsx';
import UserSession from './UserSession.jsx';
import Modal from './Modal.jsx';
import $ from 'jquery';

var ReactCSSTransitionGroup = React.addons.CSSTransitionGroup;

export default React.createClass({
  getInitialState() {
    return { set: -1, questions: [], questionsAnswered: 0, showModal: true};
  },

  handleAgree() {
    let last = this.state.questions.length - 1;
    $('.card').addClass('agree');
    QuestionsActionCreator.answerQuestion(this.state.questions[last], true);
  },

  handleDisagree() {
    let last = this.state.questions.length - 1;
    $('.card').addClass('disagree');
    QuestionsActionCreator.answerQuestion(this.state.questions[last], false);
  },

  handleModalClose() {
    this.setState({showModal: false});
  },

  componentDidMount() {
    QuestionsStore.addChangeListener(this._onChange);
    QuestionsActionCreator.getQuestions();
  },

  componentWillUnmount() {
    QuestionsStore.removeChangeListener(this._onChange);
  },

  render() {
    var cards = this.state.questions.map((q) => {
      return this._cardTemplate(q);
    });

    return (
      <div className="question">
        <div className="clearfix">
          <div className="right">
            <UserSession username={this.props.data.user.username} avatar={this.props.data.user.avatar_thumbnail_url} />
          </div>
          <img src="images/logo.svg" className="logo" />
        </div>
        <div className="question-container">
          <div className="question-body">
            <ReactCSSTransitionGroup transitionName="question-body">
              {cards}
            </ReactCSSTransitionGroup>
          </div>
        </div>
        <div className="question-actions">
          <button className="button disagree" onClick={this.handleDisagree}>Disagree</button>
          <button className="button agree" onClick={this.handleAgree}>Agree</button>
        </div>
        <Modal show={this.state.showModal} onCloseClick={this.handleModalClose} >
          <h2>Take the Quiz</h2>
          <div>You&apos;re about to be presented with <strong>20 prompts</strong>. Mark whether you agree or disagree<br/>with the statment. This is how we determine your beliefs and match you up with people<br/>who share your beliefs.</div>
          <br/>
          <br/>
          <div className="text-center">
            <button onClick={this.handleModalClose} >Get Started</button>
          </div>
        </Modal>
      </div>
    );
  },

  _cardTemplate(q) {
    if(q.prompt === undefined) {
      return '';
    }

    return (
      <Card key={q.prompt}>
        <div className="question-number">
          {this.state.questionsAnswered} of {this._maxQuestions()}
        </div>
        <div className="card-body">
          {q.prompt}
        </div>
      </Card>
    );
  },

  _onChange() {
    let question = QuestionsStore.getQuestion();
    let answered = this.state.questionsAnswered + 1;


    if(this._maxQuestions() > -1 && answered > this._maxQuestions()) {
      window.location.href = "/feed";
      return;
    }

    this.setState({questions: [question], questionsAnswered: answered});
  },

  _maxQuestions() {
    if(this.props.maxQuestions === undefined) {
      return 20;
    }

    return this.props.maxQuestions;
  }
});
