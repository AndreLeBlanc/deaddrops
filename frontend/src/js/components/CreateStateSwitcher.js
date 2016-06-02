import React from 'react';
import { connect } from 'react-redux';
import DragAndDropZone from './DragAndDropZone';
import Loader from './Loader';
import ErrorDialog from './ErrorDialog';



const Switcher = (props) => {
		let createState = props.createState;
		if (createState === "initial") {
			return <Loader />		
		} else if (createState === "loaded") {
			return (
				<div class="createLayout">
					<label class="dropLabel">{props.stashLabel}</label>
					<DragAndDropZone />
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
		createState: state.essentials.state,
		errorStateMsg: state.essentials.errorStateMsg,
		stashLabel: state.essentials.label
	}
}

export default connect(mapStateToProps)(Switcher);
