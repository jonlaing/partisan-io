import React from 'react';

import HashtagActionCreator from '../actions/HashtagActionCreator';
import HashtagStore from '../stores/HashtagStore';

import Card from './Card.jsx';
import Post from './Post.jsx';

let _ENTER = 13; // key code for pressing the ENTER/RETURN key

export default React.createClass({
  getInitialState() {
    return {search: this.props.defaultSearch, hits: []};
  },

  componentDidMount() {
    HashtagStore.addChangeListener(this._onChange);
    HashtagActionCreator.search(this.state.search);
  },

  componentWillUnmount() {
    HashtagStore.removeChangeListener(this._onChange);
  },

  handleChange(e) {
    this.setState({search: e.target.value});
  },

  handleKeyDown(e) {
    if(e.keyCode === _ENTER && e.target.value !== "") {
      HashtagActionCreator.search(this.state.search);
    }
  },

  render() {
    var list = this.state.hits.map((hit) => {
      return (
        <Card key={hit.post.id}>
          <Post data={hit} />
        </Card>
      );
    });

    return (
      <div className="hashtag-search">
        <div className="hashtag-search-bar">
          <input type="text" ref="search" onKeyDown={this.handleKeyDown} onChange={this.handleChange} defaultValue={this.props.defaultSearch} />
        </div>
        <div className="hashtag-search-list">
          {list}
        </div>
      </div>
    );
  },

  _onChange() {
    let state = HashtagStore.getAll();
    this.setState(state);
  }
});
