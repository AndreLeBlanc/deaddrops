/**
 * @FileCommandBar       Defines the FileCommandBar component which acts as a command bar for the fileObj
 *             the ui. 
 */

import React from 'react';
import { connect } from 'react-redux';
import { deleteFile, downloadFile, FileState, fetchDropData, updateOpenFileState } from '../../actions/FileActions';
import IconButton from '../IconButton';
import RippleButton from '../RippleButton';
import { CONFIG } from '../../utils/Utils';


/**
 * FileCommandBar component
 * @param  {Object} props Properties for the FileCommandBar component.
 * @return {React.Component} Returns a react component.
 */
const FileCommandBar = (props) => {
  return (
    <div class={getFileStateClass(props.state)}>
         <div name="state" class={(props.state !== FileState.FILE_READY) ? "enabled" : "disabled"}><div>{getStateObj(props.state)}</div></div>
         {getActionButton(props)}
    </div>
  )
}

function getActionButton(props)
{
  if(props.state !== FileState.FILE_READY)
  {
    return(
      <IconButton onClickFunc={props.deleteFile} clickAttr={props.name} iconName="cancel" enabled={(props.state === FileState.UPLOAD_WAITING || props.state === FileState.UPLOAD_FAILED || props.state === FileState.FILE_TOLARGE)} />
    );
  }
  else if(props.state === FileState.FILE_READY)
  {
    const buttonStyle = {
      'marginRight': '10px'
    };
    const downloadData = {
      StashID: props.stashID,
      FileID: props.fileID,
      Name: props.name
    };
    return(
      // TOOD: Take file_id as clickAttr
      //<RippleButton cssClass="stylized" action={props.downloadFile} customStyle={buttonStyle} label="Create Drop"></RippleButton>
        <a class="stylized flat" href={CONFIG.baseURL + "/download/"  + downloadData.StashID + "/" + downloadData.FileID} onClick={() => props.fetchDropData(props.stashID, downloadData.FileID)} download={downloadData.Name}>Download ({props.downloads})</a>
      //<RippleButton cssClass="stylized flat" rippleColor="#757575" action={props.downloadFile} actionAttr={downloadData} customStyle={buttonStyle} label="DOWNLOAD" />
    );
  }
  else
  {
    return (
        <div>
         </div>
      );
  }
}

/**
 * Computes the correct css class to use for displaying the current 
 * upload state of the File.
 * @param  {String} status Current upload state.
 * @return {String} Css class to use for the current upload state.
 */
function getFileStateClass(state) {
  switch(state)
  {
    case FileState.UPLOAD_WAITING:
      return "fileCommandBar";

    case FileState.UPLOADING:
      return "fileCommandBar";

    case FileState.UPLOAD_COMPLETE:
      return "fileCommandBar complete";

    case FileState.UPLOAD_FAILED:
    case FileState.FILE_TOLARGE:
      return "fileCommandBar fail";

    default:
      return "fileCommandBar";
  }
}

function getStateObj(state) {

  const opacStyle = {
    'opacity': '0.6'
  }

  const loadingStyle = {
      'marginRight': '10px'
    };

  switch(state)
  {
    case FileState.UPLOAD_WAITING:
      return (
          <i class="material-icons" style={opacStyle}>backup</i>
        );

    case FileState.FILE_DOWNLOADING:
    case FileState.UPLOADING:
      return (
            <div class="la-ball-pulse" style={loadingStyle}>
                <div></div>
                <div></div>
                <div></div>
             </div>
             );

    case FileState.UPLOAD_COMPLETE:
      return (
          <i class="material-icons" style={opacStyle}>done</i>
        );

    case FileState.UPLOAD_FAILED:
      return (
          <i class="material-icons" style={opacStyle}>error</i>
        );

    case FileState.FILE_TOLARGE:
      return (
          <i class="material-icons" style={opacStyle}>cloud_off</i>
        );

    default:
      return "";
  }

}

const mapDispatchToProps = (dispatch) => {
  return {
    preventDefault: (e) => {
    	e.preventDefault();
    	e.stopPropagation();
    },

    deleteFile: (e, fileName) => {
      e.preventDefault();
      e.stopPropagation();
      dispatch(deleteFile(fileName));
    },

    downloadFile: (e, downloadData) => {
      e.preventDefault();
      e.stopPropagation();
      dispatch(downloadFile(downloadData));
    },

    fetchDropData: (token, fileID) => {
      dispatch(updateOpenFileState(fileID, FileState.FILE_DOWNLOADING));
      setTimeout(() => dispatch(fetchDropData(token)), 2000);
    } 

  }
}


const FileCommandBarComponent = connect (
	null,
	mapDispatchToProps
)(FileCommandBar);

export default FileCommandBarComponent;