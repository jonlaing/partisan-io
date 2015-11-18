import React from 'react/addons';

var ReactCSSTransitionGroup = React.addons.CSSTransitionGroup;

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
  "patriots",
  "everyone!"
];

export default React.createClass({
  getInitialState() {
    return {
      current: [_fullList[0]]
    };
  },

  componentDidMount() {
    _timer = window.setInterval(() => {
      let i = Math.floor(Math.random() * _fullList.length);
      this.setState({current: [_fullList[i]]});
    }, 2000);
  },

  componentWillUnmount() {
    window.clearInterval(_timer);
  },

  render() {
    let text = this.state.current.map((t, i) => {
      if(i === 0) { return <h3 key={Date.now()}>{t}</h3>; }
    });

    return (
      <div className="front-ticker">
        <h3>Partisan is for:</h3>
        <div className="front-ticker-tick">
          <ReactCSSTransitionGroup transitionName="front-ticker-text">
            {text}
          </ReactCSSTransitionGroup>
        </div>
      </div>
    );
  }
});
