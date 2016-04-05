import React from 'react';
import ReactCSSTransitionGroup from 'react-addons-css-transition-group';

var _timer;

const _fullList = [
  "liberals",
  "conservatives",
  "anarchists",
  "communists",
  "libertarians",
  "Democrats",
  "Republicans",
  "socialists",
  "progressives",
  "marxists",
  "activists",
  "conspiracy theorists",
  "feminists",
  "traditionalists",
  "radicals",
  "patriots"
];

var _list = _scramble(_fullList);

export default React.createClass({
  getInitialState() {
    return {
      current: [_list[0]],
      index: 0
    };
  },

  componentDidMount() {
    _timer = window.setInterval(this._cycle, 2000);
  },

  componentWillUnmount() {
    window.clearInterval(_timer);
  },

  onWindowBlur() {
    window.clearInterval(_timer);
  },

  onWindowFocus() {
    _timer = window.setInterval(this._cycle, 2000);
  },

  render() {
    let text = this.state.current.map((t, i) => {
      if(i === 0) { return <h3 key={Date.now()}>{t}</h3>; }
    });

    return (
      <div className="front-ticker">
        <h3>Partisan is for:</h3>
        <div className="front-ticker-tick">
          <ReactCSSTransitionGroup transitionName="front-ticker-text" transitionEnterTimeout={2000} transitionLeaveTimeout={2000}>
            {text}
          </ReactCSSTransitionGroup>
        </div>
      </div>
    );
  },

  _cycle() {
    let index = this.state.index + 1;
    if(index === _list.length) {
      _list = _scramble(_fullList);
      index = 0;
    }

    this.setState({current: [_list[index]], index: index});
  }
});

function _scramble(arr) {
  var array = arr;
  var currentIndex = array.length, temporaryValue, randomIndex;

  // While there remain elements to shuffle...
  while (currentIndex !== 0) {

    // Pick a remaining element...
    randomIndex = Math.floor(Math.random() * currentIndex);
    currentIndex -= 1;

    // And swap it with the current element.
    temporaryValue = array[currentIndex];
    array[currentIndex] = array[randomIndex];
    array[randomIndex] = temporaryValue;
  }

  array.push("everyone!"); // always make sure "everyone" is at the end

  return array;
}
