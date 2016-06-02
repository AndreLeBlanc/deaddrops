/**
 * @FileInfo       Defines the FileCommandBar component which acts as a command bar for the fileObj
 *             the ui. 
 */

import React from 'react';
import { connect } from 'react-redux';
import { FileState } from '../../actions/FileActions';
import { formatBytes } from '../../utils/Utils';


/**
 * FileInfo component
 * @param  {Object} props Properties for the FileInfo component.
 * @return {React.Component} Returns a react component.
 */
const FileInfo = (props) => {
	if(props.state !== FileState.FILE_TOLARGE)
  {
    return (
      <div class="fileInfo">
        <label title={props.name}>{cutFileName(props.name)}</label>
        <p class="size">{formatBytes(props.size, 1)}</p>
      </div>
    )
  }
  else
  {
    return (
      <div class="fileInfo">
        <label title={props.name}>{cutFileName(props.name)}</label>
        <p class="size error">{"Too large (Max: " + formatBytes(props.maxSize, 1) + ") - " + formatBytes(props.size, 1)}</p>
      </div>
    )
  }
}

/**
 *  Truncates the name of a file if above 20 characters.
 * @param  {String} name Filename to possibly truncate.
 * @return {String}      Possibly truncated filename.
 */
function cutFileName(name) {
  if(name && name.length > 30)
    return name.substr(0, 30) + "...";
  else
    return name;
}

const mapDispatchToProps = (dispatch) => {
  return {
    preventDefault: (e) => {
    	e.preventDefault();
    	e.stopPropagation();
    }
  }
}


const FileInfoComponent = connect (
	null,
	mapDispatchToProps
)(FileInfo);

export default FileInfoComponent;