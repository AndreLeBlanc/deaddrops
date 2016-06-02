/**
 * @file       Defines the FileObj component which acts as a representation of a file in
 *             the ui. 
 */

import React from 'react';
import { connect } from 'react-redux';
import { FileState } from '../actions/FileActions';
import IconButton from './IconButton';
import FileInfo from './File/FileInfo';
import FileCommandBar from './File/FileCommandBar';


/**
 * FileObj component
 * @param  {Object} props Properties for the FileObj component.
 * @return {React.Component} Returns a react component.
 */
const FileObj = (props) => {

	return (
		<div class="file">
        <FileInfo state={props.state} name={props.name} size={props.size} maxSize={props.maxSize} />
        <FileCommandBar state={props.state} name={props.name} downloads={props.downloads} stashID={props.stashID} fileID={props.fileID} />
		</div>
	)
}

const mapDispatchToProps = (dispatch) => {
  return {
    preventDefault: (e) => {
      e.preventDefault();
      e.stopPropagation();
    }
  }
}

const mapStateToProps = (state) => {
  if(state.files)
  {
    return {
        maxSize: state.files.maxSize
      }
  }
  else
  {
     return {
        maxSize: 0
      }
  }
}

const FileComponent = connect (
	mapStateToProps,
	mapDispatchToProps
)(FileObj);

export default FileComponent;