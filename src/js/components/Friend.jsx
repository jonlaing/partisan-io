import React from 'react';

import FriendsActionCreator from '../actions/FriendsActionCreator';
import FriendsStore from '../stores/FriendsStore';

export default React.createClass({
  getInitialState() {
    return {friendship: {}, isFriend: false};
  },

  toggleFriend() {
    if(this.state.friendship.id === undefined) {
      FriendsActionCreator.requestFriendship(this.props.id);
    } else if (this.state.friendship.userID === this.props.id) {
      // confirm
    }
  },

  componentDidMount() {
    FriendsStore.addChangeListener(this._onChange);
    FriendsActionCreator.getFriendship(this.props.id);
  },

  componentWillUnmount() {
    FriendsStore.removeChangeListener(this._onChange);
  },

  render() {
    var text;

    console.log(this.state.friendship);

    if(this.state.friendship.id !== undefined && this.state.friendship.confirmed === true) {
      text = "Friends";
    } else if (this.state.friendship.id !== undefined
       && this.state.friendship.confirmed === false
       && this.state.friendship.user_id === this.props.id) {

      text = "Confirm";
    } else if (this.state.friendship.id !== undefined
       && this.state.friendship.confirmed === false
       && this.state.friendship.user_id !== this.props.id) {

      text = "Request Sent";
    } else {
      text = "Add Friend";
    }

    return (
      <button onClick={this.toggleFriend} className={"friendship" + (this.state.isFriend ? " active" : "")}>
        {text}
      </button>
    );
  },

  _onChange() {
    let state = FriendsStore.getFriend(this.props.id);
    this.setState({friendship: state});
  }
});
