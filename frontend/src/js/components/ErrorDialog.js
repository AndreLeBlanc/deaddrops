/**
 * @errordialog       Defines the FileObj component which acts as a representation of an error dialog
 *             the ui. 
 */

import React from 'react';
import { connect } from 'react-redux';


/**
 * FileObj component
 * @param  {Object} props Properties for the FileObj component.
 * @return {React.Component} Returns a react component.
 */
const ErrorDialogObj = (props) => {

	return (
		<div class="dialog">
      <div class="error">
        <i class="material-icons">cancel</i>
        <label>Error</label>
      </div>
      <div class="message">
        {props.message}
      </div>
		</div>
	)
}

export default ErrorDialogObj;