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
    } else if (this.state.friendship.user_id === this.props.id) {
      FriendsActionCreator.confirmFriendship(this.props.id);
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
    var text, header;

    if(this.state.friendship.id !== undefined && this.state.friendship.confirmed === true) {
      text = "Friends";
      header = "You and @" + this.props.username + " are friends";
    } else if (this.state.friendship.id !== undefined
       && this.state.friendship.confirmed === false
       && this.state.friendship.user_id === this.props.id) {

      text = "Confirm";
      header = "Confirm your friendship with @" + this.props.username;
    } else if (this.state.friendship.id !== undefined
       && this.state.friendship.confirmed === false
       && this.state.friendship.user_id !== this.props.id) {

      text = "Request Sent";
      header = "Awaiting @" + this.props.username + "'s confirmation";
    } else {
      text = "Add Friend";
      header = "You and @" + this.props.username + " are not friends";
    }

    return (
      <div className="profile-friend">
        <div className="right">
          <button onClick={this.toggleFriend} className={"friendship" + (this.state.isFriend ? " active" : "")}>
            {text}
          </button>
        </div>
        <h3>{header}</h3>
      </div>
    );
  },

  _onChange() {
    let state = FriendsStore.getFriend(this.props.id);
    this.setState({friendship: state});
  }
});
