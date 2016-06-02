import React from 'react';
import Axios from 'axios';
import {CONFIG} from '../utils/Utils';
import { Provider } from 'react-redux';
import { fetchDropData } from '../actions/FileActions';
import configureOpenStore from '../configureOpenStore';
import OpenStateSwitcher from '../components/OpenStateSwitcher';

export default class Open extends React.Component {

	constructor() {
		super();

		this.state = {
			store: configureOpenStore()
		}
	}

	componentWillMount() {
		this.state.store.dispatch(fetchDropData(this.props.params.deadropId));
	}

	render(){
		return (
			<Provider store={this.state.store}>
				<div class="createLayout">
					<OpenStateSwitcher />
				</div>
			</Provider>
		)	
	}
}

//