import React from 'react';
import { connect } from 'react-redux';
import ReactCSSTransitionGroup from 'react-addons-css-transition-group';
import { filesAdded, uploadFiles, FileState } from '../actions/FileActions';
import { finalizeDrop } from '../actions/Finalize';
import { constructFileObjects, uploadCompleteAmount, calcSize, countValidFiles } from '../utils/Utils';
import UploadCommandBar from './UploadCommandBar';
import FinalizeCommandBar from './FinalizeCommandBar';

const DragAndDrop = (props) => {

		const d = props.preventDefault;
		const f = props.files;
		return (
            <div class="dropzone" onDrag={d} onDragStart={d} onDragEnd={d} 
                onDragEnter={d} onDragExit={d} onDragLeave={d} onDragOver={d} onDrop={(e) => props.filesDropped(e, props.isOpen)}>
    			
    				  <input value={props.value} disabled={!props.isOpen} type="file" name="fileinput" id="fileinput" class="inputfile" multiple onChange={(e) => props.filesAdded(e, props.isOpen)} />
              <div class="fileWrapper">
                <FinalizeCommandBar nrOfFiles={countValidFiles(props.files)} maxFiles={props.maxFiles} size={calcSize(props.files)} />
                <ReactCSSTransitionGroup transitionName="file" transitionEnterTimeout={350} transitionLeaveTimeout={200}>
                  {props.fileObjects}
                </ReactCSSTransitionGroup>
              </div>
              {placeHolderText(props.files.size)}
              
            </div>
		);
}

//<UploadCommandBar uploadAmount={uploadCompleteAmount(f)} fileAmount={props.fileObjects.length} history={props.history} />
function placeHolderText(length)
{
  if(length < 1)
  {
    return (
        <div class="placeHolderDropZoneText">
        <span>DROP YOUR FILES HERE</span>
        </div>
      )
  }
}

const mapStateToProps = (state) => {
  return {
    files: state.files.files,
    isOpen: state.files.isOpen,
    fileObjects: constructFileObjects(state.files.files),
    maxFiles: state.files.maxFiles,
    value: ""
  }
}

const mapDispatchToProps = (dispatch) => {
  return {
    preventDefault: (e) => {
    	e.preventDefault();
    	e.stopPropagation();
    },

    filesAdded: (e, isOpen) => {
    	if (isOpen) {
     		dispatch(filesAdded(e.target.files));   		
    	}
    },

    filesDropped: (e, isOpen) => {
    	e.preventDefault();
    	e.stopPropagation();
    	if (isOpen) {
    		dispatch(filesAdded(e.dataTransfer.files));
    	}
    }
  }
}

const DragAndDropZone = connect(
  mapStateToProps,
  mapDispatchToProps
)(DragAndDrop);


export default DragAndDropZone;
      