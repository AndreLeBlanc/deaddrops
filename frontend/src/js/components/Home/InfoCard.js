/**
 * @infocard       Defines the header component
 *             the ui. 
 */
import React from 'react';
import { connect } from 'react-redux';


/**
 * InfoCardObj component
 * @param  {Object} props Properties for the InfoCardObj component.
 * @return {React.Component} Returns a react component.
 */
const InfoCardObj = (props) => {
	return (
		<div class="infoCard">
			<div class="icon"><i class="material-icons">{props.iconName}</i></div>
			<div class="textWrapper">
				<div class="infoCardHeader">{props.header}</div>
				<div class="infoCardText">{props.text}</div>
			</div>
		</div>
	)
}

export default InfoCardObj;