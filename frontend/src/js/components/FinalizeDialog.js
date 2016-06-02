import React from 'react';
import { connect } from 'react-redux';
import File from './File.js'
import { FileState } from '../actions/FileActions';
import { formatBytes } from '../utils/Utils';
import FinalizeCommandBar from './FinalizeCommandBar';


const Finalize = (props) => {

		return (
            <div class="finalizeDialogContainer">
            	<div class="finalizeDialog">
            		<div class="fileWrapper">
            		    <FinalizeCommandBar  size={calcSize(props.files)} nrOfFiles={props.fileObjects.length} />
            			{constructFileObjects(props.files)}
            		</div>
            	</div>
            </div>
		);
}

const calcSize = (files) => {
	let size = 0;
	for(let v of files.keys()) {
	    const file = files.get(v);
	    if (file.state === FileState.UPLOAD_COMPLETE) {
	    	size += file.data.size;
	    }
  	}
  	return formatBytes(size);
}

const constructFileObjects = (files) => {
  let fileObjects = [];
  for(let v of files.keys()) {
    const file = files.get(v);

    if (file.state === FileState.UPLOAD_COMPLETE) {
    	fileObjects.push(
      		<File key={file.data.name} name={file.data.name} state={file.state} size={file.data.size} />
    	);
    }
  }
  return fileObjects
}

const mapStateToProps = (state) => {
  return {
    files: state.files.files,
    fileObjects: constructFileObjects(state.files.files),
    value: ""
  }
}

const mapDispatchToProps = (dispatch) => {
  return {

  }
}

const FinalizeDialog = connect(
  mapStateToProps,
  mapDispatchToProps
)(Finalize);


export default FinalizeDialog;
