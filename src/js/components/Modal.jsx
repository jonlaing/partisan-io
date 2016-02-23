import React from 'react';
import ReactCSSTransitionGroup from 'react-addons-css-transition-group';

export default React.createClass({
  handleClose() {
    this.props.onCloseClick();
  },

  render() {
    var modal;

    if(this.props.show) {
      modal = [(
        <div key={1}>
          <div className="modal-container">
            <div className={"modal " + this.props.className}>
              <a href="javascript:void(0)" className="right" onClick={this.handleClose}><i className="fi-x"></i></a>
              {this.props.children}
            </div>
          </div>
        </div>
      )];
    } else {
      modal = [];
    }

    return (
      <div>
        <ReactCSSTransitionGroup transitionName="modal-container" transitionEnterTimeout={1000} transitionLeaveTimeout={1000}>
          {modal}
        </ReactCSSTransitionGroup>
      </div>
    );
  }
});
