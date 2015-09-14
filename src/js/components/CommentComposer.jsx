import React from 'react';

import CommentsActionCreator from '../actions/CommentsActionCreator';

export default React.createClass({
  handleSubmit() {
    let body = $(React.findDOMNode(this.refs.comment));
    CommentsActionCreator.create({
      "record_type": this.props.type,
      "record_id": this.props.id,
      "body": body.val()
    });
    body.val('');
  },

  render() {
    return (
      <div className="comment-composer">
        <div className="row collapse">
          <div className="large-10 columns">
            <textarea type="text" placeholder="Type your comment here..." ref="comment" ></textarea>
          </div>
          <div className="large-2 columns">
            <button onClick={this.handleSubmit}>Comment</button>
          </div>
        </div>
      </div>
    );
  }
});
