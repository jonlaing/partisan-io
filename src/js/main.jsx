/*global data */
import React from 'react';
import ReactDOM from 'react-dom';
import Perf from 'react-addons-perf';
import jQuery from 'jquery';

import Events from 'events';
Events.EventEmitter.prototype._maxListeners = 100;

global.$ = jQuery;

Perf.start();

import Login from './components/Login.jsx';
import SignUp from './components/SignUp.jsx';
import Feed from './components/Feed.jsx';
import Questions from './components/Questions.jsx';
import ProfileShow from './components/ProfileShow.jsx';
import Matches from './components/Matches.jsx';
import Messages from './components/Messages.jsx';
import Friends from './components/Friends.jsx';
import ProfileEdit from './components/ProfileEdit.jsx';
import Post from './components/Post.jsx';
import Card from './components/Card.jsx';
import HashtagSearch from './components/HashtagSearch.jsx';
import FrontTicker from './components/FrontTicker.jsx';

// optionally attack DOM elements to React
let login = document.getElementById('login');
let signUp = document.getElementById('sign-up');
let feed = document.getElementById('feed');
let questions = document.getElementById('questions');
let profileShow = document.getElementById('profile-show');
let profileEdit = document.getElementById('profile-edit');
let matches = document.getElementById('matches');
let messages = document.getElementById('messages');
let friends = document.getElementById('friends');
let post = document.getElementById('post');
let hashtags = document.getElementById('hashtags');
let frontTicker = document.getElementById('front-ticker');

// for static login page
if(login !== null) {
  ReactDOM.render(<Login />, login);
}

if(signUp !== null) {
  ReactDOM.render(<SignUp />, signUp);
}

if(feed !== null) {
  ReactDOM.render(<Feed data={data} />, feed);
}


if(questions !== null) {
  ReactDOM.render(<Questions data={data}/>, questions);
}

if(profileShow !== null) {
  ReactDOM.render(<ProfileShow user={data.user} match={data.match} profile={data.profile} currentUser={data.current_user} />, profileShow);
}

if(profileEdit !== null) {
  ReactDOM.render(<ProfileEdit data={data} />, profileEdit);
}

if(matches !== null) {
  ReactDOM.render(<Matches data={data} />, matches);
}

if(messages !== null) {
  ReactDOM.render(<Messages data={data} />, messages);
}

if(friends !== null) {
  ReactDOM.render(<Friends data={data} />, friends);
}

if(post !== null) {
  let pData = data;
  pData.user = data.post_user;
  ReactDOM.render(<Card><Post data={data} defaultShowComments={true}/></Card>, post);
}

if(hashtags !== null) {
  ReactDOM.render(<HashtagSearch defaultSearch={data.search} />, hashtags);
}

if(frontTicker !== null) {
  ReactDOM.render(<FrontTicker />, frontTicker);
}

Perf.stop();
Perf.printWasted(Perf.getLastMeasurements());
