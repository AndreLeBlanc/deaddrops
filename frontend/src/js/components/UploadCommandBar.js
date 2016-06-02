/**
 * @UploadCommandBar       Defines the UploadCommandBar component which acts as a command bar for the dragdropzone
 *             the ui. 
 */

import React from 'react';
import { connect } from 'react-redux';
import { uploadFiles, FileState } from '../actions/FileActions';
import { updateDropState, DropState } from '../actions/Finalize';
import { formatBytes } from '../utils/Utils';
import IconButton from './IconButton';


/**
 * UploadCommandBar component
 * @param  {Object} props Properties for the UploadCommandBar component.
 * @return {React.Component} Returns a react component.
 */
const UploadCommandBar = (props) => {
	return (
		<div class="commandBar">
        {getPotentialFinalizeButton(props)}
        {getPotentialUploadButton(props)}
        <label class="addButton" for="fileinput" title="Add files"><i class="material-icons md-36">add</i></label>
     </div>
	)
}

function getPotentialFinalizeButton(props)
{
  if(props.uploadAmount > 0)
  {
    const buttonStyle = {
      'width': '36px',
      'height': '36px'
    };

    return (
      <button onClick={(e) => props.inializeDrop(e)} class="finalizeButton" title="Initialize drop"><img src="images/finalize.png" style={buttonStyle}/></button>
      )
  }
}

function getPotentialUploadButton(props)
{
  if(props.fileAmount > 0)
  {
    return (
       <button onClick={props.uploadFiles} class="uploadButton" title="Upload files"><i class="material-icons md-36">file_upload</i></button>
      )
  }
}

const mapDispatchToProps = (dispatch) => {
  return {
    preventDefault: (e) => {
    	e.preventDefault();
    	e.stopPropagation();
    },

     uploadFiles: (e) => {
      dispatch(uploadFiles());
    },

    inializeDrop: (e) => {
      dispatch(updateDropState(DropState.INITIALIZING_DROP));
    }
  }
}


const UploadCommandBarComponent = connect (
	null,
	mapDispatchToProps
)(UploadCommandBar);

export default UploadCommandBarComponent;