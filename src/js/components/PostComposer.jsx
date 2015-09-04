import React from 'react';

export default React.createClass({
  getInitialState() {
    return {};
  },

  componentDidMount() {
  },

  render() {
    return (
      <div className="post-composer">
        <div className="post-composer-field">
          <textarea rows="1" placeholder="Write a new post" rel="body"></textarea>
        </div>
        <div className="post-composer-actions clearfix">
          <button className="button right" onClick={this.handleCreate}>Post</button>
        </div>
      </div>
    );
  }
});
