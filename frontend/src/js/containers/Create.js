import React from 'react';
import Axios from 'axios';
import {CONFIG} from '../utils/Utils';
import { Provider } from 'react-redux';
import { updateLabel } from '../actions/Finalize';
import configureCreateStore from '../configureCreateStore';
import CreateStateSwitcher from '../components/CreateStateSwitcher';
//const store = configureCreateStore();

export default class Create extends React.Component {
	
	constructor() {
		super();
		this.state = {
			store: configureCreateStore()
		}
	}

	componentWillMount (){

		this.state.store.dispatch(updateLabel(this.props.params.label));

		Axios.get(CONFIG.baseURL + '/create')
		.then((response) => {
			this.state.store.dispatch(
	      	{
	      	 	type: "SET_KEY", 
	      	 	key: response.data.Token
	      	});
	      	this.state.store.dispatch({type: "UPDATE_CREATE_STATE", state: "loaded"});
		}).catch((response) => {
			this.state.store.dispatch({type: "UPDATE_CREATE_STATE", state: "failed"});
			this.state.store.dispatch(
				{
					type: "SET_ERROR_MESSAGE", 
					message: "Unable to connect to server, Try again"
				}
			);	
		});
	}

	render () {
		
		return (
			<Provider store={this.state.store}>
					<CreateStateSwitcher />
			</Provider>
		)
	}
}