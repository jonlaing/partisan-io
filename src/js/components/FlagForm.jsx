import React from 'react';

import FlagActionCreator from '../actions/FlagActionCreator';

import Modal from './Modal.jsx';

export default React.createClass({
  handleClose() {
    FlagActionCreator.cancelReport();
  },

  render() {
    return (
      <Modal show={this.props.show} onCloseClick={this.handleClose}>
        <h4>Flag a {this.props.type}</h4>
      </Modal>
    );
  }
});
