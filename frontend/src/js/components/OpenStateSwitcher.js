import React from 'react';
import { connect } from 'react-redux';
import OpenDropZone from './OpenDropZone';
import Loader from './Loader';
import ErrorDialog from './ErrorDialog';



const Switcher = (props) => {
		let openState = props.openState;
		if (openState === "initial") {
			return <Loader />		
		} else if (openState === "loaded") {
			return (
				<div class="createLayout">
					<label class="dropLabel">{props.stashLabel}</label>
					<OpenDropZone />
				</div>
				)
		} else {
			const flexStyle = {
				'display': 'flex',
				'justifyContent': 'center',
				'alignItems': 'center',
				'width': '100%',
    			'height': '100%'
			}

			return (
				<div style={flexStyle}>
					<ErrorDialog message={props.errorStateMsg} />
				</div>
				)
		}
}
				
const mapStateToProps = (state) => {
	return {
		openState: state.essentials.openState,
		errorStateMsg: state.essentials.errorStateMsg,
		stashLabel: state.essentials.dropData.StashName
	}
}

export default connect(mapStateToProps)(Switcher);
