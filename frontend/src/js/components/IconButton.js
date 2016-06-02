import React from 'react';
import { connect } from 'react-redux';

const IconButton = (props) => {
  	return (
      <div class={(props.enabled) ? "enabled": "disabled"}>
    	 {getContent(props)}
      </div>
  	)
}

function getContent(props)
{
  if(props.enabled)
    return (
        <button onClick={(e) => props.onClickFunc(e, props.clickAttr)} class="iconButton"><i class="material-icons">{props.iconName}</i></button>
      )
}

export default IconButton;