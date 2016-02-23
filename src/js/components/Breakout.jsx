import React from 'react';

export default React.createClass({
  render() {
    if(this.props.show === false) {
      return <span/>;
    }

    return (
      <div className={"breakout " + this.props.className}>
        <div className="breakout-arrow">
          <div className="breakout-arrow-inner">&nbsp;</div>
        </div>
        <div>
          {this.props.children}
        </div>
      </div>
    );
  }
});
