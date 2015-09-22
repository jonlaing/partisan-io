import React from 'react';

export default React.createClass({
  render() {
    return (
      <div className={this.props.className}>
        <button className="comment" onClick={this.props.onClick} >
          <i className="fi-comment"></i>
          {this.props.count}
        </button>
      </div>
    );
  }
});
