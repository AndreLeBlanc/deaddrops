/**
 *@file app.js
 *@description Main application file. Acts as the entry point to the app.
 */
import React from "react";
import ReactDOM from "react-dom";
import {Router, Route, IndexRoute, hashHistory} from "react-router"
import App from './components/App';
import Home from './containers/Home';
import Create from './containers/Create'
import Open from './containers/Open'

const app = document.getElementById('app');

ReactDOM.render( 
	<Router history={hashHistory}>
		<Route path="/" component={App}>
			<IndexRoute component={Home}></IndexRoute>
			<Route path="/create/:label" component={Create}></Route>
			<Route path="/open/:deadropId" component={Open}></Route>
		</Route>
	</Router>, 
app);

