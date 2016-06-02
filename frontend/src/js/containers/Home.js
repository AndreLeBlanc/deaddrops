import React from 'react';
import {Link} from 'react-router';
import Header from '../components/Home/Header';
import NewDrop from '../components/Home/NewDrop';
import ErrorDialog from '../components/ErrorDialog';

export default class Home extends React.Component {
	
	render () {
		return (
			<div class="home">
				<Header />
				<NewDrop />
			</div>
		)
	}
}