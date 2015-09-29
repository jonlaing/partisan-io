/*global data */
import React from 'react/addons';
import jQuery from 'jquery';

import Events from 'events';
Events.EventEmitter.prototype._maxListeners = 100;

global.$ = jQuery;
var Perf = React.addons.Perf;

Perf.start();

import Login from './components/Login.jsx';
import SignUp from './components/SignUp.jsx';
import Feed from './components/Feed.jsx';
import Questions from './components/Questions.jsx';
import ProfileShow from './components/ProfileShow.jsx';
import Matches from './components/Matches.jsx';
import ProfileEdit from './components/ProfileEdit.jsx';
import Post from './components/Post.jsx';
import Card from './components/Card.jsx';
import HashtagSearch from './components/HashtagSearch.jsx';

// optionally attack DOM elements to React
let login = document.getElementById('login');
let signUp = document.getElementById('sign-up');
let feed = document.getElementById('feed');
let questions = document.getElementById('questions');
let profileShow = document.getElementById('profile-show');
let profileEdit = document.getElementById('profile-edit');
let matches = document.getElementById('matches');
let post = document.getElementById('post');
let hashtags = document.getElementById('hashtags');

// for static login page
if(login !== null) {
  React.render(<Login />, login);
}

if(signUp !== null) {
  React.render(<SignUp />, signUp);
}

if(feed !== null) {
  React.render(<Feed data={data} />, feed);
}


if(questions !== null) {
  React.render(<Questions data={data}/>, questions);
}

if(profileShow !== null) {
  React.render(<ProfileShow user={data.user} match={data.match} profile={data.profile} />, profileShow);
}

if(profileEdit !== null) {
  React.render(<ProfileEdit data={data} />, profileEdit);
}

if(matches !== null) {
  React.render(<Matches data={data} />, matches);
}

if(post !== null) {
  React.render(<Card><Post data={data} defaultShowComments={true}/></Card>, post);
}

if(hashtags !== null) {
  React.render(<HashtagSearch defaultSearch={data.search} />, hashtags);
}

Perf.stop();
Perf.printWasted(Perf.getLastMeasurements());
