/**
 * @newdrop       Defines the header component
 *             the ui. 
 */
import React from 'react';
import { Link, hashHistory } from 'react-router';
import { connect } from 'react-redux';
import RippleButton from '../RippleButton';

const buttonStyle = {
	'marginLeft': '10px',
	'marginTop': '8px',
	'height': '30px'
}

/**
 * NewDropObj class
 * @param  {Object} props Properties for the NewDropObj class.
 * @return {React.Component} Returns a react component.
 */
export default class NewDropObj extends React.Component {

	constructor() {
		super();
		this.state = {
			label: ""
		}
		this.updateLabel = this.updateLabel.bind(this);
		this.routeToCreate = this.routeToCreate.bind(this);
		this.handleKeyPress = this.handleKeyPress.bind(this);
	}
	
	updateLabel(e) {
		this.setState({
			label: e.target.value
		})
	}

	handleKeyPress(e){
		/*enter*/
		if (e.keyCode == 13) {
		    this.routeToCreate();
		}
	}

	routeToCreate() {
		if (this.state.label.length > 0) {
		    hashHistory.push('/create/' + this.state.label);
		}
	}



	render() {
		const marginStyle = {
			'marginTop': '10px'
		}
		return (
			<div class="case">
				<input class="stylized" style={marginStyle} value={this.state.label} onChange={this.updateLabel} onKeyDown={this.handleKeyPress} placeholder="Dropname..." />
				<RippleButton cssClass="stylized" action={this.routeToCreate} customStyle={buttonStyle} label="Create Drop"></RippleButton>
			</div>
		)
	}	
}